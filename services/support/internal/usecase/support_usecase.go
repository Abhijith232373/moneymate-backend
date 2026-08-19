package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/abijith/moneymate-backend/services/support/internal/domain"
)

type SupportUseCase interface {
	CreateFeedback(ctx context.Context, userID uuid.UUID, userType string, rating int, description string) (*domain.Feedback, error)
	ListFeedbacks(ctx context.Context, limit, offset int32) ([]*domain.Feedback, error)

	CreateComplaint(ctx context.Context, userID uuid.UUID, userType string, title, description string) (*domain.Complaint, error)
	ListComplaints(ctx context.Context, limit, offset int32) ([]*domain.Complaint, error)
	ListComplaintsByUser(ctx context.Context, userID uuid.UUID, userType string) ([]*domain.Complaint, error)

	CreateReport(ctx context.Context, reporterID uuid.UUID, reporterType string, reportedVPA, title, description string) (*domain.Report, error)
	ListReports(ctx context.Context, limit, offset int32) ([]*domain.Report, error)
	ListReportsByUser(ctx context.Context, reporterID uuid.UUID, reporterType string) ([]*domain.Report, error)

	CreateAuditLog(ctx context.Context, adminID uuid.UUID, adminName, adminRole, module, action string) (*domain.AuditLog, error)
	ListAuditLogs(ctx context.Context, limit, offset int32) ([]*domain.AuditLog, error)
}

type supportUseCase struct {
	repo domain.SupportRepository
}

func NewSupportUseCase(repo domain.SupportRepository) SupportUseCase {
	return &supportUseCase{
		repo: repo,
	}
}

func (u *supportUseCase) CreateFeedback(ctx context.Context, userID uuid.UUID, userType string, rating int, description string) (*domain.Feedback, error) {
	return u.repo.CreateFeedback(ctx, &domain.Feedback{
		UserID:      userID,
		UserType:    userType,
		Rating:      rating,
		Description: description,
	})
}

func (u *supportUseCase) ListFeedbacks(ctx context.Context, limit, offset int32) ([]*domain.Feedback, error) {
	return u.repo.ListFeedbacks(ctx, limit, offset)
}

func (u *supportUseCase) CreateComplaint(ctx context.Context, userID uuid.UUID, userType string, title, description string) (*domain.Complaint, error) {
	return u.repo.CreateComplaint(ctx, &domain.Complaint{
		UserID:      userID,
		UserType:    userType,
		Title:       title,
		Description: description,
	})
}

func (u *supportUseCase) ListComplaints(ctx context.Context, limit, offset int32) ([]*domain.Complaint, error) {
	return u.repo.ListComplaints(ctx, limit, offset)
}

func (u *supportUseCase) ListComplaintsByUser(ctx context.Context, userID uuid.UUID, userType string) ([]*domain.Complaint, error) {
	return u.repo.ListComplaintsByUser(ctx, userID, userType)
}

func (u *supportUseCase) CreateReport(ctx context.Context, reporterID uuid.UUID, reporterType string, reportedVPA, title, description string) (*domain.Report, error) {
	return u.repo.CreateReport(ctx, &domain.Report{
		ReporterID:   reporterID,
		ReporterType: reporterType,
		ReportedVPA:  reportedVPA,
		Title:        title,
		Description:  description,
	})
}

func (u *supportUseCase) ListReports(ctx context.Context, limit, offset int32) ([]*domain.Report, error) {
	return u.repo.ListReports(ctx, limit, offset)
}

func (u *supportUseCase) ListReportsByUser(ctx context.Context, reporterID uuid.UUID, reporterType string) ([]*domain.Report, error) {
	return u.repo.ListReportsByUser(ctx, reporterID, reporterType)
}

func (u *supportUseCase) CreateAuditLog(ctx context.Context, adminID uuid.UUID, adminName, adminRole, module, action string) (*domain.AuditLog, error) {
	return u.repo.CreateAuditLog(ctx, &domain.AuditLog{
		AdminID:   adminID,
		AdminName: adminName,
		AdminRole: adminRole,
		Module:    module,
		Action:    action,
	})
}

func (u *supportUseCase) ListAuditLogs(ctx context.Context, limit, offset int32) ([]*domain.AuditLog, error) {
	return u.repo.ListAuditLogs(ctx, limit, offset)
}
