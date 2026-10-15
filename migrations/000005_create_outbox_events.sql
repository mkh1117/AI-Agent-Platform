
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,

    payload JSONB NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'publishing', 'published')),

    attempts INTEGER NOT NULL DEFAULT 0
        CHECK (attempts >= 0),

    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_until TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    last_error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (
        (status = 'published' AND published_at IS NOT NULL)
        OR status <> 'published'
    )
);

CREATE INDEX idx_outbox_pending
    ON outbox_events (available_at, created_at)
    WHERE status = 'pending';

CREATE INDEX idx_outbox_expired_leases
    ON outbox_events (locked_until)
    WHERE status = 'publishing';
