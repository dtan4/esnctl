package v1

import (
	"testing"

	"gopkg.in/h2non/gock.v1"
)

const testClusterEndpoint = "http://example.com:9200"

func TestDisableReallocation(t *testing.T) {
	defer gock.Off()

	client := &Client{
		client:          nil,
		clusterEndpoint: testClusterEndpoint,
	}

	gock.New(testClusterEndpoint).Put("/_cluster/settings").BodyString(`{"transient":{"cluster.routing.allocation.enable":"none"}}`).Reply(200)

	if err := client.DisableReallocation(); err != nil {
		t.Errorf("error should not be raised: %s", err)
	}
}
