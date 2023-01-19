package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/guregu/null"
	"github.com/jtmilanest/cognito-restore/internal/types"
	log "github.com/sirupsen/logrus"
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

	// Switch between parameter or environment variables
	switch value := eventRaw.(type) {
	case types.Event:
		getFromEvent = true
		event = value
	default:
		getFromEvent = false
	}

	// Process AWS Region
	if awsRegion := getLookEnv("AWS_REGION", ""); awsRegion != "" {
		config.AWSRegion = awsRegion
	} else {
		log.Warn("Environment variable for AWS_REGION is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.AWSRegion != "" {
			config.AWSRegion = event.AWSRegion
		} else {
			log.Warn("Event contains empty awsRegion variable")
		}
	}
	if config.AWSRegion == "" {
		return nil, fmt.Errorf("awsRegion is empty;Configure it via 'AWS_REGION' env variable OR pass in event body")
	}
	// Process AWS Region

	// Process Cognito User Pool ID
	if cognitoUserPoolID := getLookEnv("COGNITO_USER_POOL_ID", ""); cognitoUserPoolID != "" {
		config.CognitoUserPoolID = cognitoUserPoolID
	} else {
		log.Warn("Environment variable for COGNITO_USER_POOL_ID is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.CognitoUserPoolID != "" {
			config.CognitoUserPoolID = event.CognitoUserPoolID
		} else {
			log.Warn("Event contains empty cognitoUserPoolID variable")
		}
	}
	if config.CognitoUserPoolID == "" {
		return nil, fmt.Errorf("cognitoUserPoolID is empty;Configure it via 'COGNITO_USER_POOL_ID' env variable OR pass in event body")
	}
	// Process Cognito User Pool ID

	// Process Cognito Region
	if cognitoRegion := getLookEnv("COGNITO_REGION", ""); cognitoRegion != "" {
		config.CognitoRegion = cognitoRegion
	} else {
		log.Warn("Environment variable for COGNITO_REGION is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.CognitoRegion != "" {
			config.CognitoRegion = event.CognitoRegion
		} else {
			log.Warn("Event contains empty cognitoRegion variable")
		}
	}
	if config.CognitoRegion == "" {
		return nil, fmt.Errorf("cognitoRegion is empty;Configure it via 'COGNITO_REGION' env variable OR pass in event body")
	}
	// Process Cognito Region

	// Process S3BucketName
	if s3BucketName := getLookEnv("S3_BUCKET_NAME", ""); s3BucketName != "" {
		config.S3BucketName = s3BucketName
	} else {
		log.Warn("Environment variable for S3_BUCKET_NAME is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.S3BucketName != "" {
			config.S3BucketName = event.S3BucketName
		} else {
			log.Warn("Event contains empty s3BucketName variable")
		}
	}
	if config.S3BucketName == "" {
		return nil, fmt.Errorf("s3BucketName is empty;Configure it via 'S3_BUCKET_NAME' env variable OR pass in event body")
	}
	// Process S3BucketName

	// Process S3BucketRegion
	if s3BucketRegion := getLookEnv("S3_BUCKET_REGION", ""); s3BucketRegion != "" {
		config.S3BucketRegion = s3BucketRegion
	} else {
		log.Warn("Environment variable for S3_BUCKET_REGION is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.S3BucketRegion != "" {
			config.S3BucketRegion = event.S3BucketRegion
		} else {
			log.Warn("Event contains empty s3BucketRegion variable")
		}
	}
	if config.S3BucketRegion == "" {
		return nil, fmt.Errorf("s3BucketRegion is empty;Configure it via 'S3_BUCKET_REGION' env variable OR pass in event body")
	}
	// Process S3BucketRegion

	// Process BackupDirPath
	if backupDirPath := getLookEnv("BACKUP_DIR_PATH", ""); backupDirPath != "" {
		config.BackupDirPath = backupDirPath
	} else {
		log.Warn("Environment variable for BACKUP_DIR_PATH is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.BackupDirPath != "" {
			config.BackupDirPath = event.BackupDirPath
		} else {
			log.Warn("Event contains empty backupDirPath variable")
		}
	}
	if config.BackupDirPath == "" {
		return nil, fmt.Errorf("backupDirPath is empty;Configure it via 'BACKUP_DIR_PATH' env variable OR pass in event body")
	}
	// Process BackupDirPath

	// Process RestoreUsers
	if restoreUsers := getLookEnv("RESTORE_USERS", ""); restoreUsers != "" {
		restoreUsersValue, err := strconv.ParseBool(restoreUsers)
		if err != nil {
			return nil, fmt.Errorf("Could not parse `RESTURE_USERS` variable. Error: %w", err)
		}
		config.RestoreUsers = null.NewBool(restoreUsersValue, true)
	} else {
		log.Warn("Environment variable `RESTORE_USERS` is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.RestoreUsers.Valid {
			config.RestoreUsers = event.RestoreUsers
		}
	}
	if !config.RestoreUsers.Valid {
		log.Warn("restoreUsers is not specified, Default value 'false' will be used")
		config.RestoreUsers = null.NewBool(false, true)
	}
	// Process RestoreUsers

	// Process RestoreGroups
	if restoreGroups := getLookEnv("RESTORE_GROUPS", ""); restoreGroups != "" {
		restoreGroupsValue, err := strconv.ParseBool(restoreGroups)
		if err != nil {
			return nil, fmt.Errorf("Could not parse `RESTURE_USERS` variable. Error: %w", err)
		}
		config.RestoreGroups = null.NewBool(restoreGroupsValue, true)
	} else {
		log.Warn("Environment variable `RESTORE_GROUPS` is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.RestoreGroups.Valid {
			config.RestoreGroups = event.RestoreGroups
		}
	}
	if !config.RestoreGroups.Valid {
		log.Warn("restoreGroups is not specified, Default value 'false' will be used")
		config.RestoreGroups = null.NewBool(false, true)
	}
	// Process RestoreGroups

	// Process cleanUpBeforeRestore
	if cleanUpBeforeRestore := getLookEnv("CLEANUP_BEFORE_RESTORE", ""); cleanUpBeforeRestore != "" {
		cleanUpBeforeRestoreValue, err := strconv.ParseBool(cleanUpBeforeRestore)
		if err != nil {
			return nil, fmt.Errorf("Could not parse 'CLEANUP_BEFORE_RESTORE' variable. Error: %w", err)
		}
		config.CleanUpBeforeRestore = null.NewBool(cleanUpBeforeRestoreValue, true)
	} else {
		log.Warn("Environment variable 'CLEANUP_BEFORE_RESTORE' is empty")
	}

	// pass the value to config
	if getFromEvent {
		if event.CleanUpBeforeRestore.Valid {
			log.Warn("cleanupBeforeRestore is not specified; Default value 'false' will be used")
			config.CleanUpBeforeRestore = null.NewBool(false, true)
		} else {
			if config.CleanUpBeforeRestore.Bool {
				log.Warnf("Pay attention that CLEANUP_BEFORE_RESTORE is 'true'. It means all data from %s userpool will be deleted before restore", config.CognitoUserPoolID)
			}
		}
	}
	// Process cleanUpBeforeRestore

	// Return config, no errors
	return config, nil
}
