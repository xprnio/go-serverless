package main

import (
	"log"
	"os"
	"path"

	"github.com/xprnio/go-serverless/internal/server"
	"github.com/xprnio/go-serverless/internal/utils"
)

func main() {
	dataPath := os.Getenv("DATA_PATH")

	if utils.IsDocker() {
		if dataPath == "" {
			dataPath = "/data"
		}

		log.Println("running inside of docker")
	} else {
		if dataPath == "" {
			dataPath = "./"
		}

		log.Println("running outside of docker")
	}

	server, err := server.New(server.Options{
		DatabaseName: path.Join(dataPath, "database.db"),
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Start(":9999"); err != nil {
		log.Fatal(err)
	}
}
