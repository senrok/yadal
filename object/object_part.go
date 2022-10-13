package object

type ObjectPart struct {
	PartNumber uint
	ETag       string
}

func (o ObjectPart) GetPartNumber() uint {
	return o.PartNumber
}

func (o ObjectPart) GetETag() string {
	return o.ETag
}
