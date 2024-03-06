package client

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetClients(clientsFilePath string) ([]Client, error) {
	var list []Client
	filePath, err := filepath.Abs(clientsFilePath)
	if err != nil {
		return nil, err
	}
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

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
