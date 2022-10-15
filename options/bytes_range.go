package options

import "fmt"

type bytesRange struct {
	offset *uint64
	size   *uint64
}

func NewBytesRange(offset, size *uint64) RangeBounds {
	return &bytesRange{
		offset: offset,
		size:   size,
	}
}

func (br bytesRange) Offset() *uint64 {
	return br.offset
}

func (br bytesRange) Size() *uint64 {
	return br.size
}

func (br bytesRange) String() string {
	if br.size != nil && br.offset != nil {
		return fmt.Sprintf("bytes=%d-%d", *br.offset, *br.offset+*br.size-1)
	} else if br.offset != nil {
		return fmt.Sprintf("bytes=%d-", *br.offset)
	} else if br.size != nil {
		return fmt.Sprintf("bytes=-%d", *br.size-1)
	}
	panic("unreachable")
}

type RangeBounds interface {
	Offset() *uint64
	Size() *uint64
	String() string
}

type rangeBounds struct {
	start *uint64
	end   *uint64
}

func (r rangeBounds) Offset() *uint64 {
	return r.start
}

func (r rangeBounds) Size() *uint64 {
	if r.end == nil {
		return nil
	} else if r.start == nil {
		return r.end
	}
	l := *r.end - *r.start
	return &l
}

func (r rangeBounds) String() string {
	if r.Size() != nil && r.Offset() != nil {
		return fmt.Sprintf("bytes=%d-%d", *r.Offset(), *r.Offset()+*r.Size()-1)
	} else if r.Offset() != nil {
		return fmt.Sprintf("bytes=%d-", *r.Offset())
	} else if r.Size() != nil {
		return fmt.Sprintf("bytes=-%d", *r.Size()-1)
	}
	panic("unreachable")
}

type BuildRangeBounds func(bounds *rangeBounds)

func Start(start uint64) func(bounds *rangeBounds) {
	return func(bounds *rangeBounds) {
		bounds.start = &start
	}
}

func End(end uint64) func(bounds *rangeBounds) {
	return func(bounds *rangeBounds) {
		bounds.end = &end
	}
}

func Range(start, end uint64) func(bounds *rangeBounds) {
	return func(bounds *rangeBounds) {
		bounds.start = &start
		bounds.end = &end
	}
}

func NewRangeBounds(opts ...BuildRangeBounds) RangeBounds {
	r := &rangeBounds{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
