CREATE TABLE public.ai_workflow_io_results(
    workflow_result_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    orchestration_id int8 NOT NULL REFERENCES orchestrations(orchestration_id),
    retrieval_id BIGINT NOT NULL REFERENCES ai_retrieval_library(retrieval_id),
    running_cycle_number int8 NOT NULL DEFAULT 1,
    iteration_count int8 NOT NULL DEFAULT 0,
    chunk_offset int8 NOT NULL DEFAULT 0,
    attempts int8 NOT NULL DEFAULT 0,
    search_window_unix_start int8 NOT NULL CHECK (search_window_unix_start < search_window_unix_end),
    search_window_unix_end int8 NOT NULL CHECK (search_window_unix_start < search_window_unix_end),
    skip_retrieval bool NOT NULL DEFAULT false,
    status text NOT NULL DEFAULT 'pending',
    metadata jsonb
);
CREATE INDEX wf_io_skip_idx ON public.ai_workflow_io_results("skip_retrieval");
CREATE INDEX wf_io_orchestration_id_idx ON public.ai_workflow_io_results("orchestration_id");
CREATE INDEX wf_io_retrieval_id_idx ON public.ai_workflow_io_results("retrieval_id");
CREATE INDEX wf_io_cycle_idx ON public.ai_workflow_io_results("running_cycle_number");
CREATE INDEX wf_io_search_start_idx ON public.ai_workflow_io_results("search_window_unix_start");
CREATE INDEX wf_io_search_end_idx ON public.ai_workflow_io_results("search_window_unix_end");
CREATE INDEX wf_io_status_idx ON public.ai_workflow_io_results("status");
CREATE INDEX wf_io_metadata_idx ON public.ai_workflow_io_results USING GIN (metadata);

ALTER TABLE public.ai_workflow_io_results
    ADD CONSTRAINT unique_combination_wf_io UNIQUE (orchestration_id, retrieval_id, running_cycle_number, iteration_count, chunk_offset);

