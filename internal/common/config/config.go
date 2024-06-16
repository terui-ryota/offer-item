package config

import (
	"database/sql"
	"fmt"
	"time"

	"contrib.go.opencensus.io/integrations/ocsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Database struct {
	Name        string `yaml:"dbname"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Debug       bool   `yaml:"debug"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	MaxIdleTime string `yaml:"max_idle_time"`
}

type ExternalContexts struct {
	AffiliateItem *ContextConfig `yaml:"affiliate_item"`
	MediaAmeba    *ContextConfig `yaml:"media_ameba"`
	Affiliator    *ContextConfig `yaml:"affiliator"`
}

type ContextConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func LoadDB(config *Database) *sql.DB {
	// DB初期化
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Asia%%2FTokyo",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	// OpenCensus
	driverName, err := ocsql.Register("mysql")
	if err != nil {
		fmt.Println("==================")
		fmt.Println(err)
		fmt.Println("==================")
		panic(err)
	}

	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(config.MaxIdleConn)
	db.SetMaxOpenConns(config.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(1*config.MaxOpenConn) * time.Second)
	maxIdleTime, err := time.ParseDuration(config.MaxIdleTime)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxIdleTime(maxIdleTime)

	boil.SetDB(db)
	boil.DebugMode = config.Debug

	return db
}
