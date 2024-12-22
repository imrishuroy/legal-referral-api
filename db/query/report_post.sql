-- name: AddReportReason :one
INSERT INTO report_reasons (
    reason
) VALUES (
    $1
) RETURNING reason_id;

-- name: ReportPost :exec
INSERT INTO reported_posts (
    post_id,
    reported_by,
    reason_id
) VALUES (
    $1, $2, $3
);

-- name: IsPostReported :one
SELECT
    CASE WHEN report_id IS NOT NULL THEN true ELSE false END AS is_reported
FROM reported_posts
WHERE post_id = $1 AND reported_by = $2;
