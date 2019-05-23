package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"truora/models"
)

const addr = "postgresql://maxroach@localhost:26257/domains?ssl=true&sslmode=disable"

var db *gorm.DB

func Connect() *gorm.DB {

	dbTemp, err := gorm.Open("postgres", addr)
	db = dbTemp
	if err != nil {
		log.Fatal(err)
	}
	return db

}

func ListItems() models.Items {
	var items []models.Item
	db.Find(&items)
	listItems := models.Items{items}
	return listItems
}

func SaveQueriedDomain(domain string) {

	domain = domain + " info"
	currentDomains := ListItems()
	if !Contains(currentDomains.Domains, domain) {
		currentDomains.Domains = append(currentDomains.Domains, models.Item{domain})
		db.Save(&models.Item{domain})
	}
}

func Close() {
	db.Close()
}

func Contains(a []models.Item, x string) bool {
	for _, n := range a {
		if x == n.Item {
			return true
		}
	}
	return false
}
