package main

import (
	"flag"
	"fmt"
	"github.com/usercoredev/usercore/app"
	"os"
)

func main() {
	flag.Parse()
	mainApp := app.App{
		GRPCServer: app.Server{
			Port: "9001",
		},
		HTTPServer: app.Server{
			Port: "8001",
		},
	}

	fmt.Println("App Name: ", os.Getenv("APP_NAME"))

	fmt.Println("Usercore under development..., mainApp.Debug: ", mainApp.GRPCServer.Port, mainApp.HTTPServer.Port)
	//mainApp.SetTokenOptions()
	//mainApp.ConnectToDatabase()
	//mainApp.ConnectToCache()
	//mainApp.LoadClients()
	//mainApp.StartServer()
}
