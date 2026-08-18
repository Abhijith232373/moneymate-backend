package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type staffRepository struct {
	db *pgxpool.Pool
}

func NewStaffRepo(db *pgxpool.Pool) domain.StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	query := `
		INSERT INTO auth.staff (id, email, full_name, password_hash, status, token_version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		staff.ID, staff.Email, staff.FullName, staff.PasswordHash,
		staff.Status, staff.TokenVersion, staff.CreatedAt, staff.UpdatedAt,
	)
	if err != nil {
		mappedErr := apperrors.MapDBErrors(err)
		if mappedErr != err {
			return mappedErr
		}
		return err
	}
	return nil
}

func (r *staffRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error) {
	query := `
		SELECT id, email, full_name, password_hash, status, token_version, created_at, updated_at
		FROM auth.staff
		WHERE id = $1 AND status != 'deleted'
	`
	return r.scanRow(r.db.QueryRow(ctx, query, id))
}

func (r *staffRepository) GetByEmail(ctx context.Context, email string) (*domain.Staff, error) {
	query := `
		SELECT id, email, full_name, password_hash, status, token_version, created_at, updated_at
		FROM auth.staff
		WHERE email = $1 AND status != 'deleted'
	`
	return r.scanRow(r.db.QueryRow(ctx, query, email))
}

func (r *staffRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.staff WHERE email = $1)`
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *staffRepository) UpdateStatus(ctx context.Context, staffID uuid.UUID, status domain.UserStatus) error {
	query := `UPDATE auth.staff SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.Exec(ctx, query, status, staffID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return apperrors.ErrUserNotFound
	}
	return nil
}

func (r *staffRepository) SoftDelete(ctx context.Context, staffID uuid.UUID) error {
	return r.UpdateStatus(ctx, staffID, domain.UserStatusDeleted)
}

func (r *staffRepository) GetTokenVersion(ctx context.Context, staffID uuid.UUID) (int64, error) {
	var version int64
	query := `SELECT token_version FROM auth.staff WHERE id = $1`
	err := r.db.QueryRow(ctx, query, staffID).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apperrors.ErrUserNotFound
		}
		return 0, err
	}
	return version, nil
}

func (r *staffRepository) IncrementTokenVersion(ctx context.Context, staffID uuid.UUID) (int64, error) {
	var newVersion int64
	query := `
		UPDATE auth.staff
		SET token_version = token_version + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING token_version
	`
	err := r.db.QueryRow(ctx, query, staffID).Scan(&newVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apperrors.ErrUserNotFound
		}
		return 0, err
	}
	return newVersion, nil
}

func (r *staffRepository) Update(ctx context.Context, staffID uuid.UUID, fullName, email, passwordHash *string) (*domain.Staff, error) {
	var parts []string
	var args []any
	argID := 1

	if fullName != nil {
		parts = append(parts, fmt.Sprintf("full_name = $%d", argID))
		args = append(args, *fullName)
		argID++
	}
	if email != nil {
		parts = append(parts, fmt.Sprintf("email = $%d", argID))
		args = append(args, *email)
		argID++
	}
	if passwordHash != nil {
		parts = append(parts, fmt.Sprintf("password_hash = $%d", argID))
		args = append(args, *passwordHash)
		argID++
	}

	if len(parts) == 0 {
		return r.GetByID(ctx, staffID)
	}

	parts = append(parts, "updated_at = NOW()")
	args = append(args, staffID)

	query := fmt.Sprintf(`
		UPDATE auth.staff
		SET %s
		WHERE id = $%d AND status != 'deleted'
		RETURNING id, email, full_name, password_hash, status, token_version, created_at, updated_at
	`, strings.Join(parts, ", "), argID)

	return r.scanRow(r.db.QueryRow(ctx, query, args...))
}

func (r *staffRepository) ListStaff(ctx context.Context, filter domain.ListUsersFilter, page domain.Pagination) (*domain.ListStaffResult, error) {
	baseQuery := `FROM auth.staff WHERE status != 'deleted'`
	var args []any
	argID := 1
	var conditions []string

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argID))
		args = append(args, filter.Status)
		argID++
	}

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(email ILIKE $%d OR full_name ILIKE $%d)", argID, argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	orderDir := "ASC"
	if filter.SortDesc {
		orderDir = "DESC"
	}
	sortBy := "created_at"
	if filter.SortBy == "email" || filter.SortBy == "full_name" || filter.SortBy == "status" {
		sortBy = filter.SortBy
	}
	orderClause := fmt.Sprintf("ORDER BY %s %s", sortBy, orderDir)

	limit := page.PageSize
	offset := (page.Page - 1) * limit
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf("SELECT id, email, full_name, password_hash, status, token_version, created_at, updated_at %s %s LIMIT $%d OFFSET $%d", baseQuery, orderClause, argID, argID+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staffList []domain.Staff
	for rows.Next() {
		var s domain.Staff
		if err := rows.Scan(&s.ID, &s.Email, &s.FullName, &s.PasswordHash, &s.Status, &s.TokenVersion, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		staffList = append(staffList, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &domain.ListStaffResult{Staff: staffList, TotalCount: total}, nil
}

func (r *staffRepository) AssignRole(ctx context.Context, staffID, roleID uuid.UUID) error {
	query := `
		INSERT INTO auth.staff_roles (staff_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, staffID, roleID)
	return err
}

func (r *staffRepository) RemoveRole(ctx context.Context, staffID, roleID uuid.UUID) error {
	query := `DELETE FROM auth.staff_roles WHERE staff_id = $1 AND role_id = $2`
	_, err := r.db.Exec(ctx, query, staffID, roleID)
	return err
}

func (r *staffRepository) GetRoles(ctx context.Context, staffID uuid.UUID) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM auth.roles r
		JOIN auth.staff_roles sr ON r.id = sr.role_id
		WHERE sr.staff_id = $1
	`
	rows, err := r.db.Query(ctx, query, staffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var r domain.Role
		var desc sql.NullString
		if err := rows.Scan(&r.ID, &r.Name, &desc, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			r.Description = &desc.String
		}
		roles = append(roles, r)
	}
	return roles, rows.Err()
}

func (r *staffRepository) scanRow(row pgx.Row) (*domain.Staff, error) {
	var s domain.Staff
	err := row.Scan(&s.ID, &s.Email, &s.FullName, &s.PasswordHash, &s.Status, &s.TokenVersion, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}
	return &s, nil
}
