CREATE TABLE public.ai_workflow_eval_result_responses(
    eval_result_id BIGINT NOT NULL PRIMARY KEY,
    workflow_result_id BIGINT NOT NULL REFERENCES ai_workflow_analysis_results(workflow_result_id),
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
    response_id int8 NOT NULL REFERENCES completion_responses(response_id)
);

CREATE INDEX idx_workflow_result_id_evr ON public.ai_workflow_eval_result_responses(workflow_result_id);
CREATE INDEX idx_eval_id_evr ON public.ai_workflow_eval_result_responses(eval_id);
CREATE INDEX idx_response_id_evr ON public.ai_workflow_eval_result_responses(response_id);


CREATE TABLE public.ai_workflow_trigger_result_responses(
    trigger_result_id BIGINT NOT NULL PRIMARY KEY,
    workflow_result_id BIGINT NOT NULL REFERENCES ai_workflow_analysis_results(workflow_result_id),
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    response_id int8 NOT NULL REFERENCES completion_responses(response_id)
);
CREATE INDEX idx_workflow_result_id_trr ON public.ai_workflow_trigger_result_responses(workflow_result_id);
CREATE INDEX idx_trigger_id_trr ON public.ai_workflow_trigger_result_responses(trigger_id);
CREATE INDEX idx_response_id_trr ON public.ai_workflow_trigger_result_responses(response_id);
