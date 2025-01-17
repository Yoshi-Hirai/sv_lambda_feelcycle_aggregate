// db DBシステムパッケージ
package db // パッケージ名はディレクトリ名と同じにする

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// ---- struct

// feelcycleデータベースhisotryテーブルの結果情報構造体(汎用)
type HistoryResult struct {
	Id         int    `json:"id"`
	Start      string `json:"start"`
	Studio     string `json:"studio"`
	Instructor string `json:"instructor"`
	Program    string `json:"program"`
	Count      int    `json:"count"`
}

// feelcycleデータベースhisotryテーブルの集計情報構造体(汎用)
type HistoryTotalling struct {
	Id    int    `json:"id"`
	Item  string `json:"item"`
	Count int    `json:"count"`
	Value int    `json:"value"`
}

// feelcycleデータベースhisotryテーブルの西暦別集計情報構造体(汎用)
type WesternCalenderTotalling struct {
	WesternCalender  string             `json:"westerncalender"`
	TotalInformation []HistoryTotalling `json:"totalinformation"`
}

// 汎用履歴レスポンス構造体(プログラム履歴、インストラクター履歴)
type MultiHistorResponse struct {
	IsGroup    bool            `json:"isgroup"`
	Searchword string          `json:"searchword"`
	History    []HistoryResult `json:"history"`
}

// ファーストビューレスポンス構造体
type FirstViewResponse struct {
	LastUpdate                 string                     `json:"lastupdate"`
	InstructorRanking          []HistoryResult            `json:"instructorranking"`
	ProgramRanking             []HistoryResult            `json:"programranking"`
	ProgramCategoryTotalling   []HistoryTotalling         `json:"programcategorytotalling"`
	WesternInstructorTotalling []WesternCalenderTotalling `json:"westerninstructortotalling"`
	WesternProgramTotalling    []WesternCalenderTotalling `json:"westernprogramtotalling"`
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
	sqlStr = fmt.Sprintf("SELECT DATE_FORMAT(start, '%%Y-%%m-%%d %%H:%%i'),studio,instructor,program FROM history WHERE program LIKE \"%s\" ORDER BY start DESC", program)
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
		single = HistoryResult{Id: i, Start: start, Studio: studio, Instructor: instructor, Program: program, Count: 0}
		result.History = append(result.History, single)
		i++
	}
	result.IsGroup = false
	result.Searchword = program

	var errMarshal error
	jsonBytes, errMarshal = json.Marshal(result)
	if errMarshal != nil {
		return jsonBytes, errMarshal
	}

	return jsonBytes, nil
}

// プログラム履歴(インストラクター集約)
func ProgramHistoryGroupInstructorInstructorSql(program string) ([]byte, error) {

	var sqlStr string
	var rows *sql.Rows
	var errQuery, errScan error
	var i int

	var jsonBytes []byte
	var single HistoryResult
	var result MultiHistorResponse

	// プログラム履歴
	sqlStr = fmt.Sprintf("SELECT instructor,program,COUNT(*) FROM history WHERE program LIKE \"%s\" GROUP BY instructor", program)
	slog.Info(sqlStr)
	rows, errQuery = db.Query(sqlStr)
	if errQuery != nil {
		return jsonBytes, errQuery
	}
	i = 0
	for rows.Next() {
		var instructor, program, count string
		if errScan = rows.Scan(&instructor, &program, &count); errScan != nil {
			return jsonBytes, errScan
		}
		countInt, _ := strconv.Atoi(count)
		single = HistoryResult{Id: i, Start: "", Studio: "", Instructor: instructor, Program: program, Count: countInt}
		result.History = append(result.History, single)
		i++
	}
	result.IsGroup = true
	result.Searchword = program

	var errMarshal error
	jsonBytes, errMarshal = json.Marshal(result)
	if errMarshal != nil {
		return jsonBytes, errMarshal
	}

	return jsonBytes, nil
}

// インストラクター履歴
func InstructorHistorySql(instructor string) ([]byte, error) {

	var sqlStr string
	var rows *sql.Rows
	var errQuery, errScan error
	var i int

	var jsonBytes []byte
	var single HistoryResult
	var result MultiHistorResponse

	// インストラクター履歴
	sqlStr = fmt.Sprintf("SELECT DATE_FORMAT(start, '%%Y-%%m-%%d %%H:%%i'),studio,instructor,program FROM history WHERE instructor LIKE \"%s\" ORDER BY start DESC", instructor)
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
		single = HistoryResult{Id: i, Start: start, Studio: studio, Instructor: instructor, Program: program, Count: 0}
		result.History = append(result.History, single)
		i++
	}
	result.IsGroup = false
	result.Searchword = instructor

	var errMarshal error
	jsonBytes, errMarshal = json.Marshal(result)
	if errMarshal != nil {
		return jsonBytes, errMarshal
	}

	return jsonBytes, nil
}

// インストラクター履歴(プログラム集約)
func InstructorHistoryGroupProgramSql(instructor string) ([]byte, error) {

	var sqlStr string
	var rows *sql.Rows
	var errQuery, errScan error
	var i int

	var jsonBytes []byte
	var single HistoryResult
	var result MultiHistorResponse

	// インストラクター履歴
	sqlStr = fmt.Sprintf("SELECT instructor,program,COUNT(*) FROM history WHERE instructor LIKE \"%s\" GROUP BY program", instructor)
	slog.Info(sqlStr)
	rows, errQuery = db.Query(sqlStr)
	if errQuery != nil {
		return jsonBytes, errQuery
	}
	i = 0
	for rows.Next() {
		var instructor, program, count string
		if errScan = rows.Scan(&instructor, &program, &count); errScan != nil {
			return jsonBytes, errScan
		}
		countInt, _ := strconv.Atoi(count)
		single = HistoryResult{Id: i, Start: "", Studio: "", Instructor: instructor, Program: program, Count: countInt}
		result.History = append(result.History, single)
		i++
	}
	result.IsGroup = true
	result.Searchword = instructor

	var errMarshal error
	jsonBytes, errMarshal = json.Marshal(result)
	if errMarshal != nil {
		return jsonBytes, errMarshal
	}

	return jsonBytes, nil
}

