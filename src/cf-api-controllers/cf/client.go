package cf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"code.cloudfoundry.org/capi-k8s-release/src/cf-api-controllers/cf/model"
)

func NewClient(host string, restClient Rest, uaaClient TokenFetcher) *Client {
	// TODO: We may want to consider using cloudfoundry/tlsconfig for using
	// standard TLS configs in Golang.
	return &Client{
		host:       host,
		restClient: restClient,
		uaaClient:  uaaClient,
	}
}

// TODO: remove mockery usages after refactoring everything to use Ginkgo for consistency
//go:generate mockery -case snake -name Rest
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Rest
type Rest interface {
	Patch(url string, authToken string, body io.Reader) (*http.Response, error)
}

// TODO: remove mockery usages after refactoring everything to use Ginkgo for consistency
//go:generate mockery -case snake -name TokenFetcher
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . TokenFetcher
type TokenFetcher interface {
	Fetch() (string, error)
}

// TODO: replace this with the client the cf-cli uses?
type Client struct {
	host       string
	restClient Rest
	uaaClient  TokenFetcher
}

func (c *Client) UpdateBuild(buildGUID string, build model.Build) error {
	token, err := c.uaaClient.Fetch()
	if err != nil {
		return err
	}

	raw, err := json.Marshal(build)
	if err != nil {
		return err
	}

	resp, err := c.restClient.Patch(
		fmt.Sprintf("%s/v3/builds/%s", c.host, buildGUID),
		token,
		bytes.NewReader(raw),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to patch build, received status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateDroplet(dropletGUID string, droplet model.Droplet) error {
	token, err := c.uaaClient.Fetch()
	if err != nil {
		return err
	}

	raw, err := json.Marshal(droplet)
	if err != nil {
		return err
	}

	resp, err := c.restClient.Patch(
		fmt.Sprintf("%s/v3/droplets/%s", c.host, dropletGUID),
		token,
		bytes.NewReader(raw),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to patch droplet, received status %d", resp.StatusCode)
	}

	return nil
}
