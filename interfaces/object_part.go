package interfaces

type objectPart struct {
	PartNumber int
	Etag       string
}

type ObjectPart interface {
	GetPartNumber() uint
	GetETag() string
}
