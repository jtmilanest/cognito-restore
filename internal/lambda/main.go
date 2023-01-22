package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jtmilanest/cognito-restore/internal/cloud"
	"github.com/jtmilanest/cognito-restore/internal/config"
	log "github.com/sirupsen/logrus"
)

// Decrypt S3 Data via KMS
func decryptViaKMS(ctx context.Context, client *cloud.Client, kmsKeyName string, data []byte) (*kms.DecryptOutput, error) {
	result, err := client.KMSClient.Decrypt(ctx, &kms.DecryptInput{
		KeyId:          aws.String(kmsKeyName),
		CiphertextBlob: data,
	})

	return result, err
}

// Retreive Data to restore from S3
func getDataFromS3(ctx context.Context, client *cloud.Client, bucketName, keyName string) ([]byte, error) {

	// Retrieve S3 BucketName and Filename of Cognito User Pools
	obj, err := client.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	if err != nil {
		log.Errorf("Failed to get %s object data from %s bucket", keyName, bucketName)

		return nil, err
	}

	// Return slice of bytes which is the data and error
	data, err := io.ReadAll(obj.Body)
	if err != nil {
		log.Errorf("Failed to convert %s object data to bytes", keyName)

		return nil, err
	}

	return data, err
}

// Execute Lambda Function
func Execute(ctx context.Context, config config.ConfigParam) error {

	client, err := cloud.New(ctx, config.CognitoRegion, config.S3BucketRegion, config.KMSRegion)
	if err != nil {
		return fmt.Errorf("Could not create AWS client. Error %w", err)
	}

	// Fresh cleanup of users in pool before restore
	if config.CleanUpBeforeRestore.Bool {
		users, err := client.CognitoClient.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(config.CognitoUserPoolID),
		})
		if err != nil {
			return fmt.Errorf("[CLEANUP] Failed to get lists of cognito users. Error: %w", err)
		}

		for _, user := range users.Users {
			_, err := client.CognitoClient.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
				UserPoolId: aws.String(config.CognitoUserPoolID),
				Username:   user.Username,
			})
			if err != nil {
				log.Errorf("[CLEANUP] Failed to user %s. Error: %s", *user.Username, err)
			} else {
				log.Debug("User %s has been successfully deleted from %s userpool", *user.Username, config.CognitoUserPoolID)
			}
		}

		// TODO implement groups cleanup

		time.Sleep(3 * time.Second)
		log.Infof("User pool %s has been successfully cleaned up", config.CognitoUserPoolID)
	}

	if config.RestoreUsers.Bool {

		// Get the data object in S3
		data, err := getDataFromS3(ctx, client, config.S3BucketName, fmt.Sprintf("%s/users.json", config.BackupDirPath))
		if err != nil {
			return fmt.Errorf("Failed to get users backup data from S3. Error: %w", err)
		} else {
			log.Debugf("%s/users.json data has been received successfully from S3", config.BackupDirPath)
		}

		encData, err := decryptViaKMS(ctx, client, config.KMSKeyName, data)
		if err != nil {
			return fmt.Errorf("Failed to decrypt users backup data. Error: %w", err)
		} else {
			log.Debugf("data is not able to decrypt by KMS key", config.KMSKeyName)
		}

		// Get the decrypted data
		decryptedData := encData.Plaintext

		var users cognitoidentityprovider.ListUsersOutput
		err = json.Unmarshal(decryptedData, &users)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal users backup data. Error: %w", err)
		} else {
			log.Debug("users data has been unmarshalled successfully")
		}

		for _, user := range users.Users {
			var userAttributes []types.AttributeType
			var userName *string

			for _, attribute := range user.Attributes {
				if *attribute.Name == "email" {
					userName = attribute.Value
				}

				if *attribute.Name != "sub" {
					userAttributes = append(userAttributes, attribute)
				}
			}

			_, err := client.CognitoClient.AdminCreateUser(
				ctx, &cognitoidentityprovider.AdminCreateUserInput{
					UserPoolId:     aws.String(config.CognitoUserPoolID),
					Username:       userName,
					UserAttributes: userAttributes,
				},
			)
			if err != nil {
				return fmt.Errorf("Failed to restore user %s. Error: %w", *user.Username, err)
			}
		}
	}

	// TODO restore Cognito Groups
	return nil
}
