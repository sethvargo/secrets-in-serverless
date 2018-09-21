package example

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/vault/api"
)

var username string
var password string

func init() {
	config := api.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %s", err)
	}

	secret, err := client.Logical().Unwrap(os.Getenv("VAULT_TOKEN"))
	if err != nil || secret == nil || secret.Data == nil {
		// monitoringsystem.SendAlert("...")
		log.Fatalf("failed to unwrap: %s", err)
	}

	data := secret.Data["data"].(map[string]interface{})
	username = data["username"].(string)
	password = data["password"].(string)
}

func Secrets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, username+":"+password)
}
