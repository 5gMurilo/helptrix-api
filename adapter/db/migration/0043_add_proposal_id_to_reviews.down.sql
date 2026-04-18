DROP INDEX IF EXISTS idx_reviews_unique_active_proposal;
DROP INDEX IF EXISTS idx_reviews_proposal_id;

ALTER TABLE reviews
    DROP CONSTRAINT IF EXISTS fk_reviews_proposal_id;

ALTER TABLE reviews
    DROP COLUMN IF EXISTS proposal_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_unique_active
    ON reviews(business_id, helper_id)
    WHERE deleted_at IS NULL;
