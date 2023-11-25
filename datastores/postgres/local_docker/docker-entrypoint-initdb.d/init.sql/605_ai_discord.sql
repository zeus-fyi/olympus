-- Table for Discord search query
CREATE TABLE public.ai_discord_search_query (
    search_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    search_group_name TEXT NOT NULL,
    max_results BIGINT NOT NULL CHECK (max_results <= 100),
    query TEXT NOT NULL
);
CREATE INDEX discord_search_group_name_trgm_idx ON public.ai_discord_search_query USING GIN (search_group_name gin_trgm_ops);
ALTER TABLE public.ai_discord_search_query ADD COLUMN query_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', query)) STORED;
CREATE INDEX discord_query_tsvector_idx ON public.ai_discord_search_query USING GIN (query_tsvector);
ALTER TABLE "public"."ai_discord_search_query" ADD CONSTRAINT "ai_discord_search_query_uniq" UNIQUE ("org_id", "user_id", "query");

-- Table for Discord guilds
CREATE TABLE public.ai_discord_guild (
    guild_id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE INDEX rd_guild_name_trgm_idx ON public.ai_discord_guild USING GIN (name gin_trgm_ops);

-- Table for Discord channels
CREATE TABLE public.ai_discord_channel (
    search_id BIGINT NOT NULL REFERENCES public.ai_discord_search_query(search_id),
    guild_id TEXT NOT NULL REFERENCES public.ai_discord_guild(guild_id),
    channel_id TEXT NOT NULL PRIMARY KEY,
    category_id TEXT,
    category TEXT NOT NULL,
    name TEXT NOT NULL,
    topic TEXT NOT NULL
);
ALTER TABLE public.ai_discord_channel ADD COLUMN name_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', name)) STORED;
ALTER TABLE public.ai_discord_channel ADD COLUMN topic_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', topic)) STORED;
ALTER TABLE public.ai_discord_channel ADD COLUMN category_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', category)) STORED;
CREATE INDEX idx_ai_discord_channel_guild_fk ON public.ai_discord_channel (guild_id);
CREATE INDEX idx_ai_discord_channel_search_fk ON public.ai_discord_channel (search_id);

CREATE INDEX rd_channel_name_tsvector_idx ON public.ai_discord_channel USING GIN (name_tsvector);
CREATE INDEX rd_channel_topic_tsvector_idx ON public.ai_discord_channel USING GIN (topic_tsvector);
CREATE INDEX rd_channel_category_tsvector_idx ON public.ai_discord_channel USING GIN (category_tsvector);

-- Table for incoming Discord messages
CREATE TABLE public.ai_incoming_discord_messages (
    message_id TEXT NOT NULL PRIMARY KEY,
    timestamp_creation BIGINT NOT NULL,
    search_id BIGINT NOT NULL REFERENCES public.ai_discord_search_query(search_id),
    guild_id TEXT NOT NULL REFERENCES public.ai_discord_guild(guild_id),
    channel_id TEXT NOT NULL REFERENCES public.ai_discord_channel(channel_id),
    author JSONB NOT NULL,
    content TEXT NOT NULL,
    mentions JSONB,
    reactions JSONB,
    reference JSONB,
    timestamp_edited BIGINT NOT NULL DEFAULT 0,
    type TEXT NOT NULL
);

CREATE INDEX rd_message_tssearch_idx ON public.ai_incoming_discord_messages (timestamp_creation DESC);
CREATE INDEX rd_message_search_idx ON public.ai_incoming_discord_messages (search_id);
CREATE INDEX rd_message_type_idx ON public.ai_incoming_discord_messages (type);
ALTER TABLE public.ai_incoming_discord_messages ADD COLUMN content_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', content)) STORED;
CREATE INDEX rd_content_tsvector_idx ON public.ai_incoming_discord_messages USING GIN (content_tsvector);

CREATE INDEX idx_ai_incoming_discord_messages_guild_id ON public.ai_incoming_discord_messages (guild_id);
CREATE INDEX idx_ai_incoming_discord_messages_channel_id ON public.ai_incoming_discord_messages (channel_id);
CREATE INDEX idx_ai_incoming_discord_messages_author
    ON public.ai_incoming_discord_messages USING GIN (
                                                      (author -> 'name') jsonb_path_ops,
                                                      (author -> 'nickname') jsonb_path_ops,
                                                      (jsonb_path_query_array(author, '$.roles[*]')) jsonb_path_ops
        );
CREATE INDEX idx_ai_incoming_discord_messages_mentions
    ON public.ai_incoming_discord_messages USING GIN (
                                                      (jsonb_path_query_array(mentions, '$[*].name')) jsonb_path_ops,
                                                      (jsonb_path_query_array(mentions, '$[*].nickname')) jsonb_path_ops,
                                                      (jsonb_path_query_array(mentions, '$[*].roles[*]')) jsonb_path_ops
        );
CREATE INDEX idx_ai_incoming_discord_messages_reactions
    ON public.ai_incoming_discord_messages USING GIN (
                                                      (jsonb_path_query_array(reactions, '$[*].users[*].name')) jsonb_path_ops,
                                                      (jsonb_path_query_array(reactions, '$[*].users[*].nickname')) jsonb_path_ops,
                                                      (jsonb_path_query_array(reactions, '$[*].emoji.code')) jsonb_path_ops,
                                                      (jsonb_path_query_array(reactions, '$[*].count')) jsonb_path_ops
        );
