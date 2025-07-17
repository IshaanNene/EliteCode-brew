package problem

import (
	"context"
	"fmt"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/IshaanNene/EliteCode-brew/internal/models"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	firestoreClient *firestore.Client
}

func NewService(firestoreClient *firestore.Client) *Service {
	return &Service{
		firestoreClient: firestoreClient,
	}
}

func (s *Service) ListProblems(ctx context.Context, filters map[string]interface{}) ([]models.Problem, error) {
	collRef := s.firestoreClient.Collection("problems")

	var query firestore.Query
	first := true
	for key, value := range filters {
		if first {
			query = collRef.Where(key, "==", value)
			first = false
		} else {
			query = query.Where(key, "==", value)
		}
	}

	var docs []*firestore.DocumentSnapshot
	var err error
	if first {
		docs, err = collRef.Documents(ctx).GetAll()
	} else {
		docs, err = query.Documents(ctx).GetAll()
	}

	if err != nil {
		return nil, fmt.Errorf("error fetching problems: %v", err)
	}

	problems := make([]models.Problem, 0, len(docs))
	for _, doc := range docs {
		var problem models.Problem
		if err := doc.DataTo(&problem); err != nil {
			return nil, fmt.Errorf("error parsing problem data: %v", err)
		}
		problem.ID = doc.Ref.ID
		problems = append(problems, problem)
	}

	sort.Slice(problems, func(i, j int) bool {
		if problems[i].Difficulty != problems[j].Difficulty {
			difficultyOrder := map[models.Difficulty]int{
				models.Easy:     1,
				models.Medium:   2,
				models.Hard:     3,
				models.VeryHard: 4,
			}
			return difficultyOrder[problems[i].Difficulty] < difficultyOrder[problems[j].Difficulty]
		}
		return problems[i].ID < problems[j].ID
	})

	return problems, nil
}

func (s *Service) GetProblem(ctx context.Context, id string) (*models.Problem, error) {
	doc, err := s.firestoreClient.Collection("problems").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching problem: %v", err)
	}

	var problem models.Problem
	if err := doc.DataTo(&problem); err != nil {
		return nil, fmt.Errorf("error parsing problem data: %v", err)
	}
	problem.ID = doc.Ref.ID

	return &problem, nil
}

func (s *Service) GetUserProblemStatus(ctx context.Context, userID, problemID string) (*models.UserProblemStatus, error) {
	doc, err := s.firestoreClient.Collection("users").Doc(userID).Collection("problem_status").Doc(problemID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching user problem status: %v", err)
	}

	var status models.UserProblemStatus
	if err := doc.DataTo(&status); err != nil {
		return nil, fmt.Errorf("error parsing user problem status: %v", err)
	}

	return &status, nil
}

func (s *Service) GetUserProblemStatuses(ctx context.Context, userID string) ([]models.UserProblemStatus, error) {
	iter := s.firestoreClient.Collection("users").Doc(userID).Collection("problem_status").Documents(ctx)
	defer iter.Stop()

	var statuses []models.UserProblemStatus
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating problem statuses: %v", err)
		}

		var status models.UserProblemStatus
		if err := doc.DataTo(&status); err != nil {
			return nil, fmt.Errorf("error parsing problem status: %v", err)
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (s *Service) SaveUserProblemStatus(ctx context.Context, status *models.UserProblemStatus) error {
	_, err := s.firestoreClient.Collection("users").Doc(status.UserID).Collection("problem_status").Doc(status.ProblemID).Set(ctx, status)
	if err != nil {
		return fmt.Errorf("error saving problem status: %v", err)
	}
	return nil
}

func (s *Service) GetProblemStats(ctx context.Context, problemID string) (*models.ProblemStats, error) {
	doc, err := s.firestoreClient.Collection("problem_stats").Doc(problemID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return &models.ProblemStats{}, nil
		}
		return nil, fmt.Errorf("error fetching problem stats: %v", err)
	}

	var stats models.ProblemStats
	if err := doc.DataTo(&stats); err != nil {
		return nil, fmt.Errorf("error parsing problem stats: %v", err)
	}

	return &stats, nil
}

