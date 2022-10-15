package behavior

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/senrok/yadal"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	baseTests = []testFunc{testMetadata, testObjectId, testObjectPath, testMetadata, testObjectName}
)

func testMetadata(t *testing.T, op *yadal.Operator) {
	meta := op.Metadata()
	assert.NotNil(t, meta)
}

func testObjectId(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()
	o := op.Object(path)
	assert.Equal(t, o.ID(), fmt.Sprintf("%s%s", op.Metadata().Root(), path))
}

func testObjectPath(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()
	o := op.Object(path)
	assert.Equal(t, o.Path(), path)
}

func testObjectName(t *testing.T, op *yadal.Operator) {
	// Normal
	path := uuid.New().String()
	o := op.Object(path)
	assert.Equal(t, o.Name(), path)

	//  FIle in subdirectory
	name := uuid.New().String()
	path = fmt.Sprintf("%s/%s", uuid.New().String(), name)
	o = op.Object(path)
	assert.Equal(t, o.Name(), name)

	// Dir in subdirectory
	name = uuid.New().String()
	path = fmt.Sprintf("%s/%s/", uuid.New().String(), name)

	o = op.Object(path)
	assert.Equal(t, fmt.Sprintf("%s/", name), o.Name())
}
