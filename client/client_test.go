package client

import (
	"testing"
)

func TestNewDingTalkClient(t *testing.T) {
	credential := Credential{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
	}

	client := NewDingTalkClient(credential)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.Credential.ClientID != credential.ClientID {
		t.Errorf("Expected ClientID to be %s, got %s", credential.ClientID, client.Credential.ClientID)
	}

	if client.Credential.ClientSecret != credential.ClientSecret {
		t.Errorf("Expected ClientSecret to be %s, got %s", credential.ClientSecret, client.Credential.ClientSecret)
	}
}

func TestNewDingTalkClientManager(t *testing.T) {
	credentials := []Credential{
		{
			ClientID:     "client1",
			ClientSecret: "secret1",
		},
		{
			ClientID:     "client2",
			ClientSecret: "secret2",
		},
	}

	manager := NewDingTalkClientManager(credentials)

	if manager == nil {
		t.Fatal("Expected manager to be created, got nil")
	}

	if len(manager.Clients) != 2 {
		t.Errorf("Expected 2 clients, got %d", len(manager.Clients))
	}

	client1 := manager.GetClientByOAuthClientID("client1")
	if client1 == nil {
		t.Error("Expected to get client1, got nil")
	}

	client3 := manager.GetClientByOAuthClientID("client3")
	if client3 != nil {
		t.Error("Expected nil for non-existent client, got a client")
	}
}
