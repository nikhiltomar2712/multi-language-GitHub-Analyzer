package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// Repository represents GitHub repository data
type Repository struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Stars        int       `json:"stargazers_count"`
	Forks        int       `json:"forks_count"`
	Watchers     int       `json:"watchers_count"`
	Language     string    `json:"language"`
	Topics       []string  `json:"topics"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	OpenIssues   int       `json:"open_issues_count"`
	URL          string    `json:"html_url"`
	IsArchived   bool      `json:"archived"`
	IsForked     bool      `json:"fork"`
}

// Contributor represents a repository contributor
type Contributor struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
	AvatarURL     string `json:"avatar_url"`
}

// Languages represents programming language distribution
type Languages map[string]int

// HealthScore represents repository health assessment
type HealthScore struct {
	Score   int
	Rating  string
	Insights []string
}

// AnalysisResult contains complete analysis output
type AnalysisResult struct {
	Repository  Repository
	Contributors []Contributor
	Languages   Languages
	Health      HealthScore
	Duration    time.Duration
}

// GitHubClient handles GitHub API interactions
type GitHubClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		baseURL: "https://api.github.com",
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest performs an HTTP request to GitHub API
func (gc *GitHubClient) makeRequest(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", gc.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if gc.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", gc.token))
	}

	resp, err := gc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetRepository fetches repository information
func (gc *GitHubClient) GetRepository(owner, repo string) (*Repository, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s", owner, repo)
	body, err := gc.makeRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var repository Repository
	if err := json.Unmarshal(body, &repository); err != nil {
		return nil, err
	}

	return &repository, nil
}

// GetContributors fetches repository contributors
func (gc *GitHubClient) GetContributors(owner, repo string) ([]Contributor, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/contributors?per_page=100", owner, repo)
	body, err := gc.makeRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var contributors []Contributor
	if err := json.Unmarshal(body, &contributors); err != nil {
		return nil, err
	}

	return contributors, nil
}

// GetLanguages fetches programming language distribution
func (gc *GitHubClient) GetLanguages(owner, repo string) (Languages, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/languages", owner, repo)
	body, err := gc.makeRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var languages Languages
	if err := json.Unmarshal(body, &languages); err != nil {
		return nil, err
	}

	return languages, nil
}

// CalculateHealth evaluates repository health
func CalculateHealth(repo *Repository, contributors []Contributor) HealthScore {
	score := 0
	insights := []string{}

	// Stars evaluation (0-25)
	if repo.Stars >= 1000 {
		score += 25
		insights = append(insights, "🌟 Excellent community engagement")
	} else if repo.Stars >= 500 {
		score += 20
		insights = append(insights, "⭐ Good community interest")
	} else if repo.Stars >= 100 {
		score += 15
		insights = append(insights, "📈 Growing community")
	}

	// Contributors evaluation (0-25)
	if len(contributors) >= 50 {
		score += 25
		insights = append(insights, "👥 Strong contributor base")
	} else if len(contributors) >= 20 {
		score += 20
		insights = append(insights, "🤝 Active community")
	} else if len(contributors) >= 5 {
		score += 15
		insights = append(insights, "👤 Growing team")
	}

	// Forks evaluation (0-20)
	forkRatio := float64(repo.Forks) / float64(repo.Stars)
	if repo.Stars > 0 && forkRatio >= 0.3 {
		score += 20
		insights = append(insights, "🔀 High fork-to-star ratio")
	} else if repo.Stars > 0 && forkRatio >= 0.1 {
		score += 15
		insights = append(insights, "🔀 Moderate fork activity")
	}

	// Issues evaluation (0-15)
	if repo.OpenIssues < 20 {
		score += 15
		insights = append(insights, "✅ Well-maintained issue tracker")
	} else if repo.OpenIssues < 50 {
		score += 10
		insights = append(insights, "📋 Manageable issue backlog")
	}

	// Activity evaluation (0-15)
	updatedAt, _ := time.Parse("2006-01-02T15:04:05Z", repo.UpdatedAt)
	daysSinceUpdate := time.Since(updatedAt).Hours() / 24

	if daysSinceUpdate < 7 {
		score += 15
		insights = append(insights, "🔥 Actively maintained")
	} else if daysSinceUpdate < 30 {
		score += 10
		insights = append(insights, "✓ Regularly maintained")
	} else if daysSinceUpdate < 90 {
		score += 5
		insights = append(insights, "⏱️ Maintenance may be slowing")
	}

	rating := "Needs Improvement"
	if score >= 90 {
		rating = "Excellent"
	} else if score >= 70 {
		rating = "Good"
	} else if score >= 50 {
		rating = "Fair"
	}

	return HealthScore{
		Score:    score,
		Rating:   rating,
		Insights: insights,
	}
}

// Analyze performs concurrent analysis of a GitHub repository
func (gc *GitHubClient) Analyze(owner, repo string) (*AnalysisResult, error) {
	startTime := time.Now()

	// Use WaitGroup for concurrent requests
	var wg sync.WaitGroup
	var mu sync.Mutex

	result := &AnalysisResult{}
	var repoErr, contribErr, langErr error

	// Fetch repository info
	wg.Add(1)
	go func() {
		defer wg.Done()
		r, err := gc.GetRepository(owner, repo)
		if err != nil {
			mu.Lock()
			repoErr = err
			mu.Unlock()
			return
		}
		mu.Lock()
		result.Repository = *r
		mu.Unlock()
	}()

	// Fetch contributors
	wg.Add(1)
	go func() {
		defer wg.Done()
		c, err := gc.GetContributors(owner, repo)
		if err != nil {
			mu.Lock()
			contribErr = err
			mu.Unlock()
			return
		}
		mu.Lock()
		result.Contributors = c
		mu.Unlock()
	}()

	// Fetch languages
	wg.Add(1)
	go func() {
		defer wg.Done()
		l, err := gc.GetLanguages(owner, repo)
		if err != nil {
			mu.Lock()
			langErr = err
			mu.Unlock()
			return
		}
		mu.Lock()
		result.Languages = l
		mu.Unlock()
	}()

	wg.Wait()

	// Check for errors
	if repoErr != nil {
		return nil, repoErr
	}
	if contribErr != nil {
		return nil, contribErr
	}
	if langErr != nil {
		return nil, langErr
	}

	// Calculate health
	result.Health = CalculateHealth(&result.Repository, result.Contributors)
	result.Duration = time.Since(startTime)

	return result, nil
}

// PrintResults outputs analysis results in a formatted way
func PrintResults(result *AnalysisResult) {
	fmt.Println("\n╔════════════════════════════════════════════════════════╗")
	fmt.Println("║           GITHUB REPOSITORY ANALYSIS REPORT            ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝\n")

	fmt.Printf("📦 Repository: %s\n", result.Repository.Name)
	if result.Repository.Description != "" {
		fmt.Printf("📝 Description: %s\n", result.Repository.Description)
	}

	fmt.Println("\n📊 Statistics:")
	fmt.Printf("  ⭐ Stars: %d\n", result.Repository.Stars)
	fmt.Printf("  🔀 Forks: %d\n", result.Repository.Forks)
	fmt.Printf("  👁️ Watchers: %d\n", result.Repository.Watchers)
	fmt.Printf("  ⚠️ Open Issues: %d\n", result.Repository.OpenIssues)

	fmt.Printf("\n👥 Top Contributors: %d total\n", len(result.Contributors))
	topN := 5
	if len(result.Contributors) < topN {
		topN = len(result.Contributors)
	}
	for i := 0; i < topN; i++ {
		c := result.Contributors[i]
		fmt.Printf("  %d. %s (%d commits)\n", i+1, c.Login, c.Contributions)
	}

	// Sort and display languages
	type langPair struct {
		lang  string
		bytes int
	}
	var pairs []langPair
	for lang, bytes := range result.Languages {
		pairs = append(pairs, langPair{lang, bytes})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].bytes > pairs[j].bytes
	})

	fmt.Println("\n💻 Languages:")
	totalBytes := 0
	for _, p := range pairs {
		totalBytes += p.bytes
	}
	for i := 0; i < 5 && i < len(pairs); i++ {
		p := pairs[i]
		percent := (float64(p.bytes) / float64(totalBytes)) * 100
		fmt.Printf("  %s: %.1f%%\n", p.lang, percent)
	}

	fmt.Println("\n🏥 Health Assessment:")
	fmt.Printf("  Health Score: %d/100\n", result.Health.Score)
	fmt.Printf("  Rating: %s\n", result.Health.Rating)
	fmt.Println("\n  Insights:")
	for _, insight := range result.Health.Insights {
		fmt.Printf("  %s\n", insight)
	}

	fmt.Printf("\n⏱️ Analysis completed in %v\n\n", result.Duration)
}

func main() {
	token := flag.String("token", "", "GitHub API token")
	help := flag.Bool("help", false, "Show help message")
	flag.Bool("h", false, "Show help message")

	flag.Parse()

	if *help || len(os.Args) < 2 {
		fmt.Println(`
GitHub Analyzer - Go Version
High-performance concurrent repository analysis

Usage: github-analyzer-go <owner/repo> [options]

Options:
  -token STRING     GitHub API token (or set GITHUB_TOKEN environment variable)
  -help, -h         Show this help message

Example:
  ./github-analyzer-go torvalds/linux
  ./github-analyzer-go facebook/react -token ghp_xxxxx
		`)
		os.Exit(0)
	}

	repoPath := flag.Arg(0)
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		fmt.Println("❌ Error: Invalid repository format. Use: owner/repo")
		os.Exit(1)
	}

	owner, repo := parts[0], parts[1]

	// Get token from flag or environment
	if *token == "" {
		*token = os.Getenv("GITHUB_TOKEN")
	}

	client := NewGitHubClient(*token)
	fmt.Printf("\n📊 Analyzing %s/%s...\n", owner, repo)

	result, err := client.Analyze(owner, repo)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	PrintResults(result)
}
