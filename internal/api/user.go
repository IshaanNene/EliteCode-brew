package api

import (
	"fmt"
	"time"
)

type UserStats struct {
	ProblemsSolved    int                    `json:"problemsSolved"`
	TotalSubmissions  int                    `json:"totalSubmissions"`
	Accuracy          float64                `json:"accuracy"`
	CurrentStreak     int                    `json:"currentStreak"`
	MaxStreak         int                    `json:"maxStreak"`
	Languages         map[string]int         `json:"languages"`
	DifficultyStats   map[string]int         `json:"difficultyStats"`
	TagStats          map[string]int         `json:"tagStats"`
	MonthlyProgress   []MonthlyProgress      `json:"monthlyProgress"`
	RecentActivity    []Activity             `json:"recentActivity"`
	Achievements      []Achievement          `json:"achievements"`
}

type MonthlyProgress struct {
	Month  string `json:"month"`
	Solved int    `json:"solved"`
}

type Activity struct {
	Type        string    `json:"type"` // submission, bookmark, achievement
	ProblemID   string    `json:"problemId"`
	ProblemName string    `json:"problemName"`
	Status      string    `json:"status"`
	Language    string    `json:"language"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Achievement struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	UnlockedAt  time.Time `json:"unlockedAt"`
}

type UserProblems struct {
	Solved     []ProblemSummary `json:"solved"`
	Attempted  []ProblemSummary `json:"attempted"`
	Bookmarked []ProblemSummary `json:"bookmarked"`
	Total      UserProblemCount `json:"total"`
}

type UserProblemCount struct {
	Solved     int `json:"solved"`
	Attempted  int `json:"attempted"`
	Bookmarked int `json:"bookmarked"`
}

type BookmarkRequest struct {
	ProblemID string `json:"problemId"`
}

type BookmarkResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Bookmarked bool  `json:"bookmarked"`
}

func (c *Client) GetUserStats() (*UserStats, error) {
	var stats UserStats
	err := c.Get("/users/stats", &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (c *Client) GetUserProblems(status string, page, pageSize int) (*UserProblems, error) {
	var problems UserProblems
	endpoint := fmt.Sprintf("/users/problems?status=%s&page=%d&pageSize=%d", status, page, pageSize)
	err := c.Get(endpoint, &problems)
	if err != nil {
		return nil, err
	}
	return &problems, nil
}

func (c *Client) BookmarkProblem(problemId string) (*BookmarkResponse, error) {
	var response BookmarkResponse
	req := BookmarkRequest{ProblemID: problemId}
	err := c.Post("/users/bookmarks", req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) UnbookmarkProblem(problemId string) (*BookmarkResponse, error) {
	var response BookmarkResponse
	endpoint := fmt.Sprintf("/users/bookmarks/%s", problemId)
	err := c.Delete(endpoint, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) GetBookmarks(page, pageSize int) ([]ProblemSummary, error) {
	var bookmarks []ProblemSummary
	endpoint := fmt.Sprintf("/users/bookmarks?page=%d&pageSize=%d", page, pageSize)
	err := c.Get(endpoint, &bookmarks)
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (c *Client) GetUserActivity(page, pageSize int) ([]Activity, error) {
	var activities []Activity
	endpoint := fmt.Sprintf("/users/activity?page=%d&pageSize=%d", page, pageSize)
	err := c.Get(endpoint, &activities)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (c *Client) GetUserAchievements() ([]Achievement, error) {
	var achievements []Achievement
	err := c.Get("/users/achievements", &achievements)
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

func (c *Client) UpdateUserProfile(updates map[string]interface{}) (*User, error) {
	var user User
	err := c.Put("/users/profile", updates, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}