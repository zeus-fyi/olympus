CREATE TABLE public.ai_workflow_runs(
    workflow_run_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    orchestration_id int8 NOT NULL REFERENCES orchestrations(orchestration_id),
    is_archived bool NOT NULL DEFAULT false
);
CREATE INDEX runs_ai_workflow_runs_orch_idx ON public.ai_workflow_runs("orchestration_id");
CREATE INDEX runs_ai_workflow_runs_run_idx ON public.ai_workflow_runs("workflow_run_id");
CREATE INDEX runs_ai_workflow_runs_archive_idx ON public.ai_workflow_runs("is_archived");

ALTER TABLE public.ai_workflow_runs
    ADD COLUMN total_api_requests int8 NOT NULL DEFAULT 0;

ALTER TABLE public.ai_workflow_runs
    ADD COLUMN total_csv_cells int8 NOT NULL DEFAULT 0;