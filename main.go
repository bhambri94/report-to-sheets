package main

import (
	"fmt"
	"log"

	"github.com/bhambri94/report-to-sheets/configs"
	"github.com/bhambri94/report-to-sheets/email"
	"github.com/bhambri94/report-to-sheets/googleSheets"
	"github.com/bhambri94/report-to-sheets/report"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var (
	logger, _          = zap.NewProduction()
	sugar              = logger.Sugar()
	TestCaseResultsMap = make(map[string]string)
)

func main() {
	configs.SetConfig()
	sugar.Infof("starting report-to-sheets app server...")
	defer logger.Sync() // flushes buffer, if any

	router := fasthttprouter.New()
	router.POST("/v1/report/save", handleReportToSheets)
	router.POST("/v1/report/save/SendEmail=:SendEmail/RepeatEmailWithSameResults=:RepeatEmailWithSameResults", handleReportToSheets)
	router.GET("/v1/report/sendEmail", handleReportToEmail)
	log.Fatal(fasthttp.ListenAndServe(":3002", router.Handler))
}
func handleReportToEmail(ctx *fasthttp.RequestCtx) {
	// email.SendEmail()
}

func handleReportToSheets(ctx *fasthttp.RequestCtx) {
	sugar.Infof("received a push report request to Google Sheets!")
	SendEmailFromRequest := ctx.UserValue("SendEmail")
	RepeatEmailWithSameResultsRequest := ctx.UserValue("RepeatEmailWithSameResults")
	if SendEmailFromRequest == nil {
		SendEmailFromRequest = "false"
	}
	if RepeatEmailWithSameResultsRequest == nil {
		RepeatEmailWithSameResultsRequest = "true"
	}

	fh, err := ctx.FormFile("file")
	if err != nil {
		sugar.Error(err)
		ctx.Response.SetStatusCode(500)
		successResponse := "{\"success\":false,\"response\":\"File key mentioned in request body is wrong\"}"
		ctx.Write([]byte(successResponse))
	}
	JsonFileName := fh.Filename
	if err := fasthttp.SaveMultipartFile(fh, "uploads/latestreport.json"); err != nil {
		sugar.Error(err)
		ctx.Response.SetStatusCode(500)
		successResponse := "{\"success\":false,\"response\":\"Unable to save request body file\"}"
		ctx.Write([]byte(successResponse))
	}
	finalValues, SendEmail := report.GetReport(JsonFileName, RepeatEmailWithSameResultsRequest.(string))
	if SendEmail && SendEmailFromRequest.(string) == "true" {
		email.SendEmail(finalValues)
	} else {
		fmt.Println("No Email Sent")
	}
	if len(finalValues) > 0 {
		googleSheets.BatchAppend(configs.Configurations.SheetNameWithRange, finalValues)
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	if SendEmail && SendEmailFromRequest.(string) == "true" {
		ctx.SetBody([]byte("{\"success\":true,\"response\":\"Sheet has been updated, Email report has also been sent\"}"))
	} else {
		ctx.SetBody([]byte("{\"success\":true,\"response\":\"Sheet has been updated\"}"))
	}
}
