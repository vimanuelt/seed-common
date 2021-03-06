package dockerhub

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ngageoint/seed-common/constants"
	"github.com/ngageoint/seed-common/util"
)

//DockerHubRegistry type representing a Docker Hub registry
type DockerHubRegistry struct {
	URL    string
	Client *http.Client
	Org    string
	Print  util.PrintCallback
}

//New creates a new docker hub registry from the given URL
func New(registryUrl, org string) (*DockerHubRegistry, error) {
	if util.PrintUtil == nil {
		util.InitPrinter(util.PrintErr)
	}
	url := strings.TrimSuffix(registryUrl, "/")
	registry := &DockerHubRegistry{
		URL:    url,
		Client: &http.Client{},
		Org:    org,
		Print:  util.PrintUtil,
	}

	return registry, nil
}

func (r *DockerHubRegistry) url(pathTemplate string, args ...interface{}) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

func (r *DockerHubRegistry) Name() string {
	return "DockerHubRegistry"
}

func (r *DockerHubRegistry) Ping() error {
	url := r.url("/v2/repositories/%s/", constants.DefaultOrg)
	resp, err := r.Client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return errors.New(resp.Status)
		}
	}
	return err
}
