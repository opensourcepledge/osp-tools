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
	// "os"
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

type EcosystemsCommitter struct {
	Login *string `json:"login"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Count int     `json:"count"`
}

type EcosystemsHost struct {
	Name              string     `json:"name"`
	Url               string     `json:"name"`
	Kind              string     `json:"kind"`
	LastSyncedAt      *time.Time `json:"last_synced_at"`
	RepositoriesCount int        `json:"repositories_count"`
	CommitsCount      int        `json:"commits_count"`
	ContributorsCount int        `json:"contributors_count"`
	OwnersCount       int        `json:"owners_count"`
	IconUrl           *string    `json:"icon_url"`
	HostUrl           string     `json:"host_url"`
	RepositoriesUrl   string     `json:"repositories_url"`
}

type EcosystemsCommits struct {
	FullName                   string                `json:"full_name"`
	DefaultBranch              string                `json:"default_branch"`
	Committers                 []EcosystemsCommitter `json:"committers"`
	TotalCommits               int                   `json:"total_commits"`
	TotalCommitters            int                   `json:"total_committers"`
	TotalBotCommits            int                   `json:"total_bot_commits"`
	TotalBotCommitters         int                   `json:"total_bot_committers"`
	MeanCommits                float32               `json:"mean_commits"`
	Dds                        float32               `json:"dds"`
	PastYearCommitters         []EcosystemsCommitter `json:"past_year_committers"`
	PastYearTotalCommits       int                   `json:"past_year_total_commits"`
	PastYearTotalCommitters    int                   `json:"past_year_total_committers"`
	PastYearTotalBotCommits    int                   `json:"past_year_total_bot_commits"`
	PastYearTotalBotCommitters int                   `json:"past_year_total_bot_committers"`
	PastYearMeanCommits        float32               `json:"past_year_mean_commits"`
	PastYearDds                float32               `json:"past_year_dds"`
	LastSyncedAt               *time.Time            `json:"last_synced_at"`
	LastSyncedCommit           *string               `json:"last_synced_commit"`
	CreatedAt                  time.Time             `json:"created_at"`
	UpdatedAt                  time.Time             `json:"updated_at"`
	CommitsUrl                 string                `json:"commits_url"`
	Host                       EcosystemsHost        `json:"host"`
}

type EcosystemsPackage struct {
	// Fields from https://packages.ecosyste.ms/docs/index.html
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

type EcosystemsSummary struct {
	// A bunch of fields we don't care about right now, and then...
	Commits EcosystemsCommits `json:"commits"`
}

type EcosystemsMaintainerRef struct {
	Login string
	Name  string
	Email string
}

type EcosystemsMaintainerStats struct {
	NCriticalPackages int
}

// How many packages `getCriticalPackages()` should fetch at a time.
const PACKAGES_PER_PAGE int = 10

// Whether `getCriticalPacakges()` should only fetch one page.
const GET_ONLY_ONE_PAGE bool = true

// The fraction of all repository commits that a committer must have committed
// to be considered a significant committer, which justifies inclusion in the
// list of maintainers.
const SIGNIFICANT_COMMITTER_FRACTION float32 = 0.2

var httpClient = &http.Client{}

func tryDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getUrlForCriticalPackages(page int) string {
	return fmt.Sprintf("https://packages.ecosyste.ms/api/v1/packages/critical?page=%d&per_page=%d",
		page, PACKAGES_PER_PAGE)
}

func getUrlForPackageSummary(url string) string {
	return fmt.Sprintf("https://summary.ecosyste.ms/api/v1/projects/lookup?url=%s", url)
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
		if GET_ONLY_ONE_PAGE {
			break
		}
	}
	return packages
}

func getPackageSummary(url string) EcosystemsSummary {
	summary := EcosystemsSummary{}
	err := getJson(getUrlForPackageSummary(url), &summary)
	if err != nil {
		log.Fatalln(err)
	}
	return summary
}

// Finds whether an EcosystemsMaintainerRef that is similar to `ref` already
// exists in `maintainerMap`, where two refs are similar if they have the same
// non-empty login, name or email. This is a heuristic, which may lead to false
// aggregation of actually unrelated refs. It should be a decent compromise for
// now.
func findSimilarMaintainerRef(
	maintainerMap map[EcosystemsMaintainerRef]EcosystemsMaintainerStats,
	ref EcosystemsMaintainerRef,
) (EcosystemsMaintainerRef, bool) {
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

func getMaintainersFromCommits(commits EcosystemsCommits) []EcosystemsMaintainer {
	maintainers := []EcosystemsMaintainer{}
	for _, committer := range commits.Committers {
		commitFraction := float32(committer.Count) / float32(commits.TotalCommits)
		// NOTE: This works because committers are descendingly sorted by
		// commit count in the ecosyste.ms API response. If this stops being
		// the case, this code will break.
		if commitFraction < SIGNIFICANT_COMMITTER_FRACTION {
			break
		}
		newMaintainer := EcosystemsMaintainer{
			Login: committer.Login,
			Name:  committer.Name,
			Email: committer.Email,
		}
		maintainers = append(maintainers, newMaintainer)
	}
	return maintainers
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

func mergeMaintainers(a []EcosystemsMaintainer, b []EcosystemsMaintainer) []EcosystemsMaintainer {
	for _, mb := range b {
		exists := false
		for _, ma := range a {
			isSimilar := (ma.Login != nil && mb.Login != nil && *ma.Login == *mb.Login) ||
				(ma.Name != nil && mb.Name != nil && *ma.Name == *mb.Name) ||
				(ma.Email != nil && mb.Email != nil && *ma.Email == *mb.Email)
			if isSimilar {
				exists = true
			}
		}
		if !exists {
			a = append(a, mb)
		}
	}
	return a
}

func main() {
	log.SetFlags(0)
	packages := getCriticalPackages()
	for idx := range packages {
		pkg := &packages[idx]
		if pkg.RepositoryUrl == nil || len(*pkg.RepositoryUrl) == 0 {
			continue
		}
		summary := getPackageSummary(*pkg.RepositoryUrl)
		maintainersFromCommits := getMaintainersFromCommits(summary.Commits)
		pkg.Maintainers = mergeMaintainers(pkg.Maintainers, maintainersFromCommits)
		// log.Println("Number of commits:")
		// spew.Fdump(os.Stderr, summary.Commits.TotalCommits)
		// log.Println("Committers:")
		// for idx, committer := range summary.Commits.Committers {
		// 	if idx > 4 {
		// 		break
		// 	}
		// 	spew.Fdump(os.Stderr, committer)
		// }
		// log.Println("Maintainers from commits:")
		// spew.Fdump(os.Stderr, maintainersFromCommits)
		// log.Println("All maintainers:")
		// spew.Fdump(os.Stderr, pkg.Maintainers)
	}
	maintainerMap := getMaintainerMapFromPackages(packages)
	// spew.Fdump(os.Stderr, maintainerMap)
	log.Println(len(packages))
	for k := range maintainerMap {
		log.Printf("%03d\t%s; %s; %s\n",
			maintainerMap[k].NCriticalPackages,
			k.Login, k.Name, k.Email)
	}
}
