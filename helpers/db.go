package helpers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type DbConnection struct {
	*gorm.DB
	dbUrl  string
	schema string
}

func InitDb(cfg Config) DbConnection {
	db := DbConnection{}
	db.dbUrl = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", cfg.Db.Host, cfg.Db.User, cfg.Db.Name, cfg.Db.Pass)
	db.Connect()

	return db
}

func (d *DbConnection) Connect() {
	if d.DB != nil {
		check := d.DB.DB().Ping()
		for check != nil {
			log.Println("Connection has been lost. Attempt to reconnect after 5 seconds...")
			_ = d.DB.Close()
			d.DB = nil
			time.Sleep(5 * time.Second)
			log.Println("Reconnecting...")
			d.Connect()
			check = d.DB.DB().Ping()
		}
	} else {
		db, err := gorm.Open("postgres", d.dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		d.DB = db
	}
}
