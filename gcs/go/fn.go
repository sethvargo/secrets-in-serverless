package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
)

var username string
var password string
var bucketName = os.Getenv("STORAGE_BUCKET")

func init() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create client: %s", err)
	}

	rc, err := client.Bucket(bucketName).Object("app1").NewReader(ctx)
	if err != nil {
		log.Fatalf("failed to get object: %s", err)
	}
	defer rc.Close()

	var t struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(rc).Decode(&t); err != nil {
		log.Fatalf("failed to decode object: %s", err)
	}

	username = t.Username
	password = t.Password
}

func F(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, username+":"+password)
}
