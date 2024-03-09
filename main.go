package main

import (
	usercore "github.com/usercoredev/usercore/app"
)

func main() {
	usercore.Create()
	usercore.App.ConfigureToken()
	usercore.App.ConnectToDatabase()
	usercore.App.Cache()
	usercore.App.LoadClients()
	usercore.App.StartServer()
}
