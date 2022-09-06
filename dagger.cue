package main

import (
	"dagger.io/dagger"
	"universe.dagger.io/docker"
)

dagger.#Plan & {
	client: {
		filesystem: ".": read: contents: dagger.#FS
		env: {
			GITHUB_USER:  string
			GITHUB_TOKEN: dagger.#Secret
			IMAGE:        "ghcr.io/prod9/tie:latest"
		}
	}

	_auth: {
		username: client.env.GITHUB_USER
		secret:   client.env.GITHUB_TOKEN
	}

	actions: {
		pull: docker.#Pull & {source: "alpine:edge"}

		build: docker.#Dockerfile & {
			source: client.filesystem.".".read.contents
			platforms: ["linux/amd64"]
		}

		push: docker.#Push & {
			image: actions.build.output
			dest:  "ghcr.io/prod9/api:latest"
			auth:  _auth
		}
	}
}
