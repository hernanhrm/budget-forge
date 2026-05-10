CREATE TABLE budget.months (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    year       INT NOT NULL CHECK (year >= 2000 AND year <= 2100),
    month      INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    status     TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'draft', 'active', 'closed')),
    closed_at  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_months_user_year_month UNIQUE NULLS NOT DISTINCT (user_id, year, month)
);

CREATE INDEX idx_months_user_id ON budget.months(user_id);
CREATE INDEX idx_months_status ON budget.months(status);
