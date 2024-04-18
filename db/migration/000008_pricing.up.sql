CREATE TABLE pricing (
    price_id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    service_type VARCHAR NOT NULL,
    per_hour_price DECIMAL,
    per_hearing_price DECIMAL,
    contingency_price VARCHAR ,
    hybrid_price VARCHAR
);