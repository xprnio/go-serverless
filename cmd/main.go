package main

import (
	"fmt"
	"os"

	"github.com/xprnio/go-serverless/internal/server"
)

func main() {
	_, err := server.New(server.Options{
		DatabaseName: "database.db",
		Port:         9999,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
