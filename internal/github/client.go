package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub client with additional functionality
type Client struct {
	client *github.Client
	owner  string
	repo   string
}

// NewClient creates a new GitHub client
func NewClient(token string) (*Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
	}, nil
}

// GetUser gets the authenticated user's information
func (c *Client) GetUser(ctx context.Context) (*github.User, error) {
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}

// SetRepository sets the owner and repository for the client
func (c *Client) SetRepository(owner, repo string) {
	c.owner = owner
	c.repo = repo
}

// CreateRepository creates a new GitHub repository
func (c *Client) CreateRepository(ctx context.Context, name string, description string, private bool) (*github.Repository, error) {
	repo := &github.Repository{
		Name:        github.String(name),
		Description: github.String(description),
		Private:     github.Bool(private),
		AutoInit:    github.Bool(true),
	}

	repository, _, err := c.client.Repositories.Create(ctx, "", repo)
	if err != nil {
		return nil, fmt.Errorf("error creating repository: %v", err)
	}

	return repository, nil
}

// GetRepository gets the repository information
func (c *Client) GetRepository(ctx context.Context) (*github.Repository, error) {
	repo, _, err := c.client.Repositories.Get(ctx, c.owner, c.repo)
	if err != nil {
		return nil, fmt.Errorf("error getting repository: %v", err)
	}

	return repo, nil
}

// CreatePullRequest creates a new pull request
func (c *Client) CreatePullRequest(ctx context.Context, title, body, head, base string) (*github.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Body:  github.String(body),
		Head:  github.String(head),
		Base:  github.String(base),
	}

	pr, _, err := c.client.PullRequests.Create(ctx, c.owner, c.repo, newPR)
	if err != nil {
		return nil, fmt.Errorf("error creating pull request: %v", err)
	}

	return pr, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, title, body string, labels []string) (*github.Issue, error) {
	issue := &github.IssueRequest{
		Title:  github.String(title),
		Body:   github.String(body),
		Labels: &labels,
	}

	newIssue, _, err := c.client.Issues.Create(ctx, c.owner, c.repo, issue)
	if err != nil {
		return nil, fmt.Errorf("error creating issue: %v", err)
	}

	return newIssue, nil
}

// GetProblemBranch gets or creates a branch for a problem
func (c *Client) GetProblemBranch(ctx context.Context, problemID string) (string, error) {
	branchName := fmt.Sprintf("problem/%s", problemID)

	// Get default branch
	repo, err := c.GetRepository(ctx)
	if err != nil {
		return "", err
	}

	defaultBranch := repo.GetDefaultBranch()

	// Get default branch reference
	ref, _, err := c.client.Git.GetRef(ctx, c.owner, c.repo, "refs/heads/"+defaultBranch)
	if err != nil {
		return "", fmt.Errorf("error getting default branch ref: %v", err)
	}

	// Check if problem branch exists
	_, _, err = c.client.Git.GetRef(ctx, c.owner, c.repo, "refs/heads/"+branchName)
	if err == nil {
		// Branch exists
		return branchName, nil
	}

	// Create new branch
	newRef := &github.Reference{
		Ref:    github.String("refs/heads/" + branchName),
		Object: &github.GitObject{SHA: ref.Object.SHA},
	}

	_, _, err = c.client.Git.CreateRef(ctx, c.owner, c.repo, newRef)
	if err != nil {
		return "", fmt.Errorf("error creating branch: %v", err)
	}

	return branchName, nil
}

// CommitSolution commits a solution to a problem
func (c *Client) CommitSolution(ctx context.Context, problemID, title string, files map[string][]byte) error {
	// Get problem branch
	branch, err := c.GetProblemBranch(ctx, problemID)
	if err != nil {
		return err
	}

	// Get branch reference
	ref, _, err := c.client.Git.GetRef(ctx, c.owner, c.repo, "refs/heads/"+branch)
	if err != nil {
		return fmt.Errorf("error getting branch ref: %v", err)
	}

	// Get base tree
	baseTree, _, err := c.client.Git.GetTree(ctx, c.owner, c.repo, *ref.Object.SHA, false)
	if err != nil {
		return fmt.Errorf("error getting base tree: %v", err)
	}

	// Create blobs for each file
	var entries []*github.TreeEntry
	for path, content := range files {
		blob := &github.Blob{
			Content:  github.String(string(content)),
			Encoding: github.String("utf-8"),
		}

		newBlob, _, err := c.client.Git.CreateBlob(ctx, c.owner, c.repo, blob)
		if err != nil {
			return fmt.Errorf("error creating blob for %s: %v", path, err)
		}

		entries = append(entries, &github.TreeEntry{
			Path: github.String(path),
			Mode: github.String("100644"),
			Type: github.String("blob"),
			SHA:  newBlob.SHA,
		})
	}

	// Create new tree
	tree, _, err := c.client.Git.CreateTree(ctx, c.owner, c.repo, *baseTree.SHA, entries)
	if err != nil {
		return fmt.Errorf("error creating tree: %v", err)
	}

	// Get parent commit
	parent, _, err := c.client.Git.GetCommit(ctx, c.owner, c.repo, *ref.Object.SHA)
	if err != nil {
		return fmt.Errorf("error getting parent commit: %v", err)
	}

	// Create commit
	now := github.Timestamp{Time: time.Now()}
	author := &github.CommitAuthor{
		Date:  &now,
		Name:  github.String("Elitecode"),
		Email: github.String("elitecode@example.com"),
	}

	commit := &github.Commit{
		Author:  author,
		Message: github.String(fmt.Sprintf("Solution for problem %s: %s", problemID, title)),
		Tree:    tree,
		Parents: []*github.Commit{parent},
	}

	newCommit, _, err := c.client.Git.CreateCommit(ctx, c.owner, c.repo, commit, &github.CreateCommitOptions{})
	if err != nil {
		return fmt.Errorf("error creating commit: %v", err)
	}

	// Update branch reference
	ref.Object.SHA = newCommit.SHA
	_, _, err = c.client.Git.UpdateRef(ctx, c.owner, c.repo, ref, false)
	if err != nil {
		return fmt.Errorf("error updating branch ref: %v", err)
	}

	return nil
}

// GetGitHubToken gets the GitHub token from the config file
func GetGitHubToken() (string, error) {
	configFile := filepath.Join(os.Getenv("HOME"), ".elitecode", "config.json")
	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("error reading config file: %v", err)
	}

	var config struct {
		GitHubToken string `json:"github_token"`
	}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return "", fmt.Errorf("error parsing config: %v", err)
	}

	if config.GitHubToken == "" {
		return "", fmt.Errorf("GitHub token not found in config. Please run 'elitecode github login' first")
	}

	return config.GitHubToken, nil
}

// GetRepositoryInfo gets the repository owner and name from the git config
func GetRepositoryInfo() (owner, repo string, err error) {
	// Run git remote -v
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("error getting git remote: %v", err)
	}

	// Parse remote URL
	remoteURL := strings.TrimSpace(string(output))
	remoteURL = strings.TrimSuffix(remoteURL, ".git")

	// Handle SSH and HTTPS URLs
	var parts []string
	if strings.HasPrefix(remoteURL, "git@github.com:") {
		parts = strings.Split(strings.TrimPrefix(remoteURL, "git@github.com:"), "/")
	} else if strings.HasPrefix(remoteURL, "https://github.com/") {
		parts = strings.Split(strings.TrimPrefix(remoteURL, "https://github.com/"), "/")
	} else {
		return "", "", fmt.Errorf("unsupported git remote URL format")
	}

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid git remote URL format")
	}

	return parts[0], parts[1], nil
}
