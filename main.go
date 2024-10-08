package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"lambdafunction_base/config"
	"lambdafunction_base/db"
	"lambdafunction_base/log"
)

// ---- struct

// POSTされてくるJSONデータ構造体
type Request struct {
	Action  string `json:"action"`
	Keyword string `json:"keyword"`
}

// ---- Global Variable

// ---- Package Global Variable

// ---- public function ----

// ---- private function
// POSTされてくるJSONデータを構造体に変換する
func convertPostDataToStruct(inputs string) (*Request, error) {
	var req Request
	err := json.Unmarshal([]byte(inputs), &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	jsonLogger := log.GetInstance()
	jsonLogger.Info("main")

	req, errPost := convertPostDataToStruct(request.Body)
	if errPost != nil {
		slog.Error("Convert Post Failed")
		jsonLogger.Error("Convert Post Failed", slog.String("error", errPost.Error()))
		return events.APIGatewayProxyResponse{
			Body:       "NG",
			StatusCode: 500,
		}, nil
	}

	errConfig := config.ReadConfigInformation()
	if errConfig != nil {
		slog.Error("Read Config Failed")
		jsonLogger.Error("Read Config Failed", slog.String("error", errConfig.Error()))
		return events.APIGatewayProxyResponse{
			Body:       "NG",
			StatusCode: 500,
		}, nil
	}

	errDb := db.DbBaseInit()
	if errDb != nil {
		slog.Error("DB Init Failed")
		jsonLogger.Error("DB Init Failed", slog.String("error", errDb.Error()))
		return events.APIGatewayProxyResponse{
			Body:       "NG",
			StatusCode: 500,
		}, nil
	}

	var qResult []byte
	var errQ error
	if req.Action == "Program" {
		// プログラム履歴情報の取得
		qResult, errQ = db.ProgramHistorySql(req.Keyword)
		if errQ != nil {
			slog.Error("ProgramHistorySql Error")
			jsonLogger.Error("ProgramHistorySql Error", slog.String("error", errQ.Error()))
			return events.APIGatewayProxyResponse{
				Body:       "NG",
				StatusCode: 500,
			}, nil
		}
	} else if req.Action == "Instructor" {
		// インストラクター履歴情報の取得
		qResult, errQ = db.InstructorHistorySql(req.Keyword)
		if errQ != nil {
			slog.Error("InstructorHistorySql Error")
			jsonLogger.Error("InstructorHistorySql Error", slog.String("error", errQ.Error()))
			return events.APIGatewayProxyResponse{
				Body:       "NG",
				StatusCode: 500,
			}, nil
		}
	} else {
		// ファーストビュー情報の取得
		var limitnum int = 10
		qResult, errQ = db.FirstViewSql(limitnum)
		if errQ != nil {
			slog.Error("FirstViewSql Error")
			jsonLogger.Error("FirstViewSql Error", slog.String("error", errQ.Error()))
			return events.APIGatewayProxyResponse{
				Body:       "NG",
				StatusCode: 500,
			}, nil
		}
	}

	// 返り値としてレスポンスを返す
	return events.APIGatewayProxyResponse{
		Body: string(qResult),
		Headers: map[string]string{
			"Content-Type": "application/json",
			// CORS対応
			/*
				"Access-Control-Allow-Headers":     "*",                       // CORS対応
				"Access-Control-Allow-Origin":      "http://localhost:5173/",  // CORS対応
				"Access-Control-Allow-Methods":     "GET, POST, PUT, OPTIONS", // CORS対応
				"Access-Control-Allow-Credentials": "true",                    // CORS対応
			*/
			// CORS対応　ここまで
		},
		StatusCode: 200,
	}, nil
}

// ---- main
func main() {
	lambda.Start(handleRequest)
}
