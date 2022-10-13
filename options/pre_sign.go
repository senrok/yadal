package options

import "time"

type PreSignOptions struct {
	Op PreSignOperation
	*ReadOptions
	*WriteOptions
	*WriteMultipart
	Expire time.Duration
}

type PreSignOperation int

const (
	ReadOp PreSignOperation = iota + 1
	WriteOp
	WriteMultipartOp
)

var (
	op2str = []string{"Unknown", "Read", "Write", "WriteMultipart"}
)

func (p PreSignOperation) String() string {
	return op2str[p]
}
