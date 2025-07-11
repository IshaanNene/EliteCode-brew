package api

import (
	"fmt"
	"time"
)

type Problem struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Difficulty  string            `json:"difficulty"`
	Tags        []string          `json:"tags"`
	Examples    []Example         `json:"examples"`
	Constraints string            `json:"constraints"`
	Languages   []string          `json:"languages"`
	TestCases   []TestCase        `json:"testCases"`
	Files       map[string]string `json:"files"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type Example struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation"`
}

type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
	Hidden   bool   `json:"hidden"`
}

type ProblemList struct {
	Problems   []ProblemSummary `json:"problems"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalPages int              `json:"totalPages"`
}

type ProblemSummary struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Difficulty string    `json:"difficulty"`
	Tags       []string  `json:"tags"`
	Solved     bool      `json:"solved"`
	Attempted  bool      `json:"attempted"`
	CreatedAt  time.Time `json:"createdAt"`
}

type SearchProblemsRequest struct {
	Query      string   `json:"query"`
	Tags       []string `json:"tags"`
	Difficulty string   `json:"difficulty"`
	Status     string   `json:"status"` // all, solved, unsolved, attempted
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
}

type SubmissionRequest struct {
	ProblemID string `json:"problemId"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}

type SubmissionResponse struct {
	ID        string           `json:"id"`
	Status    string           `json:"status"`
	Results   []TestResult     `json:"results"`
	Stats     SubmissionStats  `json:"stats"`
	CreatedAt time.Time        `json:"createdAt"`
}

type TestResult struct {
	TestCase int     `json:"testCase"`
	Status   string  `json:"status"` // passed, failed, timeout, error
	Input    string  `json:"input"`
	Expected string  `json:"expected"`
	Output   string  `json:"output"`
	Error    string  `json:"error"`
	Time     float64 `json:"time"`
	Memory   int64   `json:"memory"`
}

type SubmissionStats struct {
	TotalTests   int     `json:"totalTests"`
	PassedTests  int     `json:"passedTests"`
	FailedTests  int     `json:"failedTests"`
	TotalTime    float64 `json:"totalTime"`
	MaxMemory    int64   `json:"maxMemory"`
	Accuracy     float64 `json:"accuracy"`
}

type RunRequest struct {
	ProblemID string `json:"problemId"`
	Language  string `json:"language"`
	Code      string `json:"code"`
	Input     string `json:"input"`
}

type RunResponse struct {
	Output    string  `json:"output"`
	Error     string  `json:"error"`
	Time      float64 `json:"time"`
	Memory    int64   `json:"memory"`
	Status    string  `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (c *Client) GetProblems(page, pageSize int) (*ProblemList, error) {
	var problems ProblemList
	endpoint := fmt.Sprintf("/problems?page=%d&pageSize=%d", page, pageSize)
	err := c.Get(endpoint, &problems)
	if err != nil {
		return nil, err
	}
	return &problems, nil
}

func (c *Client) GetProblem(id string) (*Problem, error) {
	var problem Problem
	endpoint := fmt.Sprintf("/problems/%s", id)
	err := c.Get(endpoint, &problem)
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

func (c *Client) SearchProblems(req SearchProblemsRequest) (*ProblemList, error) {
	var problems ProblemList
	err := c.Post("/problems/search", req, &problems)
	if err != nil {
		return nil, err
	}
	return &problems, nil
}

func (c *Client) SubmitSolution(req SubmissionRequest) (*SubmissionResponse, error) {
	var response SubmissionResponse
	err := c.Post("/problems/submit", req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) RunCode(req RunRequest) (*RunResponse, error) {
	var response RunResponse
	err := c.Post("/problems/run", req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) GetSubmission(id string) (*SubmissionResponse, error) {
	var submission SubmissionResponse
	endpoint := fmt.Sprintf("/problems/submissions/%s", id)
	err := c.Get(endpoint, &submission)
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func (c *Client) GetProblemSubmissions(problemId string, page, pageSize int) ([]SubmissionResponse, error) {
	var submissions []SubmissionResponse
	endpoint := fmt.Sprintf("/problems/%s/submissions?page=%d&pageSize=%d", problemId, page, pageSize)
	err := c.Get(endpoint, &submissions)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}