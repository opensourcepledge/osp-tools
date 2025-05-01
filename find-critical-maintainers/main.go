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

type EcosystemsIssueAuthor struct {
	Login *string `json:"login"`
	Url   *string `json:"url"`
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

type EcosystemsCountTable struct {
	Table map[string]int `json:"table"`
}

type EcosystemsIssueAuthorTable struct {
	Table EcosystemsIssueAuthor `json:"table"`
}

type EcosystemsIssues struct {
	Table struct {
		FullName                                   string                       `json:"full_name"`
		HtmlUrl                                    string                       `json:"html_url"`
		LastSyncedAt                               *time.Time                   `json:"last_synced_at"`
		Status                                     *string                      `json:"status"`
		IssuesCount                                *int                         `json:"issues_count"`
		PullRequestsCount                          *int                         `json:"pull_requests_count"`
		AvgTimeToCloseIssue                        *float32                     `json:"avg_time_to_close_issue"`
		AvgTimeToClosePullRequest                  *float32                     `json:"avg_time_to_close_pull_request"`
		IssuesClosedCount                          *int                         `json:"issues_closed_count"`
		PullRequestsClosedCount                    *int                         `json:"pull_requests_closed_count"`
		PullRequestAuthorsCount                    *int                         `json:"pull_request_authors_count"`
		IssueAuthorsCount                          *int                         `json:"issue_authors_count"`
		AvgCommentsPerIssue                        *float32                     `json:"avg_comments_per_issue"`
		AvgCommentsPerPullRequest                  *float32                     `json:"avg_comments_per_pull_request"`
		MergedPullRequestsCount                    *int                         `json:"merged_pull_requests_count"`
		BotIssuesCount                             *int                         `json:"bot_issues_count"`
		BotPullRequestsCount                       *int                         `json:"bot_pull_requests_count"`
		PastYearIssuesCount                        *int                         `json:"past_year_issues_count"`
		PastYearPullRequestsCount                  *int                         `json:"past_year_pull_requests_count"`
		PastYearAvgTimeToCloseIssue                *float32                     `json:"past_year_avg_time_to_close_issue"`
		PastYearAvgTimeToClosePullRequest          *float32                     `json:"past_year_avg_time_to_close_pull_request"`
		PastYearIssuesClosedCount                  *int                         `json:"past_year_issues_closed_count"`
		PastYearPullRequestsClosedCount            *int                         `json:"past_year_pull_requests_closed_count"`
		PastYearPullRequestAuthorsCount            *int                         `json:"past_year_pull_request_authors_count"`
		PastYearIssueAuthorsCount                  *int                         `json:"past_year_issue_authors_count"`
		PastYearAvgCommentsPerIssue                *float32                     `json:"past_year_avg_comments_per_issue"`
		PastYearAvgCommentsPerPullRequest          *float32                     `json:"past_year_avg_comments_per_pull_request"`
		PastYearBotIssuesCount                     *int                         `json:"past_year_bot_issues_count"`
		PastYearBotPullRequestsCount               *int                         `json:"past_year_bot_pull_requests_count"`
		PastYearMergedPullRequestsCount            *int                         `json:"past_year_merged_pull_requests_count"`
		CreatedAt                                  time.Time                    `json:"created_at"`
		UpdatedAt                                  time.Time                    `json:"updated_at"`
		RepositoryUrl                              *string                      `json:"repository_url"`
		IssuesUrl                                  *string                      `json:"issues_url"`
		IssueLabelsCount                           EcosystemsCountTable         `json:"issue_labels_count"`
		PullRequestLabelsCount                     EcosystemsCountTable         `json:"pull_request_labels_count"`
		IssueAuthorAssociationsCount               EcosystemsCountTable         `json:"issue_author_associations_count"`
		PullRequestAuthorAssociationsCount         EcosystemsCountTable         `json:"pull_request_author_associations_count"`
		IssueAuthors                               EcosystemsCountTable         `json:"issue_authors"`
		PullRequestAuthors                         EcosystemsCountTable         `json:"pull_request_authors"`
		Host                                       EcosystemsHost               `json:"host"`
		PastYearIssueLabelsCount                   EcosystemsCountTable         `json:"past_year_issue_labels_count"`
		PastYearPullRequestLabelsCount             EcosystemsCountTable         `json:"past_year_pull_request_labels_count"`
		PastYearIssueAuthorAssociationsCount       EcosystemsCountTable         `json:"past_year_issue_author_associations_count"`
		PastYearPullRequestAuthorAssociationsCount EcosystemsCountTable         `json:"past_year_pull_request_author_associations_count"`
		PastYearIssueAuthors                       EcosystemsCountTable         `json:"past_year_issue_authors"`
		PastYearPullRequestAuthors                 EcosystemsCountTable         `json:"past_year_pull_request_authors"`
		Maintainers                                []EcosystemsIssueAuthorTable `json:"maintainers"`
		ActiveMaintainers                          []EcosystemsIssueAuthorTable `json:"active_maintainers"`
	} `json:"table"`
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
	Issues  EcosystemsIssues  `json:"issues"`
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

// See `getMaintainersFromCommits()` for more information.
const SIGNIFICANT_COMMITTER_FRACTION float32 = 0.3

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

func mergeMaintainers(a EcosystemsMaintainer, b EcosystemsMaintainer) EcosystemsMaintainer {
	if a.Login == nil {
		a.Login = b.Login
	}
	if a.Name == nil {
		a.Name = b.Name
	}
	if a.Email == nil {
		a.Email = b.Email
	}
	return a
}

func getMaintainersFromIssues(issues EcosystemsIssues) []EcosystemsMaintainer {
	maintainers := []EcosystemsMaintainer{}
	for _, issueRaiser := range issues.Table.Maintainers {
		newMaintainer := EcosystemsMaintainer{
			Login: issueRaiser.Table.Login,
			Name:  nil,
			Email: nil,
		}
		maintainers = append(maintainers, newMaintainer)
	}
	return maintainers
}

// Get all committers for a package, sorted descendingly by commit count. For
// each committer, mark them a significant commiter, and continue until all
// significant commiters have together cumulatively committed at least
// `SIGNIFICANT_COMMITTER_FRACTION` of all of the project's commits.
func getMaintainersFromCommits(commits EcosystemsCommits) []EcosystemsMaintainer {
	maintainers := []EcosystemsMaintainer{}
	nCommitsCounted := 0
	// NOTE: This works because committers are descendingly sorted by commit
	// count in the ecosyste.ms API response. If this stops being the case,
	// this code will break.
	for _, committer := range commits.Committers {
		nCommitsCounted += committer.Count
		newMaintainer := EcosystemsMaintainer{
			Login: committer.Login,
			Name:  committer.Name,
			Email: committer.Email,
		}
		maintainers = append(maintainers, newMaintainer)
		commitFractionCounted := float32(nCommitsCounted) / float32(commits.TotalCommits)
		if commitFractionCounted >= SIGNIFICANT_COMMITTER_FRACTION {
			break
		}
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

func mergeMaintainerLists(a []EcosystemsMaintainer, b []EcosystemsMaintainer) []EcosystemsMaintainer {
	for _, mb := range b {
		exists := false
		for idx, ma := range a {
			isSimilar := (ma.Login != nil && mb.Login != nil && *ma.Login == *mb.Login) ||
				(ma.Name != nil && mb.Name != nil && *ma.Name == *mb.Name) ||
				(ma.Email != nil && mb.Email != nil && *ma.Email == *mb.Email)
			if isSimilar {
				exists = true
				a[idx] = mergeMaintainers(ma, mb)
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

		// log.Println("##### Initial maintainers #####")
		// spew.Fdump(os.Stderr, pkg.Maintainers)

		maintainersFromCommits := getMaintainersFromCommits(summary.Commits)
		maintainersFromIssues := getMaintainersFromIssues(summary.Issues)
		pkg.Maintainers = mergeMaintainerLists(pkg.Maintainers, maintainersFromCommits)
		pkg.Maintainers = mergeMaintainerLists(pkg.Maintainers, maintainersFromIssues)

		// log.Println("##### Maintainers from commits #####")
		// spew.Fdump(os.Stderr, maintainersFromCommits)
		// log.Println("##### Maintainers from issues #####")
		// spew.Fdump(os.Stderr, maintainersFromIssues)
		// log.Println("##### Final maintainers #####")
		// spew.Fdump(os.Stderr, pkg.Maintainers)
		// log.Printf("\n\n\n")
	}
	maintainerMap := getMaintainerMapFromPackages(packages)
	log.Println(len(packages))
	for k := range maintainerMap {
		log.Printf("%03d\t%s; %s; %s\n",
			maintainerMap[k].NCriticalPackages,
			k.Login, k.Name, k.Email)
	}
}
