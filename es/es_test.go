package es

import (
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func TestDetectVersion(t *testing.T) {
	defer gock.Off()

	testcases := []struct {
		clusterURL string
		body       string
		expected   string
	}{
		{
			clusterURL: "http://example-v1.com:9200",
			body:       `{"OK":{}}`,
			expected:   "1.0.0",
		},
		{
			clusterURL: "http://example-v1.com:9200/",
			body:       `{"OK":{}}`,
			expected:   "1.0.0",
		},
		{
			clusterURL: "http://user:pass@example-v1.com:9200/",
			body:       `{"OK":{}}`,
			expected:   "1.0.0",
		},
		{
			clusterURL: "http://example-v2.com:9200",
			body: `{
  "name" : "ip-10-0-1-23.ap-northeast-1.compute.internal",
  "cluster_name" : "elasticsearch",
  "version" : {
    "number" : "2.3.0",
    "build_hash" : "8371be8d5fe5df7fb9c0516c474d77b9feddd888",
    "build_timestamp" : "2016-03-29T07:54:48Z",
    "build_snapshot" : false,
    "lucene_version" : "5.5.0"
  },
  "tagline" : "You Know, for Search"
}
`,
			expected: "2.3.0",
		},
	}

	for _, tc := range testcases {
		gock.New(tc.clusterURL).Get("/").Reply(200).BodyString(tc.body)

		got, err := DetectVersion(tc.clusterURL, &http.Client{})
		if err != nil {
			t.Errorf("error should not be raised: %s", err)
		}

		if got != tc.expected {
			t.Errorf("version does not match. expected: %q, got: %q", tc.expected, got)
		}
	}
}
