package client

import (
	"encoding/json"
	"io"
	"os"
)

type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetClients(clientsFilePath string) (clients []Client, err error) {
	jsonFile, err := os.Open(clientsFilePath)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteValue, &clients)
	if err != nil {
		return
	}
	return
}

func GetClient(id string, clients []Client) *Client {
	for _, client := range clients {
		if client.ID == id {
			return &client
		}
	}
	return nil
}
