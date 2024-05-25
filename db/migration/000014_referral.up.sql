-- ENUM type for status
CREATE TYPE  referrals_status AS ENUM ('active', 'awarded', 'completed', 'cancelled', 'rejected');

CREATE TABLE referrals (
    referral_id SERIAL PRIMARY KEY,
    referred_user_id VARCHAR NOT NULL,
    referrer_user_id VARCHAR NOT NULL,
    title TEXT NOT NULL,
    preferred_practice_area TEXT NOT NULL,
    preferred_practice_location TEXT NOT NULL,
    case_description TEXT NOT NULL,
    status referrals_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (referred_user_id) REFERENCES users(user_id),
    FOREIGN KEY (referrer_user_id) REFERENCES users(user_id)
);