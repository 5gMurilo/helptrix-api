ALTER TABLE reviews
    ADD COLUMN IF NOT EXISTS proposal_id UUID;

UPDATE reviews r
SET proposal_id = matched_proposals.proposal_id
FROM (
    SELECT
        r_inner.id AS review_id,
        MIN(p.id) AS proposal_id
    FROM reviews r_inner
    JOIN proposals p
        ON p.user_id = r_inner.business_id
        AND p.helper_id = r_inner.helper_id
        AND p.category_id = r_inner.category_id
        AND p.status = 'finished'
    WHERE r_inner.proposal_id IS NULL
    GROUP BY r_inner.id
    HAVING COUNT(*) = 1
) matched_proposals
WHERE r.id = matched_proposals.review_id;

ALTER TABLE reviews
    ADD CONSTRAINT fk_reviews_proposal_id
    FOREIGN KEY (proposal_id) REFERENCES proposals(id) ON DELETE CASCADE;

DROP INDEX IF EXISTS idx_reviews_unique_active;

CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_unique_active_proposal
    ON reviews(proposal_id)
    WHERE deleted_at IS NULL AND proposal_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_reviews_proposal_id ON reviews(proposal_id);
