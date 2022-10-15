package behavior

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/senrok/yadal"
	"github.com/senrok/yadal/errors"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/object"
	"github.com/senrok/yadal/options"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"testing"
)

var readWriteTests = []testFunc{
	testCreateFile,
	testCreateFileExisting,
	//testCreateFileWithSpecialChars,
	testCreateDir,
	testCreateDirExisting,
	testWrite,
	testWriteWithDirPath,
	//testWriteFileWithSpecialChars,
	testStat,
	testStatDir,
	//testStatWithSpecialChars,
	testStatNotCleanedPath,
	testStatNotExists,
	testStatRoot,
	testReadFull,
	testRangeRead,
	testReadNotExists,
	testReadWithDirPath,
}

// Creates a file with the path should be success
func testCreateFile(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()
	o := op.Object(path)

	err := o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Nilf(t, err, "%s", err)

	assert.Equal(t, interfaces.FILE, meta.Mode())
	assert.Equal(t, uint64(0), *meta.ContentLength())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// Creates a file on an existing file should be success
func testCreateFileExisting(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()
	o := op.Object(path)
	err := o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	err = o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Nilf(t, err, "%s", err)

	assert.Equal(t, interfaces.FILE, meta.Mode())
	assert.Equal(t, uint64(0), *meta.ContentLength())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// creates a file with special chars should be success
func testCreateFileWithSpecialChars(t *testing.T, op *yadal.Operator) {
	path := fmt.Sprintf("%s %s", uuid.New().String(), "!@#$%^&*()_+-=;'><,?.txt")

	o := op.Object(path)
	err := o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Nilf(t, err, "%s", err)

	assert.Equal(t, interfaces.FILE, meta.Mode())
	assert.Equal(t, uint64(0), *meta.ContentLength())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// creates a dir with dir path should be success
func testCreateDir(t *testing.T, op *yadal.Operator) {
	path := fmt.Sprintf("%s/", uuid.New().String())

	o := op.Object(path)
	err := o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Equal(t, interfaces.DIR, meta.Mode())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// creates dir on existing dir should be success
func testCreateDirExisting(t *testing.T, op *yadal.Operator) {
	path := fmt.Sprintf("%s/", uuid.New().String())

	o := op.Object(path)
	err := o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	err = o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Equal(t, interfaces.DIR, meta.Mode())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

func genBytes(size int) []byte {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return bytes
}

// write bytes into a file should be success
func testWrite(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()

	o := op.Object(path)
	content := genBytes(4096)
	err := o.Write(context.TODO(), content)
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Equal(t, interfaces.FILE, meta.Mode())
	assert.Equal(t, uint64(len(content)), *meta.ContentLength())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// write bytes into a dir path should be failed
func testWriteWithDirPath(t *testing.T, op *yadal.Operator) {
	path := fmt.Sprintf("%s/", uuid.New().String())

	o := op.Object(path)
	content := genBytes(4096)
	err := o.Write(context.TODO(), content)
	assert.NotNilf(t, err, "%s", err)
	assert.True(t, errors.Is(err, object.ErrTryWrite2Dir))
}

// writes bytes into a file with special chars should be success
func testWriteFileWithSpecialChars(t *testing.T, op *yadal.Operator) {
	path := fmt.Sprintf("%s %s", uuid.New().String(), "!@#$%^&*()_+-=;'><,?.txt")

	o := op.Object(path)
	bytes := genBytes(4096)
	err := o.Write(context.TODO(), bytes)
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Nilf(t, err, "%s", err)

	assert.Equal(t, interfaces.FILE, meta.Mode())
	assert.Equal(t, uint64(len(bytes)), *meta.ContentLength())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// stats an existing file should success
func testStat(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()

	o := op.Object(path)
	content := genBytes(4096)
	err := o.Write(context.TODO(), content)
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Equal(t, interfaces.FILE, meta.Mode())
	assert.Equal(t, uint64(len(content)), *meta.ContentLength())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// stats an existing dir should success
func testStatDir(t *testing.T, op *yadal.Operator) {
	path := fmt.Sprintf("%s/", uuid.New().String())

	o := op.Object(path)
	err := o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Equal(t, interfaces.DIR, meta.Mode())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// stats an existing file with special chars should be success
func testStatWithSpecialChars(t *testing.T, op *yadal.Operator) {
	path := fmt.Sprintf("%s %s", uuid.New().String(), "!@#$%^&*()_+-=;'><,?.txt")

	o := op.Object(path)
	content := genBytes(4096)
	err := o.Write(context.TODO(), content)
	assert.Nilf(t, err, "%s", err)

	meta, err := o.Metadata(context.TODO())
	assert.Equal(t, interfaces.FILE, meta.Mode())
	assert.Equal(t, uint64(len(content)), *meta.ContentLength())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// stats a non-cleaned path should be success
func testStatNotCleanedPath(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()
	bytes := genBytes(4096)
	o := op.Object(path)
	err := o.Write(context.TODO(), bytes)
	assert.Nilf(t, err, "%s", err)

	o2 := op.Object(fmt.Sprintf("//%s", path))
	meta, err := o2.Metadata(context.TODO())
	assert.Nilf(t, err, "%s", err)
	assert.Equal(t, interfaces.FILE, meta.Mode())
	assert.Equal(t, uint64(len(bytes)), *meta.ContentLength())

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// stat a not existing file return notfound
func testStatNotExists(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()

	o := op.Object(path)
	_, err := o.Metadata(context.TODO())
	assert.True(t, errors.Is(err, errors.ErrNotFound))
}

// stat root dir should be success
func testStatRoot(t *testing.T, op *yadal.Operator) {
	o := op.Object("")
	meta, err := o.Metadata(context.TODO())
	assert.Nilf(t, err, "%s", err)
	assert.Equal(t, interfaces.DIR, meta.Mode())

	o = op.Object("/")
	meta, err = o.Metadata(context.TODO())
	assert.Nilf(t, err, "%s", err)
	assert.Equal(t, interfaces.DIR, meta.Mode())
}

// read should be matched
func testReadFull(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()

	o := op.Object(path)
	content := genBytes(4096)
	err := o.Write(context.TODO(), content)
	assert.Nilf(t, err, "%s", err)

	o2 := op.Object(path)
	reader, err := o2.Read(context.TODO())
	assert.Nilf(t, err, "%s", err)

	bytes, err := io.ReadAll(reader)
	assert.Nilf(t, err, "%s", err)

	assert.Equal(t, bytes, content)

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

func randRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

func genOffsetLen(size int64) (uint64, uint64) {
	offset := randRange(0, size-1)
	len := randRange(1, int64(size-offset))
	return uint64(offset), uint64(len)
}

// range read should be success
func testRangeRead(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()

	o := op.Object(path)
	content := genBytes(4096)
	err := o.Write(context.TODO(), content)
	assert.Nilf(t, err, "%s", err)

	o2 := op.Object(path)
	off, l := genOffsetLen(4096)
	read, err := o2.RangeRead(context.TODO(), options.NewRangeBounds(options.Range(off, l+off)))
	assert.Nilf(t, err, "%s", err)

	bytes, err := io.ReadAll(read)
	assert.Nilf(t, err, "%s", err)
	assert.Equal(t, bytes, content[off:off+l])
	assert.Equal(t, len(bytes), len(content[off:off+l]))

	err = o.Delete(context.TODO())
	assert.Nilf(t, err, "%s", err)
}

// reads a not exists file should return not found
func testReadNotExists(t *testing.T, op *yadal.Operator) {
	path := uuid.New().String()
	o := op.Object(path)
	_, err := o.Read(context.TODO())
	assert.NotNilf(t, err, "%s", err)
	assert.True(t, errors.Is(err, errors.ErrNotFound))
}

func testReadWithDirPath(t *testing.T, op *yadal.Operator) {
	path := fmt.Sprintf("%s/", uuid.New().String())
	o := op.Object(path)
	err := o.Create(context.TODO())
	assert.Nilf(t, err, "%s", err)

	_, err = o.Read(context.TODO())
	assert.NotNilf(t, err, "%s", err)
	assert.True(t, errors.Is(err, object.ErrIsADir))
}
