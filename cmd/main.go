package main

import (
	app "transact-api/application"
	"transact-api/constants"
)

func main() {
	// Run the application
	app.Run(constants.ConfigFile)
}
