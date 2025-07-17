package models

import "time"

// SubmissionStatus represents the status of a submission
type SubmissionStatus string

const (
	StatusPending             SubmissionStatus = "pending"
	StatusRunning             SubmissionStatus = "running"
	StatusAccepted            SubmissionStatus = "accepted"
	StatusRejected            SubmissionStatus = "rejected"
	StatusError               SubmissionStatus = "error"
	StatusTimeLimitExceeded   SubmissionStatus = "time_limit_exceeded"
	StatusMemoryLimitExceeded SubmissionStatus = "memory_limit_exceeded"
	StatusCompileError        SubmissionStatus = "compile_error"
	StatusRuntimeError        SubmissionStatus = "runtime_error"
)

// TestCaseResult represents the result of a single test case
type TestCaseResult struct {
	TestCaseID     string           `json:"test_case_id" firestore:"test_case_id"`
	Status         SubmissionStatus `json:"status" firestore:"status"`
	ExecutionTime  int64            `json:"execution_time" firestore:"execution_time"` // in milliseconds
	MemoryUsed     int64            `json:"memory_used" firestore:"memory_used"`       // in kilobytes
	ExpectedOutput string           `json:"expected_output" firestore:"expected_output"`
	ActualOutput   string           `json:"actual_output" firestore:"actual_output"`
	ErrorMessage   string           `json:"error_message,omitempty" firestore:"error_message,omitempty"`
}

// Submission represents a solution submission
type Submission struct {
	ID            string           `json:"id" firestore:"id"`
	UserID        string           `json:"user_id" firestore:"user_id"`
	ProblemID     string           `json:"problem_id" firestore:"problem_id"`
	Language      string           `json:"language" firestore:"language"`
	Code          string           `json:"code" firestore:"code"`
	Status        SubmissionStatus `json:"status" firestore:"status"`
	SubmittedAt   time.Time        `json:"submitted_at" firestore:"submitted_at"`
	CompletedAt   time.Time        `json:"completed_at,omitempty" firestore:"completed_at,omitempty"`
	ExecutionTime int64            `json:"execution_time" firestore:"execution_time"` // in milliseconds
	MemoryUsed    int64            `json:"memory_used" firestore:"memory_used"`       // in kilobytes
	TestCases     []TestCaseResult `json:"test_cases" firestore:"test_cases"`
	ErrorMessage  string           `json:"error_message,omitempty" firestore:"error_message,omitempty"`
}

// TestCase represents a test case for a problem
type TestCase struct {
	ID       string `json:"id" firestore:"id"`
	Input    string `json:"input" firestore:"input"`
	Expected string `json:"expected" firestore:"expected"`
	Hidden   bool   `json:"hidden" firestore:"hidden"`
	Weight   int    `json:"weight" firestore:"weight"` // For weighted scoring
}

// SubmissionSummary represents a summary of a user's submissions for a problem
type SubmissionSummary struct {
	ProblemID        string    `json:"problem_id" firestore:"problem_id"`
	Language         string    `json:"language" firestore:"language"`
	LastSubmissionID string    `json:"last_submission_id" firestore:"last_submission_id"`
	BestSubmissionID string    `json:"best_submission_id" firestore:"best_submission_id"`
	AttemptCount     int       `json:"attempt_count" firestore:"attempt_count"`
	Solved           bool      `json:"solved" firestore:"solved"`
	FirstSolvedAt    time.Time `json:"first_solved_at,omitempty" firestore:"first_solved_at,omitempty"`
	BestTime         int64     `json:"best_time" firestore:"best_time"`     // in milliseconds
	BestMemory       int64     `json:"best_memory" firestore:"best_memory"` // in kilobytes
	LastAttemptedAt  time.Time `json:"last_attempted_at" firestore:"last_attempted_at"`
}
