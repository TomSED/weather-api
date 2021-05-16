package main

import (
	"fmt"
	"os"

	"github.com/TomSED/weather-api/pkg/postgres"
)

func main() {

	client, err := postgres.NewClient(os.Getenv("PG_RDS_HOST"), os.Getenv("PG_RDS_PORT"), os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DB_NAME"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = client.InitTables()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
