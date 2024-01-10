CREATE TABLE ai_assistants (
    assistant_id TEXT PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    object text NOT NULL DEFAULT 'assistant',
    created_at BIGINT,
    assistant_name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL,
    instructions TEXT,
    tools JSONB,
    file_ids JSONB,
    metadata JSONB
);
ALTER TABLE public.ai_assistants
ADD CONSTRAINT uniq_org_assistant_names UNIQUE (org_id, assistant_name);
CREATE INDEX ai_assistants_org_ind ON public.ai_assistants("org_id");
