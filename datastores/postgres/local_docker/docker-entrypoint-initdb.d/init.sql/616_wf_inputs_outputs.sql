CREATE TABLE public.ai_workflow_stage_references (
    input_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    workflow_run_id BIGINT NOT NULL,
    child_wf_id TEXT NOT NULL,
    running_cycle_number int8 NOT NULL DEFAULT 1,
    chunk_offset int8 NOT NULL DEFAULT 0,
    input_data JSONB NULL,
    logs TEXT NOT NULL DEFAULT '',
    CONSTRAINT fk_workflow_run_id
     FOREIGN KEY (workflow_run_id)
         REFERENCES public.ai_workflow_runs(workflow_run_id)
);

CREATE UNIQUE INDEX idx_child_wf_id ON public.ai_workflow_stage_references(input_id, child_wf_id);
CREATE INDEX idx_workflow_run_id ON public.ai_workflow_stage_references(workflow_run_id);
