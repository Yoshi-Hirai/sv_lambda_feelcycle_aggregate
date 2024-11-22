// db DBシステムパッケージ
package db // パッケージ名はディレクトリ名と同じにする

import (
	"database/sql"
	"log/slog"
	"sv_lambda_feelcycle_aggregate/config"

	"github.com/go-sql-driver/mysql"
)

// ---- Global Variable

// ---- Package Global Variable
var db *sql.DB

//---- public function ----

// DbBaseInit (public)DataBaseシステムの初期化関数。
func DbBaseInit() error {

	// 接続情報を外出しする
	dbProperty := config.GetConfigInformation()
	var errOpen, errPing error
	driverName := dbProperty.DbDriver
	//	dataSourceName := "host=127.0.0.1 port=5432 user=postgres dbname=sandbox0 sslmode=disable"
	c := mysql.Config{
		DBName:               dbProperty.DbName,
		User:                 dbProperty.DbUser,
		Passwd:               dbProperty.DbPasswd,
		Addr:                 dbProperty.DbHost,
		Net:                  dbProperty.DbNet,
		ParseTime:            dbProperty.DbParseTime,
		AllowNativePasswords: dbProperty.DbAllowNativePasswords,
		Collation:            dbProperty.DbCollation,
	}
	dataSourceName := c.FormatDSN()
	slog.Info("Database Open", "driver", driverName, "dsn", dataSourceName)

	db, errOpen = sql.Open(driverName, dataSourceName)
	if errOpen != nil {
		return errOpen
	}
	errPing = db.Ping()
	if errPing != nil {
		return errPing
	}

	slog.Info("DB Initialize.", "db", db)
	return nil
}

//---- private function ----
