-- name: CreateFeedback :one
INSERT INTO feedbacks (
  user_id, user_type, rating, description
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: ListFeedbacks :many
SELECT * FROM feedbacks
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateComplaint :one
INSERT INTO complaints (
  user_id, user_type, title, description
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: ListComplaints :many
SELECT * FROM complaints
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListComplaintsByUser :many
SELECT * FROM complaints
WHERE user_id = $1 AND user_type = $2
ORDER BY created_at DESC;


-- name: CreateReport :one
INSERT INTO reports (
  reporter_id, reporter_type, reported_vpa, title, description
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListReports :many
SELECT * FROM reports
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListReportsByUser :many
SELECT * FROM reports
WHERE reporter_id = $1 AND reporter_type = $2
ORDER BY created_at DESC;

