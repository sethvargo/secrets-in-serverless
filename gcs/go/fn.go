// Copyright 2018 Seth Vargo
// Copyright 2018 Google, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
