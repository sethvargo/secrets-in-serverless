package example

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigtable"
)

var project = os.Getenv("BIGTABLE_PROJECT")
var instance = os.Getenv("BIGTABLE_INSTANCE")
var table = os.Getenv("BIGTABLE_TABLE")
var column = "hello_world"
var client *bigtable.Client

func init() {
	var err error
	client, err = bigtable.NewClient(context.Background(), project, instance)
	if err != nil {
		log.Fatalf("failed to create client: %s", err)
	}
}

func F(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	row, err := client.Open(table).ReadRow(ctx, "my-key", bigtable.RowFilter(bigtable.ColumnFilter(column)))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "err: %s", err)
	}
	fmt.Fprintf(w, "%#v", row)
}
