package yadal

import (
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/object"
)

type Operator struct {
	accessor interfaces.Accessor
}

// Object returns an object.Object handler
func (o *Operator) Object(path string) object.Object {
	return object.NewObject(o.accessor, path)
}

// Layer appends a layers.Layer
func (o *Operator) Layer(layer interfaces.Layer) *Operator {
	o.accessor = layer(o.accessor)
	return o
}

// NewOperatorFromAccessor returns the Operator from the interfaces.Accessor
func NewOperatorFromAccessor(acc interfaces.Accessor) Operator {
	return Operator{
		accessor: acc,
	}
}
