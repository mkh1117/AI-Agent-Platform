
CREATE TABLE agent_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    agent_id UUID NOT NULL
        REFERENCES agents(id) ON DELETE CASCADE,

    version_number INTEGER NOT NULL
        CHECK (version_number > 0),

    instructions TEXT NOT NULL,

    model_config JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_agent_versions_number
        UNIQUE (agent_id, version_number),

    CONSTRAINT uq_agent_versions_id_agent
        UNIQUE (id, agent_id)
);

CREATE INDEX idx_agent_versions_agent_id
    ON agent_versions(agent_id);
