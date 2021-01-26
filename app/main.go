package main

import (
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
)

func main() {
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")
	dbSecretPath := os.Getenv("DB_SECRET_PATH")

	fmt.Printf("addr: %s, token: %s, path: %s\n", vaultAddr, vaultToken, dbSecretPath)

	config := vault.Config{Timeout: 5 * time.Second}
	err := config.ReadEnvironment()
	if err != nil {
		fmt.Printf("Can't read vault config: %v", err)
		return
	}
	vaultClient, err := vault.NewClient(&config)
	if err != nil {
		fmt.Printf("Can't create a new vault client instance: %v", err)
		return
	}

	var dbURL interface{}
	var isExists bool
	for attempt := 0; attempt < 5; attempt++ {
		dbSecret, err := vaultClient.Logical().Read(dbSecretPath)
		if err != nil {
			fmt.Printf("Can't read secret from vault: %v", err)
			return
		}

		dbURL, isExists = dbSecret.Data["db_url"]
		if isExists {
			break
		}

		time.Sleep(time.Duration(attempt+1) * time.Second)
	}

	if !isExists {
		fmt.Printf("Can't find db_url key in the secret")
	} else {
		fmt.Printf("The secret from vault is: %v", dbURL)
	}
}
