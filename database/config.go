package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

type Database struct {
	Engine       string
	Database     string
	DatabaseFile string // For SQLite
	Charset      string
	User         string
	Password     string
	PasswordFile string
	Host         string
	Port         string
}

func (d *Database) Connect() (dbError error) {
	DB, dbError = d.configuration()
	return
}

func (d *Database) configuration() (db *gorm.DB, dbError error) {
	if d.PasswordFile == "" && d.Password == "" {
		return nil, fmt.Errorf("password or password file is required")
	}
	if d.PasswordFile != "" {
		bin, err := os.ReadFile(d.PasswordFile)
		if err != nil {
			return nil, err
		}
		d.Password = string(bin)
	}

	if d.Engine == "sqlite" {
		if d.DatabaseFile == "" {
			return nil, fmt.Errorf("database file is required")
		}
		db, dbError = d.ConnectSQLite()
		return
	}

	if d.Engine == "mysql" {
		db, dbError = d.ConnectMySQL()
		return
	}

	if d.Engine == "postgres" {
		db, dbError = d.ConnectPostgres()
		return
	}

	return nil, fmt.Errorf("unsupported database engine")
}

func (d *Database) ConnectSQLite() (db *gorm.DB, dbError error) {
	dbDSN := fmt.Sprintf("%s?parseTime=true&charset=utf8mb4", d.DatabaseFile)
	db, dbError = gorm.Open(sqlite.Open(dbDSN), &gorm.Config{})
	return
}

func (d *Database) ConnectMySQL() (db *gorm.DB, dbError error) {
	dbTcp := fmt.Sprintf("tcp(%s:%s)", d.Host, d.Port)
	dbAccess := fmt.Sprintf("%s:%s@%s", d.User, d.Password, dbTcp)
	dbDSN := fmt.Sprintf("%s/%s?parseTime=true&charset=%s", dbAccess, d.Database, d.Charset)
	db, dbError = gorm.Open(mysql.Open(dbDSN), &gorm.Config{})
	return
}

func (d *Database) ConnectPostgres() (db *gorm.DB, dbError error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.User, d.Password, d.Database, d.Port)
	db, dbError = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return
}

/*

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
*/
