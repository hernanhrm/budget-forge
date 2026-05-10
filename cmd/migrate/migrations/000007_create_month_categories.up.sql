CREATE TABLE budget.month_categories (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    month_id       UUID NOT NULL REFERENCES budget.months(id) ON DELETE CASCADE,
    category_id    UUID NOT NULL REFERENCES budget.categories(id) ON DELETE CASCADE,
    planned_amount BIGINT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_month_categories_month_category UNIQUE NULLS NOT DISTINCT (month_id, category_id)
);

CREATE INDEX idx_month_categories_month_id ON budget.month_categories(month_id);
CREATE INDEX idx_month_categories_category_id ON budget.month_categories(category_id);
