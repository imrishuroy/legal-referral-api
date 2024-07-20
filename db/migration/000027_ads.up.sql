-- ENUM type for ads
CREATE TYPE ad_type AS ENUM ('image', 'video', 'other');

-- ENUM type for payment_cycle
CREATE Type payment_cycle AS ENUM ('weekly', 'monthly', 'yearly');

CREATE TABLE ads (
    ad_id SERIAL PRIMARY KEY,
    ad_type ad_type NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    link TEXT NOT NULL,
    media TEXT[] NOT NULL DEFAULT '{}',
    payment_cycle payment_cycle NOT NULL,
    author_id VARCHAR NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (author_id) REFERENCES users (user_id)
);