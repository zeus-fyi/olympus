CREATE TABLE public.ai_workflow_analysis_results(
    workflow_result_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    orchestration_id int8 NOT NULL REFERENCES orchestrations(orchestration_id),
    response_id int8 NOT NULL REFERENCES completion_responses(response_id),
    source_task_id int8 NOT NULL REFERENCES ai_task_library(task_id),
    running_cycle_number int8 NOT NULL DEFAULT 1,
    iteration_count int8 NOT NULL DEFAULT 1,
    chunk_offset int8 NOT NULL DEFAULT 0,
    task_offset int8 NOT NULL DEFAULT 0,
    search_window_unix_start int8 NOT NULL CHECK (search_window_unix_start < search_window_unix_end),
    search_window_unix_end int8 NOT NULL CHECK (search_window_unix_start < search_window_unix_end),
    skip_analysis bool NOT NULL DEFAULT false,
    metadata jsonb
);

CREATE INDEX wf_analysis_to_id_idx ON public.ai_workflow_analysis_results("task_offset");
CREATE INDEX wf_analysis_skip_idx ON public.ai_workflow_analysis_results("skip_analysis");
CREATE INDEX wf_analysis_orchestrations_id_idx ON public.ai_workflow_analysis_results("orchestration_id");
CREATE INDEX wf_analysis_resp_id_idx ON public.ai_workflow_analysis_results("response_id");
CREATE INDEX wf_analysis_cycle_idx ON public.ai_workflow_analysis_results("running_cycle_number");
CREATE INDEX wf_analysis_source_task_idx ON public.ai_workflow_analysis_results("source_task_id");
CREATE INDEX wf_analysis_source_search_start_idx ON public.ai_workflow_analysis_results("search_window_unix_start");
CREATE INDEX wf_analysis_source_search_end_idx ON public.ai_workflow_analysis_results("search_window_unix_end");
CREATE INDEX wf_analysis_metadata_idx ON public.ai_workflow_analysis_results USING GIN (metadata);

ALTER TABLE public.ai_workflow_analysis_results
    ADD CONSTRAINT unique_combination_wf_analysis_iteration UNIQUE (orchestration_id, response_id, source_task_id, running_cycle_number, iteration_count, chunk_offset);
