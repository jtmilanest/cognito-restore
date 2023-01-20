package main

import (
	"context"
	"os"

	awsLambda "github.com/aws/aws-lambda-go/lambda"
	cfg "github.com/jtmilanest/cognito-restore/internal/config"
	"github.com/jtmilanest/cognito-restore/internal/lambda"
	"github.com/jtmilanest/cognito-restore/internal/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(false)

	var formatter log.Formatter

	if formatterType, ok := os.LookupEnv("FORMATTER_TYPE"); ok {
		if formatterType == "JSON" {
			formatter = &log.JSONFormatter{PrettyPrint: false}
		}

		if formatterType == "TEXT" {
			formatter = &log.TextFormatter{DisableColors: false}
		}
	}

	if formatter == nil {
		formatter = &log.TextFormatter{DisableColors: false}
	}

	log.SetFormatter(formatter)

	var logLevel log.Level
	var err error

	if ll, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel, err = log.ParseLevel(ll)
		if err != nil {
			logLevel = log.DebugLevel
		}
	} else {
		logLevel = log.DebugLevel
	}

	log.SetLevel(logLevel)
}

// Function handler to execute lambda code to AWS
func RestoreCognitoUserPool(ctx context.Context, event types.Event) (types.Response, error) {
	log.Infof("Handling lambda for event: %v", event)
	// Instantiate new config param
	config, err := cfg.NewConfigParam(event)
	if err != nil {
		return types.Response{Message: "Lambda has been failed"}, err
	}

	// Execute Lambda with instantiated new config param
	var msg string
	err = lambda.Execute(ctx, *config)
	if err != nil {
		msg = "Lambda has been failed."
	} else {
		msg = "Lambda has been completed successfuly!"
	}

	return types.Response{Message: msg}, err
}

func main() {
	// Execute Lambda function
	log.Info("Starting lambda restore execution ...")
	awsLambda.Start(RestoreCognitoUserPool)
}

/*

Sample Payload to execute the Cognito Restore using
Lambda UI in AWS

{
  "awsRegion": "us-west-2",
  "cognitoUserPoolId": "us-west-2_Xy67PstDj",
  "cognitoRegion": "us-west-2",
  "s3BucketName": "test-cognito-backup001",
  "s3BucketRegion": "us-west-2",
  "backupDirPath": "platform",
  "restoreUsers": true,
  "restoreGroups": false, //todo
  "cleanUpBeforeRestore": true
}

*/
