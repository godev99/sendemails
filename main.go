package main

import (
	"context"
	"fmt"
	"github.com/karim-w/go-azure-communication-services/emails"
	"log"
	"net/http"
	"os"
)

var host = "zagocbscdv01ew1cs01.france.communication.azure.com"                                      // os.Getenv("ACS_HOST")
var key = "/On9q6QNWkbAg5KATebndcMy/IXuyO1qeTAEhMs3UfYtpgJVF2G9fEpLpbySTFbv0WnXzw5SfHgFPq6tSkeAAg==" // os.Getenv("ACS_KEY")
// recipient := "givanes@pm.me"                                                                      // os.Getenv("EMAIL_RECIPIENT")
var sender = "DoNotReply@d5f7e57b-6ade-4383-abac-2bf51413c4fd.azurecomm.net" // os.Getenv("EMAIL_SENDER")

func sendEmail(w http.ResponseWriter, r *http.Request) {
	message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"

	subject := r.URL.Query().Get("subject")
	body := r.URL.Query().Get("body")
	recipient := r.URL.Query().Get("recipient")
	name := r.URL.Query().Get("name")

	if name != "" {
		message = fmt.Sprintf("Hello, %s. This HTTP triggered function executed successfully.\n", name)
	}

	client := emails.NewClient(host, key, nil)
	payload := emails.Payload{
		Headers: emails.Headers{
			ClientCorrelationID:    "1234",
			ClientCustomHeaderName: "ClientCustomHeaderValue",
		},
		SenderAddress: sender,
		Content: emails.Content{
			Subject:   subject,
			PlainText: body,
		},
		Recipients: emails.Recipients{
			To: []emails.ReplyTo{
				{
					Address: recipient,
				},
			},
		},
	}

	result, err := client.SendEmail(context.TODO(), payload)
	if err != nil {
		fmt.Println(err)
	}
	if result.ID == "" {
		fmt.Println("TrackingId is empty")
	}
	if result.Status == "" {
		fmt.Println("Status is empty")
	}

	fmt.Fprint(w, message)
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/sendemail", sendEmail)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
