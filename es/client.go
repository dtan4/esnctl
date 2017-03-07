package es

// Client represents innterface of Elasticsearch API client
type Client interface {
	DisableReallocation() error
	EnableReallocation() error
	ListNodes() ([]string, error)
}
