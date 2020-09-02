package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func SendEmail(EmailContent [][]interface{}) {
	var HTMLReportName string
	i := 0
	var finalRow string
	for i < len(EmailContent) {
		HTMLReportName = EmailContent[i][0].(string)
		row := "Feature => " + EmailContent[i][1].(string) + " || Scenarios => "
		finalRow = finalRow + row
		j := 2
		for j < len(EmailContent[i]) {
			row := EmailContent[i][j].(string) + " -- " + strings.ToUpper(EmailContent[i][j+1].(string)) + " , "
			finalRow = finalRow + row
			j = j + 2
		}
		i++
		finalRow = finalRow + "\r"
	}

	from := "oneplanautomationreports@gmail.com"
	password := "MySecretPassword"
	to := []string{
		"shivambhambri94@gmail.com",
		"steven@onplan.co",
	}
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	message := []byte("Subject: OnPlan Automation Results  \r\n\r\n" + "OnPlan Automation Report Details:\r\r" + "HTML Report can be found with name " + HTMLReportName + "\r\r" + "Please find below all test case results: \r" + finalRow + "\r" + "Thanks")
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
