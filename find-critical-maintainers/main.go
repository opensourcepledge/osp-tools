// © 2025 Functional Software, Inc. d/b/a Sentry
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"log"
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

func tryDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getUrlForCriticalPackages(page int) string {
	return fmt.Sprintf("https://packages.ecosyste.ms/api/v1/packages/critical?page=%d&per_page=1000", page)
}

func getJson(url string, target interface{}) error {
	log.Printf("GET %s\n", url)
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func getCriticalPackages() []EcosystemsPackage {
	packages := []EcosystemsPackage{}
	page := 1
	for {
		packagesPage := []EcosystemsPackage{}
		err := getJson(getUrlForCriticalPackages(page), &packagesPage)
		if err != nil {
			log.Fatalln(err)
		}
		page += 1
		log.Printf("+%d\n", len(packagesPage))
		if len(packagesPage) == 0 {
			break
		}
		packages = append(packages, packagesPage...)
	}
	return packages
}

type EcosystemsMaintainerRef struct {
	Login string
	Name  string
	Email string
}

type EcosystemsMaintainerStats struct {
	NCriticalPackages int
}

func findSimilarMaintainerRef(maintainerMap map[EcosystemsMaintainerRef]EcosystemsMaintainerStats, ref EcosystemsMaintainerRef) (EcosystemsMaintainerRef, bool) {
	for k := range maintainerMap {
		isSimilar := (len(k.Login) > 0 && k.Login == ref.Login) ||
			(len(k.Name) > 0 && k.Name == ref.Name) ||
			(len(k.Email) > 0 && k.Email == ref.Email)
		if isSimilar {
			return k, true
		}
	}
	return EcosystemsMaintainerRef{}, false
}

func mergeMaintainerRefs(a EcosystemsMaintainerRef, b EcosystemsMaintainerRef) EcosystemsMaintainerRef {
	if len(a.Login) == 0 {
		a.Login = b.Login
	}
	if len(a.Name) == 0 {
		a.Name = b.Name
	}
	if len(a.Email) == 0 {
		a.Email = b.Email
	}
	return a
}

func getMaintainerMapFromPackages(packages []EcosystemsPackage) map[EcosystemsMaintainerRef]EcosystemsMaintainerStats {
	maintainerMap := make(map[EcosystemsMaintainerRef]EcosystemsMaintainerStats)
	for _, pkg := range packages {
		for _, maintainer := range pkg.Maintainers {
			ref := EcosystemsMaintainerRef{
				Login: tryDeref(maintainer.Login),
				Name:  tryDeref(maintainer.Name),
				Email: tryDeref(maintainer.Email),
			}
			if similarRef, gotSimilar := findSimilarMaintainerRef(maintainerMap, ref); gotSimilar {
				val := EcosystemsMaintainerStats{
					NCriticalPackages: maintainerMap[similarRef].NCriticalPackages + 1,
				}
				delete(maintainerMap, similarRef)
				mergedRef := mergeMaintainerRefs(ref, similarRef)
				maintainerMap[mergedRef] = val
			} else if currStats, ok := maintainerMap[ref]; ok {
				maintainerMap[ref] = EcosystemsMaintainerStats{
					NCriticalPackages: currStats.NCriticalPackages + 1,
				}
			} else {
				maintainerMap[ref] = EcosystemsMaintainerStats{
					NCriticalPackages: 1,
				}
			}
		}
	}
	return maintainerMap
}

func main() {
	log.SetFlags(0)
	packages := getCriticalPackages()
	maintainerMap := getMaintainerMapFromPackages(packages)
	log.Println(len(packages))
	for k := range maintainerMap {
		log.Printf("%03d\t%s; %s; %s\n",
			maintainerMap[k].NCriticalPackages,
			k.Login, k.Name, k.Email)
	}
}
