package configs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mini_3/models"
	"os"
)

// point 8
type MysqlDB struct {
	DB *gorm.DB
}

var Mysql MysqlDB

func OpenDB(silentLogger bool) {
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	MysqlConn, err := gorm.Open(mysql.Open(connString), &gorm.Config{})

	if silentLogger {
		MysqlConn, err = gorm.Open(mysql.Open(connString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	if err != nil {
		log.Fatal(err)
	}

	Mysql = MysqlDB{
		DB: MysqlConn,
	}

	Mysql.DB.Logger.LogMode(logger.Silent)
	MysqlConn.Logger.LogMode(logger.Silent)
	err = autoMigrate(Mysql.DB)
	if err != nil {
		return
	}
}

// point 9
func autoMigrate(db *gorm.DB) error {
	err := db.Migrator().DropTable(&models.Book{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		&models.Book{})

	if err != nil {
		return err
	}

	return err
}
