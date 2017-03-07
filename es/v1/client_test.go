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

func TestEnableReallocation(t *testing.T) {
	defer gock.Off()

	client := &Client{
		client:          nil,
		clusterEndpoint: testClusterEndpoint,
	}

	gock.New(testClusterEndpoint).Put("/_cluster/settings").BodyString(`{"transient":{"cluster.routing.allocation.enable":"all"}}`).Reply(200)

	if err := client.EnableReallocation(); err != nil {
		t.Errorf("error should not be raised: %s", err)
	}
}

func TestExcludeNodeFromAllocation(t *testing.T) {
	defer gock.Off()

	client := &Client{
		client:          nil,
		clusterEndpoint: testClusterEndpoint,
	}

	gock.New(testClusterEndpoint).Put("/_cluster/settings").BodyString(`{"transient":{"cluster.routing.allocation.exclude._name":"ip-10-0-1-23.ap-northeast-1.compute.internal"}}`).Reply(200)

	nodeName := "ip-10-0-1-23.ap-northeast-1.compute.internal"

	if err := client.ExcludeNodeFromAllocation(nodeName); err != nil {
		t.Errorf("error should not be raised: %s", err)
	}
}
