package sqlh

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Config struct {
	Dialect  string
	Host     string
	Port     int
	DbName   string
	Username string
	Password string
}

type DB struct {
	*sql.DB
}

const (
	Postgres = "postgres"
	MySql    = "mysql"
	Mock     = "mock"
)

func Open(c Config) (*sql.DB, error) {
	return sql.Open(c.Dialect, c.DNS())
}

func (c Config) DNS() string {
	switch c.Dialect {
	case Postgres:
		return c.postgres()
	case MySql:
		return c.mySql()
	default:
		return ""
	}
}

func (c Config) postgres() string {
	s := fmt.Sprintf("host=%s dbname=%s sslmode=disable", c.Host, c.DbName)

	if c.Username != "" {
		s += " user=" + c.Username
	}
	if c.Password != "" {
		s += " password=" + c.Password
	}
	if c.Port != 0 {
		s += " port=" + strconv.Itoa(c.Port)
	}
	return s
}

func (c Config) mySql() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port, c.DbName)
}

func Default(dialect string) Config {
	switch dialect {
	case Postgres:
		return Config{
			Dialect: dialect,
			Host:    "localhost",
			Port:    5432,
		}
	case MySql:
		return Config{
			Dialect: dialect,
			Host:    "localhost",
			Port:    3306,
		}
	default:
		return Config{
			Dialect: dialect,
		}
	}

}
