package main

import (
	app "TransactEase/application"
	"TransactEase/constants"
)

func main() {
	// Run the application
	app.Run(constants.ConfigFile)
}