func (s *Service) GetRecentSubmissions(ctx context.Context, problemID string, limit int) ([]models.Submission, error) {
	iter := s.firestoreClient.Collection("submissions").
		Where("problem_id", "==", problemID).
		OrderBy("submitted_at", firestore.Desc).
		Limit(limit).
		Documents(ctx)
	defer iter.Stop()

	var submissions []models.Submission
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating submissions: %v", err)
		}

		var submission models.Submission
		if err := doc.DataTo(&submission); err != nil {
			return nil, fmt.Errorf("error parsing submission: %v", err)
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (s *Service) GetTopSolutions(ctx context.Context, problemID string) ([]models.Submission, error) {
	iter := s.firestoreClient.Collection("submissions").
		Where("problem_id", "==", problemID).
		Where("status", "==", models.StatusAccepted).
		OrderBy("execution_time", firestore.Asc).
		OrderBy("memory_used", firestore.Asc).
		Limit(5).
		Documents(ctx)
	defer iter.Stop()

	var submissions []models.Submission
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating submissions: %v", err)
		}

		var submission models.Submission
		if err := doc.DataTo(&submission); err != nil {
			return nil, fmt.Errorf("error parsing submission: %v", err)
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (s *Service) GetProblemRankings(ctx context.Context, problemID string, startTime time.Time) ([]models.UserRanking, error) {
	query := s.firestoreClient.Collection("submissions").
		Where("problem_id", "==", problemID).
		Where("status", "==", models.StatusAccepted)

	if !startTime.IsZero() {
		query = query.Where("submitted_at", ">=", startTime)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	userBest := make(map[string]models.Submission)
	usernames := make(map[string]string)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating submissions: %v", err)
		}

		var submission models.Submission
		if err := doc.DataTo(&submission); err != nil {
			return nil, fmt.Errorf("error parsing submission: %v", err)
		}

		if _, ok := usernames[submission.UserID]; !ok {
			userDoc, err := s.firestoreClient.Collection("users").Doc(submission.UserID).Get(ctx)
			if err != nil {
				return nil, fmt.Errorf("error getting user: %v", err)
			}
			var user struct {
				Username string `firestore:"username"`
			}
			if err := userDoc.DataTo(&user); err != nil {
				return nil, fmt.Errorf("error parsing user: %v", err)
			}
			usernames[submission.UserID] = user.Username
		}

		if best, ok := userBest[submission.UserID]; !ok || submission.ExecutionTime < best.ExecutionTime {
			userBest[submission.UserID] = submission
		}
	}

	var rankings []models.UserRanking
	for userID, submission := range userBest {
		timeScore := 1000000.0 / float64(submission.ExecutionTime)
		memoryScore := 1000000.0 / float64(submission.MemoryUsed)
		score := timeScore*0.7 + memoryScore*0.3 // Weight time more heavily than memory

		ranking := models.UserRanking{
			UserID:         userID,
			Username:       usernames[userID],
			Score:          score,
			ProblemsSolved: 1,
			TotalTime:      submission.ExecutionTime,
			TotalMemory:    submission.MemoryUsed,
		}
		rankings = append(rankings, ranking)
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	return rankings, nil
}

func (s *Service) GetGlobalRankings(ctx context.Context, startTime time.Time) ([]models.UserRanking, error) {
	iter := s.firestoreClient.Collection("users").Documents(ctx)
	defer iter.Stop()

	var rankings []models.UserRanking
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating users: %v", err)
		}

		var user struct {
			Username string `firestore:"username"`
		}
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("error parsing user: %v", err)
		}

		summariesIter := s.firestoreClient.Collection("users").Doc(doc.Ref.ID).Collection("submission_summaries").Documents(ctx)

		var totalTime int64
		var totalMemory int64
		var problemsSolved int
		for {
			summaryDoc, err := summariesIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("error iterating submission summaries: %v", err)
			}

			var summary models.SubmissionSummary
			if err := summaryDoc.DataTo(&summary); err != nil {
				return nil, fmt.Errorf("error parsing submission summary: %v", err)
			}

			if summary.Solved && (startTime.IsZero() || summary.FirstSolvedAt.After(startTime)) {
				problemsSolved++
				totalTime += summary.BestTime
				totalMemory += summary.BestMemory
			}
		}

		if problemsSolved > 0 {
			avgTime := float64(totalTime) / float64(problemsSolved)
			avgMemory := float64(totalMemory) / float64(problemsSolved)
			timeScore := 1000000.0 / avgTime
			memoryScore := 1000000.0 / avgMemory
			score := float64(problemsSolved)*50 + timeScore*0.3 + memoryScore*0.2

			ranking := models.UserRanking{
				UserID:         doc.Ref.ID,
				Username:       user.Username,
				Score:          score,
				ProblemsSolved: problemsSolved,
				TotalTime:      totalTime,
				TotalMemory:    totalMemory,
			}
			rankings = append(rankings, ranking)
		}
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	return rankings, nil
}
