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
	Action string `json:"action"`
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

	slog.Info("Request", "body", request.Body)
	req, errPost := convertPostDataToStruct(request.Body)
	if errPost != nil {
		slog.Error("Convert Post Failed")
		jsonLogger.Error("Convert Post Failed", slog.String("error", errPost.Error()))
		return events.APIGatewayProxyResponse{
			Body:       "NG",
			StatusCode: 500,
		}, nil
	}
	slog.Info("Post", "Action", req.Action)

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

	// ファーストビュー情報の取得
	qResult0, errQ0 := db.FirstViewSql(10)
	if errQ0 != nil {
		slog.Error("Query Error")
		jsonLogger.Error("Query Error", slog.String("error", errQ0.Error()))
		return events.APIGatewayProxyResponse{
			Body:       "NG",
			StatusCode: 500,
		}, nil
		//		return "NG", dbErr
	}

	// 返り値としてレスポンスを返す
	return events.APIGatewayProxyResponse{
		Body: string(qResult0),
		Headers: map[string]string{
			"Content-Type": "application/json",
			// CORS対応(仮)
			"Access-Control-Allow-Headers":     "*",                       // CORS対応
			"Access-Control-Allow-Origin":      "http://localhost:5173",   // CORS対応
			"Access-Control-Allow-Methods":     "GET, POST, PUT, OPTIONS", // CORS対応
			"Access-Control-Allow-Credentials": "true",                    // CORS対応
			// CORS対応(仮)　ここまで
		},
		StatusCode: 200,
	}, nil
}

// ---- main
func main() {
	lambda.Start(handleRequest)
}
