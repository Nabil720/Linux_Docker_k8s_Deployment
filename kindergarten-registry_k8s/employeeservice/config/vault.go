package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Address string
	Token   string
	Client  *api.Client
}

type MongoDBConfig struct {
	URI      string `json:"uri"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type APMConfig struct {
	ServerURL   string `json:"server_url"`
	SecretToken string `json:"secret_token"`
	Environment string `json:"environment"`
}

type ServiceConfig struct {
	MongoDB MongoDBConfig
	APM     APMConfig
	Port    int
}

func InitVaultClient() (*VaultConfig, error) {
	// Environment variables থেকে Vault configuration নিবে
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://192.168.121.132:8200"
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		return nil, fmt.Errorf("VAULT_TOKEN environment variable is required")
	}

	config := &api.Config{
		Address: vaultAddr,
		Timeout: 10 * time.Second,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %v", err)
	}

	client.SetToken(vaultToken)

	return &VaultConfig{
		Address: vaultAddr,
		Token:   vaultToken,
		Client:  client,
	}, nil
}

func GetSecrets(vc *VaultConfig, serviceName string) (*ServiceConfig, error) {
	ctx := context.Background()

	// MongoDB secrets পড়ো
	mongoSecret, err := vc.Client.KVv2("kindergarten").Get(ctx, "mongodb")
	if err != nil {
		return nil, fmt.Errorf("failed to read MongoDB secrets: %v", err)
	}

	// APM secrets পড়ো
	apmSecret, err := vc.Client.KVv2("kindergarten").Get(ctx, "apm")
	if err != nil {
		return nil, fmt.Errorf("failed to read APM secrets: %v", err)
	}

	// Services secrets পড়ো
	servicesSecret, err := vc.Client.KVv2("kindergarten").Get(ctx, "services")
	if err != nil {
		return nil, fmt.Errorf("failed to read services secrets: %v", err)
	}

	// Ports secrets পড়ো
	portsSecret, err := vc.Client.KVv2("kindergarten").Get(ctx, "ports")
	if err != nil {
		return nil, fmt.Errorf("failed to read ports secrets: %v", err)
	}

	// Parse MongoDB config
	mongoData, err := json.Marshal(mongoSecret.Data)
	if err != nil {
		return nil, err
	}

	var mongoConfig MongoDBConfig
	if err := json.Unmarshal(mongoData, &mongoConfig); err != nil {
		return nil, err
	}

	// Parse APM config
	apmData, err := json.Marshal(apmSecret.Data)
	if err != nil {
		return nil, err
	}

	var apmConfig APMConfig
	if err := json.Unmarshal(apmData, &apmConfig); err != nil {
		return nil, err
	}

	// Get service-specific port
	var port int
	switch serviceName {
	case "student":
		port = portsSecret.Data["student"].(int)
	case "teacher":
		port = portsSecret.Data["teacher"].(int)
	case "employee":
		port = portsSecret.Data["employee"].(int)
	default:
		port = 5000
	}

	return &ServiceConfig{
		MongoDB: mongoConfig,
		APM:     apmConfig,
		Port:    port,
	}, nil
}