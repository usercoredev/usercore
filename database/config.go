package database

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB
var dbError error
var dbName string
var dbUser string
var dbPassword string
var dbPort string

func configuration() (*gorm.DB, error) {
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	fmt.Println(dbName, dbUser, dbHost, dbPort)
	bin, err := os.ReadFile(os.Getenv("DB_PASSWORD"))
	if err != nil {
		return nil, err
	}
	dbPassword = string(bin)
	dbTcp := fmt.Sprintf("tcp(%s:%s)", dbHost, dbPort)
	dbAccess := fmt.Sprintf("%s:%s@%s", dbUser, dbPassword, dbTcp)

	dbDSN := fmt.Sprintf("%s/%s?parseTime=true&charset=utf8mb4", dbAccess, dbName)

	sqlDB, err := sql.Open("mysql", dbDSN)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	return gormDB, err
}

func Connect() error {
	DB, dbError = configuration()
	return dbError
}
func Migration() {
	err := DB.AutoMigrate(
		User{},
		Profile{},
		PasswordReset{},
		Device{},
		Session{},
		Role{},
		Permission{},
		SocialProvider{},

	)
	if err != nil {
		panic(err)
	}
}
