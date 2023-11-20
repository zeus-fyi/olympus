CREATE TABLE public.ai_search_analysis_results(
    analysis_id int8 NOT NULL DEFAULT next_id(),
    response_id int8 NOT NULL REFERENCES completion_responses(response_id),
    search_hash text NOT NULL, -- hash of the search parameters and results
    analysis_hash text NOT NULL, -- hash of the search parameters, results, and response
    search_params JSONB NOT NULL DEFAULT '{}'::jsonb,
    search_results JSONB NOT NULL DEFAULT '{}'::jsonb,
    UNIQUE(analysis_hash)
);

ALTER TABLE "public"."ai_search_analysis_results" ADD CONSTRAINT "ai_search_analysis_results_pk" PRIMARY KEY ("analysis_id");
CREATE INDEX idx_sh ON public.ai_search_analysis_results("search_hash");

-- Create a composite index for the subject and contents columns to facilitate full text search
CREATE INDEX search_params_idx ON public.ai_search_analysis_results USING GIN (search_params);
CREATE INDEX search_results_idx ON public.ai_search_analysis_results USING GIN (search_results);
