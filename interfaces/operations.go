package interfaces

type Operation int

const (
	MetadataOP Operation = iota + 1
	CreateOp
	ReadOp
	WriteOp
	StatOp
	DeleteOp
	ListOp
	PreSignOp
	CreateMultipartOp
	WriteMultipartOp
	CompleteMultipartOp
	AbortMultipartOp
)

var (
	op2Str = []string{
		"Unknown",
		"Metadata",
		"Create",
		"Read",
		"Write",
		"Stat",
		"Delete",
		"List",
		"PreSign",
		"CreateMultipart",
		"WriteMultipart",
		"CompleteMultipart",
		"AbortMultipart",
	}
)

func (o Operation) String() string {
	return op2Str[o]
}
