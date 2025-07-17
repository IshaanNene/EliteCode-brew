package models

import "time"

type Difficulty string

const (
	Easy     Difficulty = "easy"
	Medium   Difficulty = "medium"
	Hard     Difficulty = "hard"
	VeryHard Difficulty = "very_hard"
)

type StoragePaths struct {
	StarterCode string `json:"starter_code" firestore:"starter_code"`
	TestCases   string `json:"test_cases" firestore:"test_cases"`
}

type Problem struct {
	ID             string       `json:"id" firestore:"id"`
	Title          string       `json:"title" firestore:"title"`
	Description    string       `json:"description" firestore:"description"`
	Difficulty     Difficulty   `json:"difficulty" firestore:"difficulty"`
	Tags           []string     `json:"tags" firestore:"tags"`
	TimeLimit      int          `json:"time_limit" firestore:"time_limit"`     // in milliseconds
	MemoryLimit    int          `json:"memory_limit" firestore:"memory_limit"` // in megabytes
	SupportedLangs []string     `json:"supported_langs" firestore:"supported_langs"`
	StoragePaths   StoragePaths `json:"storage_paths" firestore:"storage_paths"`
	CreatedAt      time.Time    `json:"created_at" firestore:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" firestore:"updated_at"`
}

type ProblemStats struct {
	TotalSubmissions int       `json:"total_submissions" firestore:"total_submissions"`
	AcceptedCount    int       `json:"accepted_count" firestore:"accepted_count"`
	AcceptanceRate   float64   `json:"acceptance_rate" firestore:"acceptance_rate"`
	AverageTime      int64     `json:"average_time" firestore:"average_time"`     // in milliseconds
	AverageMemory    int64     `json:"average_memory" firestore:"average_memory"` // in kilobytes
	UpdatedAt        time.Time `json:"updated_at" firestore:"updated_at"`
}

type UserProblemStatus struct {
	UserID            string             `json:"user_id" firestore:"user_id"`
	ProblemID         string             `json:"problem_id" firestore:"problem_id"`
	AttemptCount      int                `json:"attempt_count" firestore:"attempt_count"`
	Bookmarked        bool               `json:"bookmarked" firestore:"bookmarked"`
	LastAttemptedAt   time.Time          `json:"last_attempted_at" firestore:"last_attempted_at"`
	SubmissionSummary *SubmissionSummary `json:"submission_summary" firestore:"submission_summary"`
}
