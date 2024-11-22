// config configシステムパッケージ
package config // パッケージ名はディレクトリ名と同じにする

import (
	"log/slog"

	"sv_lambda_feelcycle_aggregate/fileio"
)

// ---- struct

// configストラクト
type Config struct {
	DbDriver               string
	DbName                 string
	DbUser                 string
	DbPasswd               string
	DbHost                 string
	DbNet                  string
	DbParseTime            bool
	DbAllowNativePasswords bool
	DbCollation            string
}

// ---- Global Variable

// ---- Package Global Variable

var configInformation Config

// ---- public function ----

// ReadConfigInformation (public)config.jsonファイルを読み込み、環境設定値のを取得する関数。
func ReadConfigInformation() error {

	var errRead error
	errRead = fileio.FileIoJsonRead("config.json", &configInformation)
	if errRead != nil {
		return errRead
	}

	slog.Info("Config Read Success.")
	return nil
}

// GetConfigInformation (public)config.jsonファイルの環境設定値を読み込んだ構造体のポインタを取得する関数。
func GetConfigInformation() *Config {
	return &configInformation
}
