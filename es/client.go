package es

// Client represents innterface of Elasticsearch API client
type Client interface {
	ListNodes() ([]string, error)
}
