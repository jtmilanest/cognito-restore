package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	awsLambda 
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

func RestoreCognitoUserPool(ctx context.Context, event Event) (string, error) {

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

func main() {
	// Execute Lambda function
	// lambda.Start(RestoreCognitoUserPool)

	log.Info("Starting lambda restore execution ...")

}
