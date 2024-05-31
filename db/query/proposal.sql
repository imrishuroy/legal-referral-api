-- name: CreateProposal :one
INSERT INTO proposals (
    project_id,
    user_id,
    title,
    proposal,
    status
) VALUES (
    $1, $2, $3, $4, 'active'
) RETURNING *;
--
-- name: UpdateProposal :one
UPDATE proposals
SET
    title = $3,
    proposal = $4
WHERE proposal_id = $1 AND user_id = $2
RETURNING *;
--
-- name: GetProposal :one
SELECT *
FROM proposals
WHERE project_id = $1 AND user_id = $2;
