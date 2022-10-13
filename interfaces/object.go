package interfaces

//type object struct {
//	accessor *Accessor
//	path     string
//}

type Object interface {
	GetAccessor() Accessor
	GetPath() string
}
