package interfaces

type Provider int

var (
	provider2Str = []string{"Unknown", "S3"}
)

const (
	S3 Provider = iota + 1
)

func (p Provider) String() string {
	return provider2Str[p]
}
