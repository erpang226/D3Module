package main

import (
	"main/cmd/app"
	"os"
)

func main() {
	command := app.NewGatewayCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
