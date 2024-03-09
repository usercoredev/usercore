package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-sql-driver/mysql"
	gMysql "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/url"
	"os"
	"strings"
	"time"
)

var DB *gorm.DB

type Database struct {
	Engine          string
	Database        string
	DatabaseFile    string // For SQLite
	Charset         string
	User            string
	Password        string
	PasswordFile    string
	Host            string
	Port            string
	Certificate     string
	EnableMigration string
}

func (d *Database) Connect() (err error) {
	DB, err = d.configuration()
	if err != nil {
		return
	} else {
		fmt.Println("Database connection successful")
		if d.EnableMigration == "true" {
			fmt.Println("Migrating database")
			Migrate()
			fmt.Println("Database migration successful")
		}
	}

	return
}

func (d *Database) configuration() (*gorm.DB, error) {
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

	d.Password = url.QueryEscape(strings.TrimSpace(d.Password))

	if d.Engine == "sqlite" {
		if d.DatabaseFile == "" {
			return nil, fmt.Errorf("database file is required")
		}
		db, dbError := d.connectSQLite()
		return db, dbError
	}

	if d.Engine == "mysql" {
		db, dbError := d.connectMySQL()
		if dbError != nil {
			return nil, dbError
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		if err = sqlDB.Ping(); err != nil {
			return nil, err
		}
		return db, nil
	}

	if d.Engine == "postgres" {
		db, dbError := d.connectPostgres()
		return db, dbError
	}

	return nil, fmt.Errorf("unsupported database engine")
}

func (d *Database) connectSQLite() (db *gorm.DB, dbError error) {
	dbDSN := fmt.Sprintf("%s?parseTime=true&charset=utf8mb4", d.DatabaseFile)
	db, dbError = gorm.Open(sqlite.Open(dbDSN), &gorm.Config{})
	return
}

func (d *Database) connectMySQL() (*gorm.DB, error) {
	dbTcp := fmt.Sprintf("tcp(%s:%s)", d.Host, d.Port)
	dbAccess := fmt.Sprintf("%s:%s@%s/%s", d.User, d.Password, dbTcp, d.Database)

	tlsEnabled := "false"
	if d.Certificate != "" {
		caCertPool := x509.NewCertPool()
		caCert, err := os.ReadFile(d.Certificate)
		if err != nil {
			return nil, fmt.Errorf("failed to read certificate file: %w", err)
		}
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to append CA certificate to pool")
		}

		// Generate a unique identifier for this TLS configuration
		tlsConfigID := fmt.Sprintf("customTLS_%v", time.Now().UnixNano())

		// Setup TLS configuration
		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}

		// Register the TLS configuration with the unique identifier
		if err := mysql.RegisterTLSConfig(tlsConfigID, tlsConfig); err != nil {
			return nil, fmt.Errorf("failed to register TLS config: %v", err)
		}

		tlsEnabled = tlsConfigID
	}

	dbDSN := fmt.Sprintf("%s?parseTime=true&charset=%s&tls=%s", dbAccess, d.Charset, tlsEnabled)
	db, dbError := gorm.Open(gMysql.Open(dbDSN), &gorm.Config{})
	if dbError != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", dbError)
	}
	return db, nil
}

func (d *Database) connectPostgres() (db *gorm.DB, dbError error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.User, d.Password, d.Database, d.Port)
	db, dbError = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return
}

func Migrate() {
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
