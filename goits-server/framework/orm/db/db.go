package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cpekyaman/goits/config"
	"github.com/cpekyaman/goits/framework/monitoring"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type serverConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
}

type connectionConfig struct {
	MaxOpen  int           `mapstructure:"maxOpen"`
	MaxIdle  int           `mapstructure:"maxIdle"`
	LifeTime time.Duration `mapstructure:"lifeTime"`
}

type dbConfig struct {
	Server serverConfig     `mapstructure:"server"`
	Conn   connectionConfig `mapstructure:"conn"`
}

var appDB *sqlx.DB
var dbURL string
var conf dbConfig

// NewDB creates and initializes db layer of the application.
func NewDB() {
	config.ReadInto("db", &conf)

	dbURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Server.UserName, conf.Server.Password,
		conf.Server.Host, conf.Server.Port, conf.Server.Dbname)

	db, err := sqlx.Open("pgx", dbURL)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		monitoring.RootLogger().With(zap.Error(err)).Fatal("could not open database")
	}

	db.DB.SetMaxOpenConns(conf.Conn.MaxOpen)
	db.DB.SetMaxIdleConns(conf.Conn.MaxIdle)
	db.DB.SetConnMaxLifetime(conf.Conn.LifeTime * time.Second)

	appDB = db
}

// WithDB creates the db layer with pre-initialized db.
func WithDB(db *sql.DB, driverName string) {
	appDB = sqlx.NewDb(db, driverName)
}

// DB returns the db created.
func DB() *sqlx.DB {
	return appDB
}
