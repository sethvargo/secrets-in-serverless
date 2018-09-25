package example

import (
	"fmt"
	"net/http"
	"os"
)

var username = os.Getenv("DB_USER")
var password = os.Getenv("DB_PASS")

func F(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, username+":"+password)
}
