package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/karim-w/go-azure-communication-services/emails"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

var products map[int]Product = map[int]Product{
	1: {"Milk"},
	2: {"Butter"},
}

type Product struct {
	Name string
}

func init() {
	flag.StringVar(&instrumentationKey, "instrumentationKey", os.Getenv("INSTRUMENTATION_KEY"), "set instrumentation key from azure portal")
	telemetryClient = appinsights.NewTelemetryClient(instrumentationKey)
	/*Set role instance name globally -- this is usually the name of the service submitting the telemetry*/
	telemetryClient.Context().Tags.Cloud().SetRole("aifapi")
	/*turn on diagnostics to help troubleshoot problems with telemetry submission. */
	appinsights.NewDiagnosticsMessageListener(func(msg string) error {
		log.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
		return nil
	})
}

var (
	telemetryClient    appinsights.TelemetryClient
	instrumentationKey string
)
var host = "zagocbscdv01ew1cs01.france.communication.azure.com"                                      // os.Getenv("ACS_HOST")
var key = "/On9q6QNWkbAg5KATebndcMy/IXuyO1qeTAEhMs3UfYtpgJVF2G9fEpLpbySTFbv0WnXzw5SfHgFPq6tSkeAAg==" // os.Getenv("ACS_KEY")
// recipient := "givanes@pm.me"                                                                      // os.Getenv("EMAIL_RECIPIENT")
var sender = "DoNotReply@d5f7e57b-6ade-4383-abac-2bf51413c4fd.azurecomm.net" // os.Getenv("EMAIL_SENDER")

func getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid product identifier")
		return
	}
	p, ok := products[id]
	if !ok {
		c.String(http.StatusNotFound, "Product not found")
		return
	}
	c.IndentedJSON(http.StatusOK, p)

}

func get_port() string {
	port := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = ":" + val
	}
	return port
}

func sendmail(c *gin.Context) {
	fmt.Println("send email")
	trace := appinsights.NewTraceTelemetry("begin create application", appinsights.Information)
	telemetryClient.Track(trace)

	// define custom type
	type Input struct {
		Subject   string `json:"subject" binding:"required"`
		Body      string `json:"body" binding:"required"`
		Recipient string `json:"recipient" binding:"required"`
	}

	var input Input

	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		trace := appinsights.NewTraceTelemetry("failed to bind request params to go struct for create application", appinsights.Critical)
		telemetryClient.Track(trace)
		telemetryClient.TrackException(err)
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
	} else {
		trace := appinsights.NewTraceTelemetry("email send", appinsights.Information)
		telemetryClient.Track(trace)
	}

	c.IndentedJSON(http.StatusOK, result)
}

func sendmail2(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func main() {
	r := gin.Default()
	r.GET("/api/products", getProducts)
	r.GET("/api/products/:id", getProduct)
	r.POST("/api/sendmail", sendmail)
	r.GET("/api/sendmail", sendmail2)
	r.Run(get_port())
}
