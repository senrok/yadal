package s3

type Options struct {
	Bucket   string
	Endpoint string
	Root     string
	Region   string

	AccessKey string
	SecretKey string

	SSEncryption               *string
	SSEncryptionAwsKmsKeyId    *string
	SSEncryptionCustomerAlgo   *string
	SSEncryptionCustomerKey    *string
	SSEncryptionCustomerKeyMD5 *string

	EnableVirtualHostStyle bool
}
