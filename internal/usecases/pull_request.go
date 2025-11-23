package usecases

import (
	"context"
	"fmt"
	"slices"

	"pr-manager-service/internal/domain"

	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
)

func (u *Usecases) CreatePullRequest(
	ctx context.Context,
	request domain.CreatePullRequestRequest,
) (domain.PullRequest, error) {
	// NOTE: проверка существования пользователя
	user, err := u.storage.GetUserShort(ctx, request.AuthorUserID)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("storage.GetUserShort: %w", err)
	}

	if !user.IsActive {
		return domain.PullRequest{}, domain.ErrUserInactive
	}

	var pr domain.PullRequest

	if err := u.storage.UnitOfWork(ctx, func(s Storage) error {
		activeColleagues, err := s.GetActiveColleagues(ctx, request.AuthorUserID)
		if err != nil {
			return fmt.Errorf("GetActiveColleagues: %w", err)
		}

		reviewers := SelectRandomElements(activeColleagues, 2)
		reviewersIDs := lo.Map(reviewers, func(u domain.User, _ int) string {
			return u.ID
		})

		createdPr, err := s.CreatePullRequest(ctx, request, reviewersIDs)
		if err != nil {
			return fmt.Errorf("CreatePullRequest: %w", err)
		}
		pr = createdPr

		if err := s.UserAssignmentsIncrementBatch(ctx, reviewersIDs); err != nil {
			return fmt.Errorf("UserAssignmentIncrementMany: %w", err)
		}

		if err := s.PullRequestStatsCreate(ctx, createdPr.ID, len(reviewersIDs)); err != nil {
			return fmt.Errorf("PullRequestStatsCreate: %w", err)
		}

		return nil
	}); err != nil {
		return domain.PullRequest{}, fmt.Errorf("UnitOfWork: %w", err)
	}

	return pr, nil
}

func (u *Usecases) MergePullRequest(ctx context.Context, prID string) (domain.PullRequest, error) {
	if _, err := u.storage.GetPullRequestByID(ctx, prID); err != nil {
		return domain.PullRequest{}, fmt.Errorf("GetPullRequestByID: %w", err)
	}

	if err := u.storage.UpdatePullRequestStatus(ctx, prID, domain.StatusMerged); err != nil {
		return domain.PullRequest{}, fmt.Errorf("storage.UpdatePullRequestStatus: %w", err)
	}

	pullRequest, err := u.storage.GetPullRequestByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("storage.GetPullRequestByID: %w", err)
	}

	return pullRequest, nil
}

func (u *Usecases) ReassignPullRequest(
	ctx context.Context,
	prID, oldUserID string,
) (domain.PullRequest, string, error) {
	// NOTE: проверка существования пользователя
	if _, err := u.storage.GetUserShort(ctx, oldUserID); err != nil {
		return domain.PullRequest{}, "", fmt.Errorf("storage.GetUserShort: %w", err)
	}

	var newReviewerID string

	if err := u.storage.UnitOfWork(ctx, func(s Storage) error {
		pr, err := s.GetPullRequestByID(ctx, prID)
		if err != nil {
			return fmt.Errorf("GetPullRequestByID: %w", err)
		}

		candidates, err := s.GetActiveColleagues(ctx, oldUserID)
		if err != nil {
			return fmt.Errorf("GetActiveColleagues: %w", err)
		}

		candidates = lo.Filter(candidates, func(u domain.User, _ int) bool {
			notAuthor := u.ID != pr.AuthorUserID
			notReviewer := !slices.ContainsFunc(pr.ReviewersUsersIDs, func(reviewerID string) bool {
				return reviewerID == u.ID
			})
			return notAuthor && notReviewer
		})

		if pr.Status == domain.StatusMerged {
			return domain.ErrPRMerged
		}

		if !slices.Contains(pr.ReviewersUsersIDs, oldUserID) {
			return domain.ErrNotAssigned
		}

		if len(candidates) == 0 {
			return domain.ErrNoCandidate
		}

		newReviewerID := SelectRandomElements(candidates, 1)[0].ID
		updatedReviewersIDs := make([]string, 0, 2)

		for _, id := range pr.ReviewersUsersIDs {
			if id == oldUserID {
				continue
			}

			if id == pr.AuthorUserID {
				continue
			}

			updatedReviewersIDs = append(updatedReviewersIDs, id)
		}
		updatedReviewersIDs = append(updatedReviewersIDs, newReviewerID)

		if err := s.UpdatePullRequestReviewersIDs(ctx, prID, updatedReviewersIDs); err != nil {
			return fmt.Errorf("UpdatePullRequestReviewersIDs: %w", err)
		}

		if err := s.UserAssignmentsIncrementBatch(ctx, []string{newReviewerID}); err != nil {
			return fmt.Errorf("UserAssignmentIncrementMany: %w", err)
		}

		if err := s.PullRequestAssignmentsIncrement(ctx, prID); err != nil {
			return fmt.Errorf("PullRequestAssignmentIncrement: %w", err)
		}

		return nil
	}); err != nil {
		return domain.PullRequest{}, "", fmt.Errorf("UnitOfWork: %w", err)
	}

	pr, err := u.storage.GetPullRequestByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, "", fmt.Errorf("GetPullRequestByID: %w", err)
	}

	return pr, newReviewerID, nil
}

func SelectRandomElements[T any](elements []T, count int) []T {
	if len(elements) <= count {
		result := make([]T, len(elements))
		copy(result, elements)
		return result
	}

	shuffled := make([]T, len(elements))
	copy(shuffled, elements)
	mutable.Shuffle(shuffled)

	return shuffled[:count]
}
