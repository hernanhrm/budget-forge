CREATE TABLE budget.transactions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    date              DATE NOT NULL,
    amount            BIGINT NOT NULL,
    category_id       UUID REFERENCES budget.categories(id) ON DELETE SET NULL,
    account_id        UUID REFERENCES budget.accounts(id) ON DELETE SET NULL,
    description       TEXT,
    transfer_pair_id  UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_user_id ON budget.transactions(user_id);
CREATE INDEX idx_transactions_category_id ON budget.transactions(category_id);
CREATE INDEX idx_transactions_account_id ON budget.transactions(account_id);
CREATE INDEX idx_transactions_date ON budget.transactions(date);
CREATE INDEX idx_transactions_transfer_pair_id ON budget.transactions(transfer_pair_id);
CREATE INDEX idx_transactions_user_date ON budget.transactions(user_id, date);
