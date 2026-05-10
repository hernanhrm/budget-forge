CREATE TABLE budget.categories (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    direction  TEXT NOT NULL CHECK (direction IN ('debit', 'credit')),
    group_id   UUID REFERENCES budget.category_groups(id) ON DELETE SET NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_categories_user_id ON budget.categories(user_id);
CREATE INDEX idx_categories_group_id ON budget.categories(group_id);
