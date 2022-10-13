package options

type ObjectPart interface {
	GetPartNumber() uint
	GetETag() string
}

type CompleteMultipart struct {
	UploadId    string
	ObjectParts []ObjectPart
}
