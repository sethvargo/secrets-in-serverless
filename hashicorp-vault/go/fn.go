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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var apikey string
var vaultAddr = os.Getenv("VAULT_ADDR")

func init() {
	jwt, err := fetchJwt()
	if err != nil {
		log.Fatal(err)
	}

	token, err := fetchToken(jwt)
	if err != nil {
		log.Fatal(err)
	}

	apikey, err = fetchApikey(token)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchJwt() (string, error) {
	client := new(http.Client)

	req, err := http.NewRequest(http.MethodGet, "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=http://vault/socialmedia&format=full", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetchToken(jwt string) (string, error) {
	client := new(http.Client)

	j := `{"role":"socialmedia", "jwt":"` + jwt + `"}`

	req, err := http.NewRequest(http.MethodPost, vaultAddr+"/v1/auth/gcp/login", bytes.NewBufferString(j))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	s := struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return "", err
	}

	return s.Auth.ClientToken, nil
}

func fetchApikey(token string) (string, error) {
	client := new(http.Client)

	req, err := http.NewRequest(http.MethodGet, vaultAddr+"/v1/secret/apikeys/twitter", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("x-vault-token", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	s := struct {
		Data struct {
			Value string `json:"value"`
		} `json:"data"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return "", err
	}

	return s.Data.Value, nil
}

func F(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, apikey)
}
