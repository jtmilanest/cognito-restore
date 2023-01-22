package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	CognitoClient *cognitoidentityprovider.Client
	S3Client      *s3.Client
	KMSClient     *kms.Client
}

// Factory function to instantiate Cognito and S3 process
func New(ctx context.Context, cognitoRegion, s3BucketRegion, kmsKeyRegion string) (*Client, error) {
	// Cognito
	cognitoCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cognitoRegion))
	if err != nil {
		return nil, err
	}

	// S3
	s3Cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(s3BucketRegion))
	if err != nil {
		return nil, err
	}

	// KMS
	kmsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(kmsKeyRegion))
	if err != nil {
		return nil, err
	}

	return &Client{
		CognitoClient: cognitoidentityprovider.NewFromConfig(cognitoCfg),
		S3Client:      s3.NewFromConfig(s3Cfg),
		KMSClient:     kms.NewFromConfig(kmsCfg),
	}, nil

}
