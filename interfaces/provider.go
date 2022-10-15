package interfaces

type Provider int

var (
	provider2Str = []string{"Unknown", "S3", "FS"}
)

const (
	S3 Provider = iota + 1
	Fs
)

func (p Provider) String() string {
	return provider2Str[p]
}
