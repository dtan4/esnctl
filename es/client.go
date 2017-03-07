package es

// Client represents innterface of Elasticsearch API client
type Client interface {
	DisableReallocation() error
	EnableReallocation() error
	ExcludeNodeFromAllocation(nodeName string) error
	ListNodes() ([]string, error)
	ListShardsOnNode(nodeName string) ([]string, error)
	Shutdown(nodeName string) error
}
