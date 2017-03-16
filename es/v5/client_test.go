package v5

import (
	"context"
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

const testClusterEndpoint = "http://example.com:9200"

func TestDisableReallocation(t *testing.T) {
	defer gock.Off()

	client := &Client{
		client:          nil,
		clusterEndpoint: testClusterEndpoint,
		httpClient:      &http.Client{},
		ctx:             context.Background(),
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
		httpClient:      &http.Client{},
		ctx:             context.Background(),
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
		httpClient:      &http.Client{},
		ctx:             context.Background(),
	}

	gock.New(testClusterEndpoint).Put("/_cluster/settings").BodyString(`{"transient":{"cluster.routing.allocation.exclude._name":"ip-10-0-1-23.ap-northeast-1.compute.internal"}}`).Reply(200)

	nodeName := "ip-10-0-1-23.ap-northeast-1.compute.internal"

	if err := client.ExcludeNodeFromAllocation(nodeName); err != nil {
		t.Errorf("error should not be raised: %s", err)
	}
}

func TestListShardsOnNode(t *testing.T) {
	defer gock.Off()

	client := &Client{
		client:          nil,
		clusterEndpoint: testClusterEndpoint,
		httpClient:      &http.Client{},
		ctx:             context.Background(),
	}

	gock.New(testClusterEndpoint).Get("/_cat/shards").Reply(200).BodyString(`wiki1 0 p STARTED 3014 31.1mb 192.168.56.10 ip-10-0-1-23.ap-northeast-1.compute.internal
wiki1 1 p STARTED 3013 29.6mb 192.168.56.30 Frankie Raye
wiki1 2 p STARTED 3973 38.1mb 192.168.56.20 Commander Kraken`)

	nodeName := "ip-10-0-1-23.ap-northeast-1.compute.internal"

	shards, err := client.ListShardsOnNode(nodeName)
	if err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	if len(shards) != 1 {
		t.Errorf("number of shards does not match. expected: 1, got: %d", len(shards))
	}

	expected := "wiki1 0 p STARTED 3014 31.1mb 192.168.56.10 ip-10-0-1-23.ap-northeast-1.compute.internal"

	if shards[0] != expected {
		t.Errorf("shard does not match. expected: %q, got: %q", expected, shards[0])
	}
}
