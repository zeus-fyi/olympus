CREATE TABLE public.ai_workflow_io(
    response_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    approval_id BIGINT NOT NULL REFERENCES ai_trigger_actions_approvals(approval_id),
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    retrieval_id BIGINT NOT NULL REFERENCES ai_retrieval_library(retrieval_id),
    req_payload JSONB,
    resp_payload JSONB
);

CREATE INDEX ai_trigger_actions_api_trg_resp_indx ON public.ai_trigger_actions_api_reqs_responses("trigger_id");
CREATE INDEX ai_trigger_actions_api_resp_rest_indx ON public.ai_trigger_actions_api_reqs_responses("retrieval_id");
CREATE INDEX ai_trigger_actions_api_apprv_indx ON public.ai_trigger_actions_api_reqs_responses("approval_id");

ALTER TABLE public.ai_trigger_actions_api_reqs_responses
    ADD CONSTRAINT ai_trigger_actions_api_responses_ret_apprv_uniq UNIQUE (response_id, approval_id, trigger_id, retrieval_id);