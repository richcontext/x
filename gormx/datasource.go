package gormx

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var db *gorm.DB
var dbOnce sync.Once

func Datasource() (*gorm.DB, error) {
	var err error
	if db == nil {
		dbOnce.Do(func() {
			db, err = gorm.Open(viper.GetString("database.dialect"), viper.GetString("database.url"))
			if err != nil {
				log.Error().Msgf("Failed to connect database. Error: %s", err.Error())

			}

			if err = db.DB().Ping(); err != nil {
				log.Error().Msgf("Failed to connect database. Error: %s", err.Error())
			}
		})
	}
	return db, err
}

func OverrideDatasource(d *gorm.DB) {
	db = d
}
