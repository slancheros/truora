package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"truora/models"
)

const addr = "postgresql://maxroach@localhost:26257/bank?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.maxroach.key&sslcert=certs/client.maxroach.crt"

var db *gorm.DB

func connect() {

	db, err := gorm.Open("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func autoMigrate() {

	db.AutoMigrate(&models.DomainHistory{})

}

func listItems() {
	// Print out the balances.
	var items []models.DomainHistory
	db.Find(&items)
	fmt.Println("Initial balances:")
	for _, item := range items {
		fmt.Println("%d %d\n", item.Domain, item.RequestTime)
	}

}
