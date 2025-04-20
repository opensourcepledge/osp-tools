// © 2025 Functional Software, Inc. d/b/a Sentry
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"time"
)

type EcosystemsMaintainer struct {
	Uuid           string    `json:"uuid"`
	Login          *string   `json:"login"`
	Name           *string   `json:"name"`
	Email          *string   `json:"email"`
	Url            *string   `json:"url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PackagesCount  int       `json:"packages_count"`
	PackagesUrl    string    `json:"packages_url"`
	TotalDownloads int       `json:"total_downloads"`
	HtmlUrl        *string   `json:"html_url"`
	Role           *string   `json:"role"`
}

type EcosystemsPackage struct {
	Id                       int        `json:"id"`
	Name                     string     `json:"name"`
	Ecosystem                string     `json:"ecosystem"`
	Description              *string    `json:"description"`
	Homepage                 *string    `json:"homepage"`
	Licenses                 *string    `json:"licenses"`
	NormalizedLicenses       []string   `json:"normalized_licenses"`
	RepositoryUrl            *string    `json:"repository_url"`
	KeywordsArray            []string   `json:"keywords_array"`
	Namespace                *string    `json:"namespace"`
	VersionsCount            int        `json:"versions_count"`
	FirstReleasePublishedAt  *time.Time `json:"first_release_published_at"`
	LatestReleasePublishedAt *time.Time `json:"latest_release_published_at"`
	LatestReleaseNumber      *string    `json:"latest_release_number"`
	LastSyncedAt             *time.Time `json:"last_synced_at"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	RegistryUrl              *string    `json:"registry_url"`
	DocumentationUrl         *string    `json:"documentation_url"`
	InstallCommand           *string    `json:"install_command"`
	// Metadata Interface{} `json:"metadata"`
	// RepoMetadata Interface{} `json:"repo_metadata"`
	RepoMetadataUpdatedAt  time.Time `json:"repo_metadata_updated_at"`
	DependentPackagesCount int       `json:"dependent_packages_count"`
	Downloads              int       `json:"downloads"`
	DownloadsPeriod        *string   `json:"downloads_period"`
	DependentReposCount    int       `json:"dependent_repos_count"`
	// Rankings Interface{} `json:"rankings"`
	Purl *string `json:"purl"`
	// Advisories Interface{} `json:"advisories"`
	VersionsUrl              string                 `json:"versions_url"`
	VersionNumbersUrl        string                 `json:"version_numbers_url"`
	DependentPackagesUrl     string                 `json:"dependent_packages_url"`
	RelatedPackagesUrl       string                 `json:"related_packages_url"`
	DockerUsageUrl           string                 `json:"docker_usage_url"`
	DockerDependentsCount    int                    `json:"docker_dependents_count"`
	DockerDownloadsCount     int                    `json:"docker_downloads_count"`
	UsageUrl                 string                 `json:"usage_url"`
	DependentRepositoriesUrl string                 `json:"dependent_repositories_url"`
	Status                   *string                `json:"status"`
	FundingLinks             []string               `json:"funding_links"`
	Maintainers              []EcosystemsMaintainer `json:"maintainers"`
	Critical                 bool                   `json:"critical"`
}

var httpClient = &http.Client{}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	packages := new([]EcosystemsPackage)
	getJson("https://packages.ecosyste.ms/api/v1/packages/critical?page=1&per_page=1000", packages)
	spew.Dump(packages)
}
