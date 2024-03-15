package main

import (
	usercore "github.com/usercoredev/usercore/app"
)

func main() {
	usercoreApp := usercore.Create()
	usercoreApp.ConfigureToken()
	usercoreApp.ConnectToDatabase()
	usercoreApp.SetupCache()
	usercoreApp.LoadClients()
	usercoreApp.StartServer()
}
