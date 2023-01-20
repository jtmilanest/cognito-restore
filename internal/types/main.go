package types

import "github.com/guregu/null"

// Response struct
type Response struct {
	Message string `json:"answer"`
}

// Event struct
type Event struct {
	AWSRegion string `json:"awsRegion"`

	CognitoUserPoolID string `json:"cognitoUserPoolID"`
	CognitoRegion     string `json:"cognitoRegion"`

	S3BucketName   string `json:"s3BucketName"`
	S3BucketRegion string `json:"s3BucketRegion"`

	BackupDirPath string `json:"backupDirPath"`

	RestoreUsers         null.Bool `json:"restoreUsers"`
	RestoreGroups        null.Bool `json:"restoreGroups"`
	CleanUpBeforeRestore null.Bool `json:"cleanUpBeforeRestore"`
}
