package gormx

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var db *gorm.DB
var dbOnce sync.Once

func Datasource() *gorm.DB {
	if db == nil {
		dbOnce.Do(func() {
			var err error
			db, err = gorm.Open(viper.GetString("database.dialect"), viper.GetString("database.url"))
			if err != nil {
				log.Fatal().Msgf("Failed to connect database. Error: %s", err.Error())
			}

			if err := db.DB().Ping(); err != nil {
				log.Fatal().Msgf("Failed to connect database. Error: %s", err.Error())
			}
		})
	}
	return db
}

func OverrideDatasource(d *gorm.DB) {
	db = d
}
