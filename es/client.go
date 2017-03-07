package es

// Client represents innterface of Elasticsearch API client
type Client interface {
	EnableReallocation() error
	ListNodes() ([]string, error)
}
