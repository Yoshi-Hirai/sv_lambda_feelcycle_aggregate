// db DBシステムパッケージ
package db // パッケージ名はディレクトリ名と同じにする

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
)

// ---- struct

// feelcycleデータベースhisotryテーブルの結果情報構造体(汎用)
type HistoryResult struct {
	Start      string `json:"start"`
	Studio     string `json:"studio"`
	Instructor string `json:"instructor"`
	Program    string `json:"program"`
	Count      int    `json:"count"`
}

// 汎用履歴レスポンス構造体(プログラム履歴、インストラクター履歴)
type MultiHistorResponse struct {
	Searchword string          `json:"searchword"`
	history    []HistoryResult `json:"history"`
}

// ファーストビューレスポンス構造体
type FirstViewResponse struct {
	LastUpdate        string          `json:"lastupdate"`
	InstructorRanking []HistoryResult `json:"instructorranking"`
	ProgramRanking    []HistoryResult `json:"programranking"`
}

// ---- Global Variable

// ---- Package Global Variable

// ---- public function ----
// プログラム履歴
func ProgramHistorySql(program string) ([]byte, error) {

	var sqlStr string
	var rows *sql.Rows
	var errQuery, errScan error
	var i int

	var jsonBytes []byte
	var single HistoryResult
	var result MultiHistorResponse

	// プログラム履歴
	sqlStr = fmt.Sprintf("SELECT DATE_FORMAT(start, '%Y-%m-%d %H:%i'),studio,instructor,program FROM history WHERE program LIKE %s ORDER BY start DESC", program)
	slog.Info(sqlStr)
	rows, errQuery = db.Query(sqlStr)
	if errQuery != nil {
		return jsonBytes, errQuery
	}
	i = 0
	for rows.Next() {
		var start, studio, instructor, program string
		if errScan = rows.Scan(&start, &studio, &instructor, &program); errScan != nil {
			return jsonBytes, errScan
		}
		single = HistoryResult{Start: start, Studio: studio, Instructor: instructor, Program: program, Count: 0}
		result.history = append(result.history, single)
		i++
	}
	result.Searchword = program

	var errMarshal error
	jsonBytes, errMarshal = json.Marshal(result)
	if errMarshal != nil {
		return jsonBytes, errMarshal
	}

	return jsonBytes, nil
}

// インストラクター別受講回数ランキング取得SQL
func FirstViewSql(limit int) ([]byte, error) {

	var sqlStr string
	var rows *sql.Rows
	var errQuery, errScan error
	var i int

	var jsonBytes []byte
	var single HistoryResult
	var result FirstViewResponse

	// 最終更新日
	sqlStr = "SELECT DATE_FORMAT(MAX(start), '%Y-%m-%d %H:%i') FROM history"
	slog.Info(sqlStr)
	rows, errQuery = db.Query(sqlStr)
	if errQuery != nil {
		return jsonBytes, errQuery
	}
	for rows.Next() {
		if errScan = rows.Scan(&result.LastUpdate); errScan != nil {
			return jsonBytes, errScan
		}
	}

	// インストラクターランキング
	sqlStr = fmt.Sprintf("SELECT instructor,COUNT(*) AS Num FROM history GROUP BY instructor ORDER BY Num DESC LIMIT %d", limit)
	slog.Info(sqlStr)
	rows, errQuery = db.Query(sqlStr)
	if errQuery != nil {
		return jsonBytes, errQuery
	}
	i = 0
	for rows.Next() {
		var instructor, count string
		var countInt int
		if errScan = rows.Scan(&instructor, &count); errScan != nil {
			return jsonBytes, errScan
		}
		countInt, _ = strconv.Atoi(count)
		single = HistoryResult{Instructor: instructor, Count: countInt}
		result.InstructorRanking = append(result.InstructorRanking, single)
		i++
	}

	// プログラムランキング
	sqlStr = fmt.Sprintf("SELECT program,COUNT(*) AS Num FROM history GROUP BY program ORDER BY Num DESC LIMIT %d", limit)
	slog.Info(sqlStr)
	rows, errQuery = db.Query(sqlStr)
	if errQuery != nil {
		return jsonBytes, errQuery
	}
	i = 0
	for rows.Next() {
		var program, count string
		var countInt int
		if errScan = rows.Scan(&program, &count); errScan != nil {
			return jsonBytes, errScan
		}
		countInt, _ = strconv.Atoi(count)
		single = HistoryResult{Program: program, Count: countInt}
		result.ProgramRanking = append(result.ProgramRanking, single)
		i++
	}

	var errMarshal error
	jsonBytes, errMarshal = json.Marshal(result)
	if errMarshal != nil {
		return jsonBytes, errMarshal
	}

	return jsonBytes, nil
}

//---- private function ----
