CREATE TABLE public.ai_trigger_actions_api(
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    retrieval_id BIGINT NOT NULL REFERENCES ai_retrieval_library(retrieval_id),
    PRIMARY KEY (trigger_id, retrieval_id)
);
CREATE INDEX ai_trigger_actions_api_trg_indx ON public.ai_trigger_actions_api("trigger_id");
CREATE INDEX ai_trigger_actions_api_ret_indx ON public.ai_trigger_actions_api("retrieval_id");