package database

import (
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var db *gorm.DB
var once sync.Once

// Init initializes the database connection and performs migrations
func Init() *gorm.DB {
	once.Do(func() {
		// Load configuration
		env := "development"
		env = os.Getenv("ENVIRONMENT")
		configFile := "config/" + env + ".yaml"
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}

		// Get database configuration
		dbConfig := viper.GetStringMapString("database")

		// Connect to PostgreSQL
		db, err = gorm.Open("postgres", "host="+dbConfig["host"]+
			" port="+dbConfig["port"]+
			" user="+dbConfig["user"]+
			" dbname="+dbConfig["dbname"]+
			" password="+dbConfig["password"]+
			" sslmode=disable")
		if err != nil {
			panic(err)
		}
	})

	return db
}

// GetDB returns the initialized database instance
func GetDB() *gorm.DB {
	return db
}
