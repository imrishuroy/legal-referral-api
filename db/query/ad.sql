-- name: CreateAd :one
INSERT INTO ads (
    ad_type,
    title,
    description,
    link,
    media,
    payment_cycle,
    author_id,
    start_date,
    end_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: ListPlayingAds :many
SELECT *
FROM ads
WHERE end_date > NOW()
ORDER BY start_date DESC;

-- name: ListExpiredAds :many
SELECT *
FROM ads
WHERE end_date < NOW()
ORDER BY start_date DESC;

-- name: ExtendAdPeriod :one
UPDATE ads
SET end_date = $2
WHERE ad_id = $1
RETURNING *;

-- name: GetRandomAd :one
SELECT *
FROM ads
WHERE end_date > NOW()
ORDER BY RANDOM()
LIMIT 1;

-- name: ListRandomAds :many
SELECT *
FROM ads
WHERE end_date > NOW()
ORDER BY RANDOM()
LIMIT $1;