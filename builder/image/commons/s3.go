package commons

import (
	"github.com/deis/deis/boot/logger"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
)

// ConnectS3Store Connnect to a S3 compatible store
func ConnectS3Store(accessKey string, secretKey string, host string) *s3.S3 {
	logger.Log.Debug("connecting to the ceph data store")
	auth := aws.Auth{AccessKey: accessKey, SecretKey: secretKey}
	return s3.New(auth, aws.Region{Name: "deis-region-1", S3Endpoint: host})
}
