package model

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"errors"
	"fmt"
	"log"
	"time"
)

var (
	app         = "x"
	Database    *gorm.DB
	HasDatabase bool

	DataType = "sqlite"
	//DatabaseConn        = "root:rootpass@/wp?charset=utf8"
	//DatabaseConn = "./content/data/sqlite.db"
	DatabaseConn        = "../content/data/sqlite_test.db"
	DatabaseTablePrefix = "zp_"
)

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      model.Model
//    }
type Model struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func ConnDatabase() (*gorm.DB, error) {
	switch {
	case DataType == "sqlite":
		if DatabaseConn[:1] == "." {
			fmt.Println("DatabaseConn>", DatabaseConn)
			return gorm.Open("sqlite3", DatabaseConn)
		}
		return gorm.Open("sqlite3", app+"/data/sqlite.db")

	case DataType == "mysql":
		return gorm.Open("mysql", DatabaseConn)

	case DataType == "postgres":
		return gorm.Open("postgres", DatabaseConn)
	}
	return nil, errors.New("unknown database type.")
}

func SetDatabase(maxIdleConns, maxOpenConns int) (*gorm.DB, error) {
	var _error error
	if Database, _error = ConnDatabase(); _error != nil {
		return nil, fmt.Errorf("Failed to connect database: %s,%s", _error.Error(), DatabaseConn)
	}

	//Database Connection Pool
	Database.DB().SetMaxIdleConns(maxIdleConns)
	Database.DB().SetMaxOpenConns(maxOpenConns)
	Database.Debug().LogMode(true)
	return Database, nil
}

func NewDatabase(maxIdleConns, maxOpenConns int) error {
	var _error error
	Database, _error = SetDatabase(maxIdleConns, maxOpenConns)
	return _error
}

func init() {

	var _error error
	if Database, _error = SetDatabase(10, 100); _error != nil {
		log.Fatal("app.models.init() errors:", _error.Error())
	}

	CreateTables(Database)

	message()

}

func CreateTables(Database *gorm.DB) {
	if DataType == "mysql" {
		Database.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&AppVersions{}, &Apps{}, &Commentmeta{}, &Comments{}, &Links{}, &Options{}, &Postmeta{}, &Posts{}, &RegistrationLog{}, &Signups{}, &Site{}, &Sitemeta{}, &TermRelationships{}, &TermTaxonomy{}, &Termmeta{}, &Terms{}, &Usermeta{}, &User{}, &Role{}, &Permission{}, &RolePermissions{})
	} else {
		Database.AutoMigrate(&AppVersions{}, &Apps{}, &Commentmeta{}, &Comments{}, &Links{}, &Options{}, &Postmeta{}, &Posts{}, &RegistrationLog{}, &Signups{}, &Site{}, &Sitemeta{}, &TermRelationships{}, &TermTaxonomy{}, &Termmeta{}, &Terms{}, &Usermeta{}, &User{}, &Role{}, &Permission{}, &RolePermissions{})
	}
}

func Ping() error {
	return Database.DB().Ping()
}

func message() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("The x system has started!")
}
