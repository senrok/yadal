package s3

import (
	"encoding/xml"
	"github.com/senrok/yadal/interfaces"
)

type Part struct {
	PartNumber uint   `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}

type CompleteMultipartUpload struct {
	XMLName xml.Name `xml:"CompleteMultipartUpload"`
	Parts   []Part   `xml:"Part"`
}

func NewCompleteMultipartUploadFromObjectParts(input []interfaces.ObjectPart) CompleteMultipartUpload {
	parts := make([]Part, 0, len(input))
	for _, part := range input {
		parts = append(parts, Part{
			PartNumber: part.GetPartNumber(),
			ETag:       part.GetETag(),
		})
	}

	return CompleteMultipartUpload{Parts: parts}
}
