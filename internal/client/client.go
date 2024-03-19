package client

import (
	"encoding/json"
	"io"
	"os"
)

type clientKey string

var Key clientKey = "client"

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Settings struct {
	Clients        []Item
	ClientFilePath string
}

func (s *Settings) GetClient(id string) *Item {
	for _, client := range s.Clients {
		if client.ID == id {
			return &client
		}
	}
	return nil
}

func (s *Settings) LoadClients() error {
	var clients []Item
	jsonFile, err := os.Open(s.ClientFilePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &clients)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	s.Clients = clients
	return nil
}
