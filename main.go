package main

import "github.com/csh0101/netagent.git/cmd"

func main() {
	rootCmd := cmd.RootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
