CREATE TABLE public.ai_workflow_runs(
    workflow_run_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    orchestration_id int8 NOT NULL REFERENCES orchestrations(orchestration_id)
);
CREATE INDEX runs_ai_workflow_runs_orch_idx ON public.ai_workflow_runs("orchestration_id");