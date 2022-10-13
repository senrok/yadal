package interfaces

import "context"

type ObjectStream interface {
	HasNext() bool
	Next(ctx context.Context) (Entry, error)
}

type ObjectPageStream interface {
	NextPage(ctx context.Context) ([]Entry, error)
}
