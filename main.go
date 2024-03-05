package main

import (
	"flag"
	"github.com/usercoredev/usercore/app"
)

var (
	isDevelopment bool
	appMode       string
)

func init() {
	flag.StringVar(&appMode, "mode", "development", "Application mode")
}

func initializeApp() {
	switch appMode {
	case "development":
		isDevelopment = true
	case "production":
		isDevelopment = false
	default:
		panic("mode not set")
	}
}

func main() {
	flag.Parse()
	initializeApp()
	mainApp := app.App{
		GRPCServer: app.Server{
			Port: "9001",
		},
		HTTPServer: app.Server{
			Port: "8001",
		},
		Debug: isDevelopment,
	}

	//mainApp.SetTokenOptions()
	mainApp.ConnectToDatabase()
	//mainApp.ConnectToCache()
	//mainApp.LoadClients()
	//mainApp.StartServer()
}
