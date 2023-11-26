CREATE TABLE public.ai_reddit_search_query (
    search_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    search_group_name TEXT NOT NULL,
    max_results BIGINT NOT NULL CHECK (max_results <= 100),
    query TEXT NOT NULL
);

CREATE INDEX rd_search_group_name_trgm_idx ON public.ai_reddit_search_query USING GIN (search_group_name gin_trgm_ops);
ALTER TABLE public.ai_reddit_search_query ADD COLUMN query_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', query)) STORED;
CREATE INDEX rd_query_tsvector_idx ON public.ai_reddit_search_query USING GIN (query_tsvector);
ALTER TABLE "public"."ai_reddit_search_query" ADD CONSTRAINT "ai_reddit_search_query_uniq" UNIQUE ("org_id", "user_id", "query");

CREATE TABLE public.ai_reddit_incoming_posts (
    post_id TEXT NOT NULL PRIMARY KEY,
    post_full_id TEXT NOT NULL,
    search_id BIGINT NOT NULL,
    permalink TEXT NOT NULL,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT,
    created_at BIGINT NOT NULL,
    edited_at BIGINT,
    score BIGINT,
    upvote_ratio FLOAT,
    number_of_comments BIGINT,
    reddit_meta JSONB NOT NULL,
    author TEXT NOT NULL,
    author_id TEXT NOT NULL,
    subreddit TEXT NOT NULL,
    FOREIGN KEY (search_id) REFERENCES public.ai_reddit_search_query(search_id)
);

-- Add tsvector columns and indexes for full-text search
ALTER TABLE public.ai_reddit_incoming_posts ADD COLUMN title_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', title)) STORED;
ALTER TABLE public.ai_reddit_incoming_posts ADD COLUMN body_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', body)) STORED;
CREATE INDEX idx_red_post_created_at ON public.ai_reddit_incoming_posts (created_at);
CREATE INDEX idx_subreddit ON public.ai_reddit_incoming_posts (subreddit);

CREATE INDEX rd_title_tsvector_idx ON public.ai_reddit_incoming_posts USING GIN (title_tsvector);
CREATE INDEX rd_body_tsvector_idx ON public.ai_reddit_incoming_posts USING GIN (body_tsvector);
CREATE INDEX rd_search_id_idx ON public.ai_reddit_incoming_posts (search_id);
CREATE INDEX rd_post_created_ts_idx ON public.ai_reddit_incoming_posts (created_at DESC);