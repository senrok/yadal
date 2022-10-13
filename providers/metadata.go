package providers

import "github.com/senrok/yadal/interfaces"

type Metadata struct {
	provider   interfaces.Provider
	root       string
	name       string
	capability interfaces.Capability
}

func (m Metadata) Provider() interfaces.Provider {
	return m.provider
}

func (m Metadata) Root() string {
	return m.root
}

func (m Metadata) Name() string {
	return m.name
}

func (m Metadata) Capability() interfaces.Capability {
	return m.capability
}

func NewMetadata(provider interfaces.Provider, root string, name string, capability interfaces.Capability) *Metadata {
	return &Metadata{
		provider:   provider,
		root:       root,
		name:       name,
		capability: capability,
	}
}
