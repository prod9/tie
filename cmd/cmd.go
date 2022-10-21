package cmd

import "log"

func Main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
