package gormx

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

var db *gorm.DB
var dbOnce sync.Once

type DatasourceConf struct {
	Dialect string
	ConnURL string
}

type Config interface {
	Config() (DatasourceConf, error)
}

func Datasource(c Config) *gorm.DB {
	dbOnce.Do(func() {
		var err error
		conf, err := c.Config()
		if err != nil {
			log.Fatal().Msgf("Failed to get config %s", err.Error())
		}
		db, err = gorm.Open(conf.Dialect, conf.ConnURL)
		if err != nil {
			log.Fatal().Msgf("Failed to connect database. Error: %s", err.Error())
		}

		if err := db.DB().Ping(); err != nil {
			log.Fatal().Msgf("Failed to connect database. Error: %s", err.Error())
		}
	})
	return db
}

func OverrideDatasource(d *gorm.DB) {
	db = d
}
