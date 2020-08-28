package main

import (
	"log"

	"github.com/bhambri94/report-to-sheets/configs"
	"github.com/bhambri94/report-to-sheets/googleSheets"
	"github.com/bhambri94/report-to-sheets/report"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
)

func main() {
	configs.SetConfig()
	sugar.Infof("starting report-to-sheets app server...")
	defer logger.Sync() // flushes buffer, if any

	router := fasthttprouter.New()
	router.POST("/v1/report/save", handleReportToSheets)
	log.Fatal(fasthttp.ListenAndServe(":8010", router.Handler))
}

func handleReportToSheets(ctx *fasthttp.RequestCtx) {
	sugar.Infof("received a push report request to Google Sheets!")
	fh, err := ctx.FormFile("file")
	JsonFileName := fh.Filename
	if err != nil {
		sugar.Error(err)
		ctx.Response.SetStatusCode(500)
		successResponse := "{\"success\":false,\"response\":\"File key mentioned in request body is wrong\"}"
		ctx.Write([]byte(successResponse))
	}
	if err := fasthttp.SaveMultipartFile(fh, "uploads/latestreport.json"); err != nil {
		sugar.Error(err)
		ctx.Response.SetStatusCode(500)
		successResponse := "{\"success\":false,\"response\":\"Unable to save request body file\"}"
		ctx.Write([]byte(successResponse))
	}
	finalValues := report.GetReport(JsonFileName)
	if len(finalValues) > 0 {
		googleSheets.BatchAppend(configs.Configurations.SheetNameWithRange, finalValues)
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetBody([]byte("{\"success\":true,\"response\":\"Sheet has been updated\"}"))
	// sugar.Infof(string(ctx.Request.Body()))
}
