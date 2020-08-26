package config

import "github.com/jmoiron/sqlx"

var Mysql *sqlx.DB

func initMysql() {
	//var err error

}

func closeMysql() {
	Mysql.Close()
}
