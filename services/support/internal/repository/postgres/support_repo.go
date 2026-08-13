package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/abijith/moneymate-backend/services/support/internal/domain"
)

type SupportRepo struct {
	db      *sql.DB
	querier Querier
}

func NewSupportRepo(db *sql.DB) *SupportRepo {
	return &SupportRepo{
		db:      db,
		querier: New(db),
	}
}

func (r *SupportRepo) CreateFeedback(ctx context.Context, f *domain.Feedback) (*domain.Feedback, error) {
	row, err := r.querier.CreateFeedback(ctx, CreateFeedbackParams{
		UserID:      f.UserID,
		UserType:    f.UserType,
		Rating:      int32(f.Rating),
		Description: f.Description,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Feedback{
		ID:          row.ID,
		UserID:      row.UserID,
		UserType:    row.UserType,
		Rating:      int(row.Rating),
		Description: row.Description,
		CreatedAt:   row.CreatedAt.Time,
	}, nil
}

func (r *SupportRepo) ListFeedbacks(ctx context.Context, limit, offset int32) ([]*domain.Feedback, error) {
	rows, err := r.querier.ListFeedbacks(ctx, ListFeedbacksParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var res []*domain.Feedback
	for _, row := range rows {
		res = append(res, &domain.Feedback{
			ID:          row.ID,
			UserID:      row.UserID,
			UserType:    row.UserType,
			Rating:      int(row.Rating),
			Description: row.Description,
			CreatedAt:   row.CreatedAt.Time,
		})
	}
	return res, nil
}

func (r *SupportRepo) CreateComplaint(ctx context.Context, c *domain.Complaint) (*domain.Complaint, error) {
	row, err := r.querier.CreateComplaint(ctx, CreateComplaintParams{
		UserID:      c.UserID,
		UserType:    c.UserType,
		Title:       c.Title,
		Description: c.Description,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Complaint{
		ID:          row.ID,
		UserID:      row.UserID,
		UserType:    row.UserType,
		Title:       row.Title,
		Description: row.Description,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}, nil
}

func (r *SupportRepo) ListComplaints(ctx context.Context, limit, offset int32) ([]*domain.Complaint, error) {
	rows, err := r.querier.ListComplaints(ctx, ListComplaintsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var res []*domain.Complaint
	for _, row := range rows {
		res = append(res, &domain.Complaint{
			ID:          row.ID,
			UserID:      row.UserID,
			UserType:    row.UserType,
			Title:       row.Title,
			Description: row.Description,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		})
	}
	return res, nil
}

func (r *SupportRepo) ListComplaintsByUser(ctx context.Context, userID uuid.UUID, userType string) ([]*domain.Complaint, error) {
	rows, err := r.querier.ListComplaintsByUser(ctx, ListComplaintsByUserParams{
		UserID:   userID,
		UserType: userType,
	})
	if err != nil {
		return nil, err
	}

	var res []*domain.Complaint
	for _, row := range rows {
		res = append(res, &domain.Complaint{
			ID:          row.ID,
			UserID:      row.UserID,
			UserType:    row.UserType,
			Title:       row.Title,
			Description: row.Description,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		})
	}
	return res, nil
}



func (r *SupportRepo) CreateReport(ctx context.Context, rp *domain.Report) (*domain.Report, error) {
	row, err := r.querier.CreateReport(ctx, CreateReportParams{
		ReporterID:   rp.ReporterID,
		ReporterType: rp.ReporterType,
		ReportedVpa:  rp.ReportedVPA,
		Title:        rp.Title,
		Description:  rp.Description,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Report{
		ID:           row.ID,
		ReporterID:   row.ReporterID,
		ReporterType: row.ReporterType,
		ReportedVPA:  row.ReportedVpa,
		Title:        row.Title,
		Description:  row.Description,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}, nil
}

func (r *SupportRepo) ListReports(ctx context.Context, limit, offset int32) ([]*domain.Report, error) {
	rows, err := r.querier.ListReports(ctx, ListReportsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var res []*domain.Report
	for _, row := range rows {
		res = append(res, &domain.Report{
			ID:           row.ID,
			ReporterID:   row.ReporterID,
			ReporterType: row.ReporterType,
			ReportedVPA:  row.ReportedVpa,
			Title:        row.Title,
			Description:  row.Description,
			CreatedAt:    row.CreatedAt.Time,
			UpdatedAt:    row.UpdatedAt.Time,
		})
	}
	return res, nil
}

func (r *SupportRepo) ListReportsByUser(ctx context.Context, reporterID uuid.UUID, reporterType string) ([]*domain.Report, error) {
	rows, err := r.querier.ListReportsByUser(ctx, ListReportsByUserParams{
		ReporterID:   reporterID,
		ReporterType: reporterType,
	})
	if err != nil {
		return nil, err
	}

	var res []*domain.Report
	for _, row := range rows {
		res = append(res, &domain.Report{
			ID:           row.ID,
			ReporterID:   row.ReporterID,
			ReporterType: row.ReporterType,
			ReportedVPA:  row.ReportedVpa,
			Title:        row.Title,
			Description:  row.Description,
			CreatedAt:    row.CreatedAt.Time,
			UpdatedAt:    row.UpdatedAt.Time,
		})
	}
	return res, nil
}


