package models

type UserRanking struct {
	UserID         string  `json:"user_id" firestore:"user_id"`
	Username       string  `json:"username" firestore:"username"`
	Score          float64 `json:"score" firestore:"score"`
	ProblemsSolved int     `json:"problems_solved" firestore:"problems_solved"`
	TotalTime      int64   `json:"total_time" firestore:"total_time"`     // in milliseconds
	TotalMemory    int64   `json:"total_memory" firestore:"total_memory"` // in kilobytes
	Rank           int     `json:"rank" firestore:"rank"`
}
