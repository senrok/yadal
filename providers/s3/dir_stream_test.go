package s3

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeOutput(t *testing.T) {
	input := `<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Name>example-bucket</Name>
  <Prefix>photos/2006/</Prefix>
  <KeyCount>3</KeyCount>
  <MaxKeys>1000</MaxKeys>
  <Delimiter>/</Delimiter>
  <IsTruncated>false</IsTruncated>
  <Contents>
    <Key>photos/2006</Key>
    <LastModified>2016-04-30T23:51:29.000Z</LastModified>
    <ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag>
    <size>56</size>
    <StorageClass>STANDARD</StorageClass>
  </Contents>
  <Contents>
    <Key>photos/2007</Key>
    <LastModified>2016-04-30T23:51:29.000Z</LastModified>
    <ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag>
    <size>100</size>
    <StorageClass>STANDARD</StorageClass>
  </Contents>

  <CommonPrefixes>
    <Prefix>photos/2006/February/</Prefix>
  </CommonPrefixes>
  <CommonPrefixes>
    <Prefix>photos/2006/January/</Prefix>
  </CommonPrefixes>
</ListBucketResult>`
	output := Output{}
	err := xml.Unmarshal([]byte(input), &output)
	assert.Nil(t, err)
	assert.Nil(t, output.NextContinuationToken)
	assert.False(t, *output.IsTruncated)
	expected := "[{photos/2006 56 2016-04-30 23:51:29 +0000 UTC \"d41d8cd98f00b204e9800998ecf8427e\"} {photos/2007 100 2016-04-30 23:51:29 +0000 UTC \"d41d8cd98f00b204e9800998ecf8427e\"}]"
	assert.Equal(t, expected, fmt.Sprintf("%v", output.Contents))
	expected = `[{photos/2006/February/} {photos/2006/January/}]`
	assert.Equal(t, expected, fmt.Sprintf("%v", output.CommonPrefixes))
}
