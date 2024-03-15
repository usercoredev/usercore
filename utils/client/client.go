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
	clients, err := getClients(s.ClientFilePath)
	if err != nil {
		return err
	}
	s.Clients = clients
	return nil
}

func getClients(clientsFilePath string) (clients []Item, err error) {
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

func GetClient(id string, clients []Item) *Item {
	for _, client := range clients {
		if client.ID == id {
			return &client
		}
	}
	return nil
}
