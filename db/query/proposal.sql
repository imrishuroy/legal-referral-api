-- name: CreateProposal :one
INSERT INTO proposals (
    referral_id,
    user_id,
    title,
    proposal
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateProposal :one
UPDATE proposals
SET
    title = $3,
    proposal = $4
WHERE proposal_id = $1 AND user_id = $2
RETURNING *;

-- name: GetProposal :one
SELECT *
FROM proposals
WHERE referral_id = $1 AND user_id = $2;
