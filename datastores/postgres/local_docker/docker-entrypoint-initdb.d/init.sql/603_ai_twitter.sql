CREATE TABLE public.ai_twitter_search_query (
    search_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "user_id" int8 NOT NULL REFERENCES users(user_id),
    search_group_name TEXT NOT NULL,
    max_results int8 NOT NULL CHECK (max_results >= 10 AND max_results <= 100),
    query TEXT NOT NULL,
    active bool NOT NULL DEFAULT true
);
CREATE INDEX tw_search_group_name_trgm_idx ON public.ai_twitter_search_query USING GIN (search_group_name gin_trgm_ops);
-- Add tsvector column and index for full-text search on query
ALTER TABLE public.ai_twitter_search_query ADD COLUMN query_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', query)) STORED;
CREATE INDEX tw_query_tsvector_idx ON public.ai_twitter_search_query USING GIN (query_tsvector);
ALTER TABLE "public"."ai_twitter_search_query" ADD CONSTRAINT "ai_twitter_search_query_uniq" UNIQUE ("org_id", "user_id", "query");
CREATE INDEX ai_twitter_search_query_active_idx ON public.ai_twitter_search_query (active);

-- AI Incoming Tweets Table
CREATE TABLE public.ai_incoming_tweets (
    tweet_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    search_id int8 NOT NULL,
    message_text TEXT NOT NULL,
    FOREIGN KEY (search_id) REFERENCES public.ai_twitter_search_query(search_id)
);

-- Create a composite index for the subject and contents columns to facilitate full text search
ALTER TABLE public.ai_incoming_tweets ADD COLUMN message_text_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', message_text)) STORED;
CREATE INDEX tw_message_text_idx ON public.ai_incoming_tweets USING GIN (to_tsvector('english', message_text));
CREATE INDEX tw_search_id_idx ON public.ai_incoming_tweets (search_id);
CREATE INDEX tweet_id_desc_idx ON public.ai_incoming_tweets (tweet_id DESC);
