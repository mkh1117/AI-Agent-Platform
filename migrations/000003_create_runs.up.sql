
CREATE TABLE runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID NOT NULL,
    agent_version_id UUID NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'queued'
        CHECK (
            status IN (
                'queued',
                'running',
                'completed',
                'failed',
                'cancelled'
            )
        ),

    input JSONB NOT NULL,
    output JSONB,

    error_code VARCHAR(100),
    error_message TEXT,

    attempts INTEGER NOT NULL DEFAULT 0
        CHECK (attempts >= 0),

    queued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_runs_agent
        FOREIGN KEY (agent_id)
        REFERENCES agents(id),

    CONSTRAINT fk_runs_agent_version
        FOREIGN KEY (agent_version_id, agent_id)
        REFERENCES agent_versions(id, agent_id)
);

CREATE INDEX idx_runs_agent_created
    ON runs(agent_id, created_at DESC);

CREATE INDEX idx_runs_status_queued
    ON runs(status, queued_at)
    WHERE status = 'queued';
