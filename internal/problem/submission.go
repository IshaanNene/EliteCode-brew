package problem

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/IshaanNene/EliteCode-brew/internal/docker"
	"github.com/IshaanNene/EliteCode-brew/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubmissionService struct {
	firestoreClient *firestore.Client
	dockerClient    *docker.Client
	problemService  *Service
}

func NewSubmissionService(firestoreClient *firestore.Client, dockerClient *docker.Client, problemService *Service) *SubmissionService {
	return &SubmissionService{
		firestoreClient: firestoreClient,
		dockerClient:    dockerClient,
		problemService:  problemService,
	}
}

func (s *SubmissionService) Submit(ctx context.Context, userID string, problemID string, code string, language string) (*models.Submission, error) {
	problem, err := s.problemService.GetProblem(ctx, problemID)
	if err != nil {
		return nil, fmt.Errorf("error getting problem: %v", err)
	}

	submission := &models.Submission{
		ID:          uuid.New().String(),
		UserID:      userID,
		ProblemID:   problemID,
		Language:    language,
		Code:        code,
		Status:      models.StatusPending,
		SubmittedAt: time.Now(),
	}

	_, err = s.firestoreClient.Collection("submissions").Doc(submission.ID).Set(ctx, submission)
	if err != nil {
		return nil, fmt.Errorf("error saving submission: %v", err)
	}

	testCasesSnap, err := s.firestoreClient.Collection("problems").Doc(problemID).Collection("test_cases").Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting test cases: %v", err)
	}

	var testCases []models.TestCase
	for _, doc := range testCasesSnap {
		var tc models.TestCase
		if err := doc.DataTo(&tc); err != nil {
			return nil, fmt.Errorf("error parsing test case: %v", err)
		}
		tc.ID = doc.Ref.ID
		testCases = append(testCases, tc)
	}

	submission.Status = models.StatusRunning
	_, err = s.firestoreClient.Collection("submissions").Doc(submission.ID).Set(ctx, submission)
	if err != nil {
		return nil, fmt.Errorf("error updating submission status: %v", err)
	}

	var totalTime int64
	var maxMemory int64
	var failedTests int
	submission.TestCases = make([]models.TestCaseResult, len(testCases))

	for i, tc := range testCases {
		result := models.TestCaseResult{
			TestCaseID: tc.ID,
			Status:     models.StatusRunning,
		}

		output, execTime, memUsed, err := s.runTestCase(ctx, code, language, tc.Input, problem.TimeLimit, problem.MemoryLimit)
		if err != nil {
			result.Status = models.StatusError
			result.ErrorMessage = err.Error()
			failedTests++
		} else {
			result.ExecutionTime = execTime.Milliseconds()
			result.MemoryUsed = int64(memUsed * 1024) // Convert MB to KB
			totalTime += result.ExecutionTime
			if result.MemoryUsed > maxMemory {
				maxMemory = result.MemoryUsed
			}

			if result.ExecutionTime > int64(problem.TimeLimit) {
				result.Status = models.StatusTimeLimitExceeded
				failedTests++
			} else if result.MemoryUsed > int64(problem.MemoryLimit)*1024 { // Convert MB to KB
				result.Status = models.StatusMemoryLimitExceeded
				failedTests++
			} else {
				result.ExpectedOutput = tc.Expected
				result.ActualOutput = output
				if output == tc.Expected {
					result.Status = models.StatusAccepted
				} else {
					result.Status = models.StatusRejected
					failedTests++
				}
			}
		}

		submission.TestCases[i] = result
	}

	submission.CompletedAt = time.Now()
	submission.ExecutionTime = totalTime / int64(len(testCases)) // Average time
	submission.MemoryUsed = maxMemory

	if failedTests == 0 {
		submission.Status = models.StatusAccepted
	} else {
		submission.Status = models.StatusRejected
	}

	_, err = s.firestoreClient.Collection("submissions").Doc(submission.ID).Set(ctx, submission)
	if err != nil {
		return nil, fmt.Errorf("error saving final submission: %v", err)
	}

	if err := s.updateSubmissionSummary(ctx, userID, problemID, submission); err != nil {
		return nil, fmt.Errorf("error updating submission summary: %v", err)
	}

	return submission, nil
}

func (s *SubmissionService) runTestCase(ctx context.Context, code string, language string, input string, timeLimit int, memoryLimit int) (string, time.Duration, float64, error) {
	files := map[string][]byte{
		"code":      []byte(code),
		"input.txt": []byte(input),
	}

	buildCtx, err := docker.CreateBuildContext(files)
	if err != nil {
		return "", 0, 0, fmt.Errorf("error creating build context: %v", err)
	}

	imageName := fmt.Sprintf("elitecode/submission:%s", uuid.New().String())
	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	}

	if err := s.dockerClient.BuildImage(ctx, buildCtx, buildOptions); err != nil {
		return "", 0, 0, fmt.Errorf("error building Docker image: %v", err)
	}

	containerConfig := &container.Config{
		Image: imageName,
		Cmd:   []string{"./run.sh"},
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory:    int64(memoryLimit) * 1024 * 1024, // Convert MB to bytes
			CPUPeriod: 100000,
			CPUQuota:  50000, // Limit to 0.5 CPU
		},
	}

	containerName := fmt.Sprintf("elitecode_submission_%s", uuid.New().String())
	output, execTime, memUsed, err := s.dockerClient.RunContainer(ctx, containerConfig, hostConfig, containerName)
	if err != nil {
		return "", 0, 0, fmt.Errorf("error running container: %v", err)
	}

	return string(output), execTime, memUsed, nil
}

func (s *SubmissionService) updateSubmissionSummary(ctx context.Context, userID string, problemID string, submission *models.Submission) error {
	summaryRef := s.firestoreClient.Collection("users").Doc(userID).Collection("submission_summaries").Doc(problemID)

	return s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var summary models.SubmissionSummary
		doc, err := tx.Get(summaryRef)
		if err != nil && status.Code(err) != codes.NotFound {
			return fmt.Errorf("error getting submission summary: %v", err)
		}

		if doc != nil && doc.Exists() {
			if err := doc.DataTo(&summary); err != nil {
				return fmt.Errorf("error parsing submission summary: %v", err)
			}
		}

		summary.ProblemID = problemID
		summary.LastSubmissionID = submission.ID
		summary.AttemptCount++
		summary.LastAttemptedAt = submission.SubmittedAt

		if submission.Status == models.StatusAccepted {
			summary.Solved = true
			if summary.FirstSolvedAt.IsZero() {
				summary.FirstSolvedAt = submission.CompletedAt
			}
			if submission.ExecutionTime < summary.BestTime || summary.BestTime == 0 {
				summary.BestTime = submission.ExecutionTime
				summary.BestSubmissionID = submission.ID
			}
			if submission.MemoryUsed < summary.BestMemory || summary.BestMemory == 0 {
				summary.BestMemory = submission.MemoryUsed
			}
		}

		return tx.Set(summaryRef, summary)
	})
}
