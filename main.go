package main

import (
	"github.com/usercoredev/usercore/app"
	"os"
)

func main() {
	mainApp := app.App{
		TokenOptions: app.TokenOptions{
			Scheme:             os.Getenv("TOKEN_SCHEME"),
			PrivateKeyPath:     os.Getenv("PRIVATE_KEY_PATH"),
			PublicKeyPath:      os.Getenv("PUBLIC_KEY_PATH"),
			AccessTokenExpire:  os.Getenv("ACCESS_TOKEN_EXPIRE"),
			RefreshTokenExpire: os.Getenv("REFRESH_TOKEN_EXPIRE"),
		},
		GRPCServer: app.Server{
			Port: os.Getenv("GRPC_SERVER_PORT"),
		},
		HTTPServer: app.Server{
			Port: os.Getenv("HTTP_SERVER_PORT"),
		},
		DatabaseOptions: app.DatabaseOptions{
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			User:         os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			PasswordFile: os.Getenv("DB_PASSWORD_FILE"),
			Database:     os.Getenv("DB_NAME"),
			DatabaseFile: os.Getenv("DB_FILE_PATH"),
			Engine:       os.Getenv("DB_ENGINE"),
		},
		CacheOptions: app.CacheOptions{
			Enabled: os.Getenv("CACHE_ENABLED"),
			Host:    os.Getenv("CACHE_HOST"),
			Port:    os.Getenv("CACHE_PORT"),
		},
		Client: app.Client{
			ClientFile: os.Getenv("CLIENTS_FILE_PATH"),
		},
	}
	mainApp.ConfigureToken()
	//mainApp.ConnectToDatabase()
	mainApp.Cache()
	mainApp.LoadClients()
	mainApp.StartServer()
}
