package layers

import (
	"github.com/senrok/yadal/interfaces"
)

type Middleware func(Layer) Layer

type Layer func(interfaces.Accessor) interfaces.Accessor
