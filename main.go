package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	awsLambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/s3"
	cfg "github.com/jtmilanest/cognito-restore/internal/config"
	"github.com/jtmilanest/cognito-restore/internal/lambda"
	"github.com/jtmilanest/cognito-restore/internal/types"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	S3Bucket         string `json:"bucket"`
	S3BucketFileName string `json:"key"`
}

// Specifies whether the attribute is standard or custom.
type AttributeType struct {

	// The name of the attribute.
	//
	// This member is required.
	Name *string

	// The value of the attribute.
	Value *string
}

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

// Not in use
func Handler(ctx context.Context, event Event) (string, error) {

	event.S3Bucket = "test-cognito-backup001"
	event.S3BucketFileName = "us-west-2_test"

	// Initialize a session that SDK used to load
	// create credentials from shared credentialfile ~/.aws/credentials
	// and configuration from the shared config file ~/.aws/config
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatal("Error creating session:", err)
	}

	// Create Cognito Identity Provider Client
	cip := cognitoidentityprovider.New(sess)

	// Create S3 client
	s3Client := s3.New(sess)

	// Get the user pool data from the S3 bucket
	// Iterate through the user pools
	obj, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(event.S3Bucket),
		Key:    aws.String(event.S3BucketFileName),
	})
	if err != nil {
		log.Fatal("Error listing user pools:", err)
	}
	defer obj.Body.Close()

	// Restore the user pool data
	data, err := io.ReadAll(obj.Body)
	if err != nil {
		log.Errorf("Failed to convert %s object to bytes", event.S3BucketFileName)
	}

	users := &cognitoidentityprovider.ListUsersOutput{}
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Errorf("Failed to unmarshal users backup data. Error: %w", err)
	} else {
		log.Debug("users data has been unmarshalled successfully")
	}

	for _, user := range users.Users {
		fmt.Println("User:", *user.Username)
		var userAttributes []*cognitoidentityprovider.AttributeType
		var userName *string

		for _, attribute := range user.Attributes {
			if *attribute.Name == "email" {
				userName = attribute.Value
			}

			if *attribute.Name != "sub" {
				userAttributes = append(userAttributes, attribute)
			}
		}
		_, err = cip.AdminCreateUser(
			&cognitoidentityprovider.AdminCreateUserInput{
				UserPoolId:     aws.String("us-west-2_6QToEIL3v"),
				Username:       userName,
				UserAttributes: userAttributes,
			},
		)
		if err != nil {
			log.Errorf("Failed to restore users %s. Error: %w", *user.Username, err)
		}
	}
	return "successful", nil
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
	// lambda.Start(RestoreCognitoUserPool)

	// config, err := cfg.NewConfigParam(nil)
	// if err != nil {
	// 	log.Errorf("Lambda execution failed. Error: %s", err)
	// 	os.Exit(1)
	// }

	// err = lambda.Execute(context.TODO(), *config)
	// if err != nil {
	// 	log.Errorf("Lambda has been failed. Error: %s", err)
	// 	os.Exit(1)
	// } else {
	// 	log.Info("Lambda cognito-restore has been completed successfully!")
	// }

	// Execute Lambda function
	log.Info("Starting lambda restore execution ...")
	awsLambda.Start(RestoreCognitoUserPool)
}

/*

Payload to execute Cognito Restore

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
