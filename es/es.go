package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dtan4/esnctl/es/v1"
	"github.com/pkg/errors"
)

func New(clusterURL string) (*Client, error) {
	version, err := DetectVersion(clusterURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect Elasticsearch version")
	}

	if version == "1.0.0" {
		client, err := v1.NewClient(clusterURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create Elasticsearch API")
		}

		return client, nil
	}

	return nil, errors.Errorf("version %q does not supported", version)
}

// DetectVersion returns Elasticsearch version of the given endpoint
func DetectVersion(clusterURL string) (string, error) {
	u, err := url.Parse(clusterURL)
	if err != nil {
		return "", errors.Wrap(err, "cluster URL is invalid")
	}

	var endpoint string

	if u.User == nil {
		endpoint = fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
	} else {
		endpoint = fmt.Sprintf("%s://%s@%s/", u.Scheme, u.User.String(), u.Host)
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return "", errors.Wrap(err, "failed to access to root API")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}

	var bodyStruct map[string]interface{}

	if err := json.Unmarshal(body, &bodyStruct); err != nil {
		return "", errors.Wrap(err, "invalid response body")
	}

	if _, ok := bodyStruct["OK"]; ok {
		return "1.0.0", nil
	}

	v, ok := bodyStruct["version"]
	if !ok {
		return "", errors.New("version field not found")
	}

	version, ok := v.(map[string]interface{})
	if !ok {
		return "", errors.New("invalid version field")
	}

	if _, ok := version["number"]; !ok {
		return "", errors.New("version number field not found")
	}

	vn, ok := version["number"].(string)
	if !ok {
		return "", errors.New("invalid version number field")
	}

	return vn, nil
}
