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

func GetClients() ([]Client, error) {
	var list []Client
	jsonFile, err := os.Open("/run/secrets/clients")
	if err != nil {
		return nil, err
	}
	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &list)
	if err != nil {
		return nil, err
	}

	return list, err
}

func GetClient(id string, clients []Client) *Client {
	for _, client := range clients {
		if client.ID == id {
			return &client
		}
	}
	return nil
}
