package config

import (
	"os"

	"config/types"

	"github.com/guregu/null"
)

// Config parameter when launching lambda
/*
For example e.g.,
{
  "awsRegion": "us-west-2",
  "cognitoUserPoolId": "us-west-2_EP1dk34",
  "s3BucketName": "mycognitotest",
  "backupDirPath": "2023-01-19T9:00:00Z/",
  "restoreUsers": true,
  "restoreGroups": true,
  "cleanUpBeforeRestore": true
}
*/

type ConfigParam struct {
	AWSRegion string

	CognitoUserPoolID string
	CognitoRegion     string

	S3BucketName   string
	S3BucketRegion string

	BackupDirPath string

	RestoreUsers         null.Bool
	RestoreGroups        null.Bool
	CleanUpBeforeRestore null.Bool
}

// Helper function to verify Environment Variables
func getLookEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewConfigParam(eventRaw interface{}) (*ConfigParam, error) {
	var config = &ConfigParam{}
	var getFromEvent bool
	var event types.Event

	// Process AWS Region
	if awsRegion := getLookEnv("AWS_REGION", ""); awsRegion != "" {

	}

}
