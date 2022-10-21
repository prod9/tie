package main

import (
	"dagger.io/dagger"
	"dagger.io/dagger/core"
	"encoding/yaml"
	"universe.dagger.io/docker"
	"strings"
)

let params = {
	github_user: "chakrit"
	name:        "prod9/tie"
	image:       "ghcr.io/\(name)"
}

dagger.#Plan & {
	client: {
		filesystem: {
			".": {
				read: contents:  dagger.#FS
				write: contents: actions.github.ci_yaml.output
			}
		}
		commands: {
			git_rev: {
				name: "git"
				args: ["rev-parse", "--short", "HEAD"]
			}
		}
		env: {
			GITHUB_USER:  string | *params.github_user
			GITHUB_TOKEN: dagger.#Secret
		}
	} // client

	actions: {
		build: docker.#Dockerfile & {
			source: client.filesystem.".".read.contents
			platforms: ["linux/amd64"]
		}

		push: {
			let _rev = strings.TrimSpace(client.commands.git_rev.stdout)
			let _push = docker.#Push & {
				image: actions.build.output
				auth: {
					username: client.env.GITHUB_USER
					secret:   client.env.GITHUB_TOKEN
				}
			}

			latest: _push & {dest: "\(params.image):latest"}
			rev:    _push & {dest: "\(params.image):\(_rev)"}
		}

		github: {
			mkdir: core.#Mkdir & {
				input: dagger.#Scratch
				path:  "./.github/workflows"
			}
			ci_yaml: core.#WriteFile & {
				input:    mkdir.output
				path:     "./.github/workflows/ci.yaml"
				contents: yaml.Marshal(github_ci)
			}
		} // github
	} // actions
}

let github_ci = {
	name: "dagger"
	on:   "push"
	jobs: {
		deploy: {
			"runs-on": "self-hosted"
			env: {
				GITHUB_USER:       params.github_user
				GITHUB_TOKEN:      "${{ github.token }}"
				DAGGER_CACHE_FROM: "type=gha,scope=\(params.name)"
				DAGGER_CACHE_TO:   "type=gha,mode=max,scope=\(params.name)"
			}
			steps: [
				{
					name: "checkout"
					uses: "actions/checkout@v3"
				},
				{
					name: "dagger update"
					uses: "dagger/dagger-for-github@v3"
					with: cmds: "project update"
				},
				{
					name: "dagger push"
					if:   "${{ github.ref_type == 'branch' && github.ref_name == 'main' }}"
					uses: "dagger/dagger-for-github@v3"
					with: cmds: "do push"
				},
			] // steps
		} // deploy
	} // jobs
}
