
CREATE TABLE run_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    run_id UUID NOT NULL
        REFERENCES runs(id) ON DELETE CASCADE,

    step_number INTEGER NOT NULL
        CHECK (step_number > 0),

    step_type VARCHAR(50) NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (
            status IN (
                'pending',
                'running',
                'completed',
                'failed',
                'skipped'
            )
        ),

    input JSONB,
    output JSONB,

    error_code VARCHAR(100),
    error_message TEXT,

    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_run_steps_number
        UNIQUE (run_id, step_number)
);

CREATE INDEX idx_run_steps_run_id
    ON run_steps(run_id, step_number);
