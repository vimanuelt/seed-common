package registry

import (
	"strings"

	"github.com/ngageoint/seed-common/objects"
	"github.com/ngageoint/seed-common/registry/containeryard"
	"github.com/ngageoint/seed-common/registry/dockerhub"
	"github.com/ngageoint/seed-common/registry/v2"
)

type RepositoryRegistry interface {
	Name() string
	Ping() error
	Repositories() ([]string, error)
	Tags(repository string) ([]string, error)
	Images() ([]string, error)
	ImagesWithManifests() ([]objects.Image, error)
	GetImageManifest(repoName, tag string) (string, error)
}

type RepoRegistryFactory func(url, org, username, password string) (RepositoryRegistry, error)

func NewV2Registry(url, org, username, password string) (RepositoryRegistry, error) {
	v2registry, err := v2.New(url, org, username, password)
	if err != nil {
		if strings.Contains(url, "https://") {
			httpFallback := strings.Replace(url, "https://", "http://", 1)
			v2registry, err = v2.New(httpFallback, org, username, password)
		}
	}

	return v2registry, err
}

func NewDockerHubRegistry(url, org, username, password string) (RepositoryRegistry, error) {
	hub, err := dockerhub.New(url, org)
	if err != nil {
		if strings.Contains(url, "https://") {
			httpFallback := strings.Replace(url, "https://", "http://", 1)
			hub, err = dockerhub.New(httpFallback, org)
		}
	}

	return hub, err
}

func NewContainerYardRegistry(url, org, username, password string) (RepositoryRegistry, error) {
	yard, err := containeryard.New(url, org, username, password)
	if err != nil {
		if strings.Contains(url, "https://") {
			httpFallback := strings.Replace(url, "https://", "http://", 1)
			yard, err = containeryard.New(httpFallback, org, username, password)
		}
	}

	return yard, err
}

func CreateRegistry(url, org, username, password string) (RepositoryRegistry, error) {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	v2, err1 := NewV2Registry(url, org, username, password)
	if err1 == nil && v2 != nil && v2.Ping() == nil {
		return v2, nil
	}

	hub, err2 := NewDockerHubRegistry(url, org, username, password)
	if err2 == nil && hub != nil && hub.Ping() == nil {
		return hub, nil
	}

	yard, err3 := NewContainerYardRegistry(url, org, username, password)
	if err3 == nil && yard != nil && yard.Ping() == nil {
		return yard, nil
	}

	return nil, err1
}
