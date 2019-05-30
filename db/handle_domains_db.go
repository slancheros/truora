package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
	"truora/models"
	"truora/properties"
)

var db *gorm.DB

func Connect() *gorm.DB {
	config := properties.ObtainConfig()
	connString := fmt.Sprintf("postgresql://%s@%s:%s/%s?ssl=true&sslmode=disable", config.Database.User, config.Database.Server, config.Database.Port, config.Database.Database)
	dialect := "postgres"

	dbTemp, err := gorm.Open(dialect, connString)
	db = dbTemp
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ListDomainItems() models.Items {
	var items []models.Item
	db.Select("item").Find(&items)
	listItems := models.Items{Domains: items}
	return listItems
}

func ListItems() models.Items {
	var items []models.Item
	db.Select("item").Find(&items)
	listItems := models.Items{Domains: items}
	return listItems
}

func SaveQueriedDomain(domain string, info models.DomainInfo) {

	domain = domain + " info"
	currentDomains := ListItems()
	if !Contains(currentDomains.Domains, domain) {
		currentDomains.Domains = append(currentDomains.Domains, models.Item{Item: domain, QueryTime: time.Now(), SSLGrade: info.ServersSSLGrade})
		db.Save(&models.Item{Item: domain, QueryTime: time.Now(), SSLGrade: info.ServersSSLGrade})
	}
}

func Close() {
	db.Close()
}

func UpdateQueriedDomain(domain string, info models.DomainInfo) {

	domain_info := domain + " info"
	var item_to_update models.Item

	if db.Where("item=?", domain_info).First(&item_to_update).RecordNotFound() {
		SaveQueriedDomain(domain, info)
	} else {
		item_to_update.QueryTime = time.Now()
		item_to_update.SSLGrade = info.ServersSSLGrade
		db.Save(&item_to_update)
	}
}
func Contains(a []models.Item, x string) bool {
	for _, n := range a {
		if x == n.Item {
			return true
		}
	}
	return false
}
