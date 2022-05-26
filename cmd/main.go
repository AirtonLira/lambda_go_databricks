package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/AirtonLira/lambda_go_databricks/pkg/domain"
	"github.com/aws/aws-lambda-go/lambda"
)

type NotebookDatabricks struct {
	NOTEBOOK_NAME string `json:"notebook_name"`
	PATH          string `json:"path"`
	DOMAIN        string `json:"domain"`
	TOKEN         string `json:"token"`
}

func initDatabricks() (NotebookDatabricks, error) {
	var noteBook NotebookDatabricks
	noteBook.DOMAIN = os.Getenv("DATABRICKS_DOMAIN")
	noteBook.TOKEN = os.Getenv("DATABRICKS_TOKEN")

	return noteBook, nil
}

func callJobDatabricks(ctx context.Context, event domain.S3Event) error {
	noteBook, err := initDatabricks()
	if err != nil {
		log.Fatal(err)
	}

	var bearer = "Bearer " + noteBook.TOKEN

	urlJob := fmt.Sprintf("https://%s/api/2.0/jobs/run-now", noteBook.DOMAIN)

	bodyRequest := strings.NewReader(`{"job_id": 593687828680760}`)

	log.Printf("Requested API: %s", urlJob)

	req, err := http.NewRequest("POST", urlJob, bodyRequest)
	if err != nil {
		log.Fatalf("%v", err)
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
		return err
	}
	log.Printf("Result request databricks API status-code: %v trace: %v", resp.StatusCode, string([]byte(body)))

	return nil
}

func HandleRequest(ctx context.Context, event domain.S3Event) (string, error) {
	fmt.Println("Started function to reprocessing files")
	var err error
	for _, c := range event.Records {
		fmt.Printf("INFO Bucket: %v", c.S3.Object.Key)
		if strings.Contains(c.S3.Object.Key, ".zip") || strings.Contains(c.S3.Object.Key, ".fat") {
			if err = callJobDatabricks(ctx, event); err != nil {
				fmt.Printf("ERROR: %v", err)
				log.Panic(err)
				return "ERROR", err
			}
		} else {
			fmt.Printf("Not exist zip or fat files in this event to reprocessing")
			return "SUCCESS", nil
		}
	}
	fmt.Printf("Finished function reprocessing")
	return "SUCCESS", err
}

func main() {
	lambda.Start(HandleRequest)
}
