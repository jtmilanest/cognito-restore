package lambda

import (
	"context"

	"github.com/jtmilanest/cognito-restore/internal/cloud"
)

func getDataFromS3(ctx context.Context, client *cloud.Client, bucketName, keyName string) ([]byte, error) {

}