// ファーストビュー
// インストラクター別プログラム別受講回数ランキング取得SQL
// 西暦別インストラクター受講回数ランキング
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
		single = HistoryResult{Id: i, Instructor: instructor, Count: countInt}
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
		single = HistoryResult{Id: i, Program: program, Count: countInt}
		result.ProgramRanking = append(result.ProgramRanking, single)
		i++
	}

	result.ProgramCategoryTotalling, errQuery = programCategoryTotallingSql()
	if errQuery != nil {
		return jsonBytes, errQuery
	}

	result.WesternInstructorTotalling, result.WesternProgramTotalling, errQuery = instructorWesternCalenderTotallingSql()
	if errQuery != nil {
		return jsonBytes, errQuery
	}

	var errMarshal error
	jsonBytes, errMarshal = json.Marshal(result)
	if errMarshal != nil {
		return jsonBytes, errMarshal
	}

	return jsonBytes, nil
}

// ---- private function ----
// プログラムカテゴリ毎の集計クエリ
func programCategoryTotallingSql() ([]HistoryTotalling, error) {
	var programCategory []string = []string{
		"",
		"BB3 ",
		"BSBi ",
		"BSWi ",
		"BB2 ",
		"BSL ",
		"BSB ",
		"BSW ",
		"BB1 ",
	}

	var result []HistoryTotalling
	var num, totalnum, value int

	i := 0
	for _, v := range programCategory {

		sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM history WHERE program LIKE \"%s%%\"", v)
		slog.Info(sqlStr)
		rows, errQuery := db.Query(sqlStr)
		if errQuery != nil {
			return result, errQuery
		}
		var single HistoryTotalling
		for rows.Next() {
			if errScan := rows.Scan(&single.Count); errScan != nil {
				return result, errScan
			}
		}
		if v == "" {
			totalnum = single.Count
			single.Value = int((float64(single.Count) / float64(totalnum)) * 100)
		} else {
			num = num + single.Count
			if totalnum > 0 {
				single.Value = int((float64(single.Count) / float64(totalnum)) * 100)
				value = value + single.Value
			}
		}
		single.Id = i
		single.Item = strings.TrimSpace(v)
		result = append(result, single)
		i++
	}
	lengthCategory := len(programCategory)
	otherCount := totalnum - num
	otherValue := 100 - value
	otherdata := HistoryTotalling{Id: lengthCategory, Item: "others", Count: otherCount, Value: otherValue}
	result = append(result, otherdata)
	return result, nil
}

// 西暦別インストラクタ毎の集計クエリ
func instructorWesternCalenderTotallingSql() ([]WesternCalenderTotalling, []WesternCalenderTotalling, error) {

	var westernCalender []string = []string{
		"2016",
		"2017",
		"2018",
		"2019",
		"2020",
		"2021",
		"2022",
		"2023",
		"2024",
		"2025",
	}
	var resultInstructor, resultProgram []WesternCalenderTotalling

	for _, v := range westernCalender {

		var sqlStr string
		var idx int
		var totalI, totalP []HistoryTotalling
		var rows *sql.Rows
		var errQuery error

		// インストラクター
		sqlStr = fmt.Sprintf("SELECT instructor,COUNT(*) AS num FROM history WHERE DATE_FORMAT(start, '%%Y')=\"%s\" GROUP BY instructor ORDER BY num DESC", v)
		slog.Info(sqlStr)
		rows, errQuery = db.Query(sqlStr)
		if errQuery != nil {
			return resultInstructor, resultProgram, errQuery
		}
		idx = 0
		for rows.Next() {
			var item string
			var count int
			if errScan := rows.Scan(&item, &count); errScan != nil {
				return resultInstructor, resultProgram, errScan
			}
			singleResult := HistoryTotalling{Id: idx, Item: item, Count: count}
			totalI = append(totalI, singleResult)
			idx++
		}
		singleI := WesternCalenderTotalling{WesternCalender: v, TotalInformation: totalI}
		resultInstructor = append(resultInstructor, singleI)

		// プログラム
		sqlStr = fmt.Sprintf("SELECT program,COUNT(*) AS num FROM history WHERE DATE_FORMAT(start, '%%Y')=\"%s\" GROUP BY program ORDER BY num DESC", v)
		slog.Info(sqlStr)
		rows, errQuery = db.Query(sqlStr)
		if errQuery != nil {
			return resultInstructor, resultProgram, errQuery
		}
		idx = 0
		for rows.Next() {
			var item string
			var count int
			if errScan := rows.Scan(&item, &count); errScan != nil {
				return resultInstructor, resultProgram, errScan
			}
			singleResult := HistoryTotalling{Id: idx, Item: item, Count: count}
			totalP = append(totalP, singleResult)
			idx++
		}
		singleP := WesternCalenderTotalling{WesternCalender: v, TotalInformation: totalP}
		resultProgram = append(resultProgram, singleP)
	}
	return resultInstructor, resultProgram, nil
}
