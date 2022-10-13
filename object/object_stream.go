package object

import (
	"context"
	"github.com/senrok/yadal/interfaces"
)

type Stream struct {
	done          bool
	stream        interfaces.ObjectPageStream
	cachedEntries []interfaces.Entry
}

func (o *Stream) HasNext() bool {
	return len(o.cachedEntries) > 0 || !o.done
}

func (o *Stream) Next(ctx context.Context) (entry interfaces.Entry, err error) {
	for {
		if len(o.cachedEntries) > 0 {
			entry, o.cachedEntries = o.cachedEntries[0], o.cachedEntries[1:]
			if len(o.cachedEntries) == 0 {
				o.done = true
			}
			return
		}
		// fetches more
		if !o.done {
			if o.cachedEntries, err = o.stream.NextPage(ctx); err != nil {
				return nil, err
			}
			if len(o.cachedEntries) == 0 {
				o.done = true
				// TODO: handles empty dir?
				return
			}
		} else {
			return nil, nil
		}
	}
}

func NewObjectStream(stream interfaces.ObjectPageStream) interfaces.ObjectStream {
	return &Stream{
		stream: stream,
	}
}
