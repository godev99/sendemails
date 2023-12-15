package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/karim-w/go-azure-communication-services/emails"
)

var (
	ctx    = context.Background()
	appId  *string
	tenant *string
	subId  *string
)

var host = "zagocbscdv01ew1cs01.france.communication.azure.com"                                      // os.Getenv("ACS_HOST")
var key = "/On9q6QNWkbAg5KATebndcMy/IXuyO1qeTAEhMs3UfYtpgJVF2G9fEpLpbySTFbv0WnXzw5SfHgFPq6tSkeAAg==" // os.Getenv("ACS_KEY")
// recipient := "givanes@pm.me"                                                                      // os.Getenv("EMAIL_RECIPIENT")
var sender = "DoNotReply@d5f7e57b-6ade-4383-abac-2bf51413c4fd.azurecomm.net" // os.Getenv("EMAIL_SENDER")

func sendmail(c *gin.Context) {

	// define custom type
	type Input struct {
		Subject   string `json:"subject" binding:"required"`
		Body      string `json:"body" binding:"required"`
		Recipient string `json:"recipient" binding:"required"`
	}

	var input Input

	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	client := emails.NewClient(host, key, nil)
	payload := emails.Payload{
		Headers: emails.Headers{
			ClientCorrelationID:    "1234",
			ClientCustomHeaderName: "ClientCustomHeaderValue",
		},
		SenderAddress: sender,
		Content: emails.Content{
			Subject:   input.Subject,
			PlainText: input.Body,
		},
		Recipients: emails.Recipients{
			To: []emails.ReplyTo{
				{
					Address: input.Recipient,
				},
			},
		},
	}

	result, err := client.SendEmail(context.TODO(), payload)
	if err != nil {
		c.String(http.StatusNotFound, "email not send")
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

func liveness(c *gin.Context) {
	result := "api is live!"
	c.IndentedJSON(http.StatusOK, result)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func get_port() string {
	port := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = ":" + val
	}
	return port
}

func main() {
	r := gin.Default()
	r.POST("/api/sendmail", sendmail)
	r.GET("/api/liveness", liveness)
	r.Run(get_port())
}
