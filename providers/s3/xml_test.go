package s3

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeCompleteMultipartUpload(t *testing.T) {
	data := CompleteMultipartUpload{
		Parts: []Part{
			{
				PartNumber: 1,
				ETag:       "a54357aff0632cce46d942af68356b38",
			},
			{
				PartNumber: 2,
				ETag:       "0c78aef83f66abc1fa1e8477f296d394",
			},
		},
	}
	b, err := xml.Marshal(data)
	assert.Nil(t, err)
	expected := `<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>a54357aff0632cce46d942af68356b38</ETag></Part><Part><PartNumber>2</PartNumber><ETag>0c78aef83f66abc1fa1e8477f296d394</ETag></Part></CompleteMultipartUpload>`
	assert.Equal(t, expected, string(b))
}
