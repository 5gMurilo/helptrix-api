CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    helper_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id),
    rate INTEGER NOT NULL CHECK (rate >= 1 AND rate <= 5),
    review TEXT NOT NULL,
    service_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_reviews_business_id ON reviews(business_id);
CREATE INDEX idx_reviews_helper_id ON reviews(helper_id);
CREATE INDEX idx_reviews_deleted_at ON reviews(deleted_at);
CREATE UNIQUE INDEX idx_reviews_unique_active
    ON reviews(business_id, helper_id)
    WHERE deleted_at IS NULL;
