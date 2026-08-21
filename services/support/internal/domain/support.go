package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	UserType    string    `json:"user_type"`
	Rating      int       `json:"rating"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Complaint struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	UserType    string    `json:"user_type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Report struct {
	ID           uuid.UUID `json:"id"`
	ReporterID   uuid.UUID `json:"reporter_id"`
	ReporterType string    `json:"reporter_type"`
	ReportedVPA  string    `json:"reported_vpa"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID        uuid.UUID              `json:"id"`
	AdminID   uuid.UUID              `json:"admin_id"`
	AdminName string                 `json:"admin_name"`
	AdminRole string                 `json:"admin_role"`
	Module    string                 `json:"module"`
	Action    string                 `json:"action"`
	Changes   map[string]interface{} `json:"changes"`
	CreatedAt time.Time              `json:"created_at"`
}

type SupportRepository interface {
	CreateFeedback(ctx context.Context, f *Feedback) (*Feedback, error)
	ListFeedbacks(ctx context.Context, limit, offset int32) ([]*Feedback, error)

	CreateComplaint(ctx context.Context, c *Complaint) (*Complaint, error)
	ListComplaints(ctx context.Context, limit, offset int32) ([]*Complaint, error)
	ListComplaintsByUser(ctx context.Context, userID uuid.UUID, userType string) ([]*Complaint, error)

	CreateReport(ctx context.Context, r *Report) (*Report, error)
	ListReports(ctx context.Context, limit, offset int32) ([]*Report, error)
	ListReportsByUser(ctx context.Context, reporterID uuid.UUID, reporterType string) ([]*Report, error)

	CreateAuditLog(ctx context.Context, log *AuditLog) (*AuditLog, error)
	ListAuditLogs(ctx context.Context, limit, offset int32) ([]*AuditLog, error)
}
