CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE public.ai_incoming_telegram_msgs (
          telegram_msg_id int8 NOT NULL DEFAULT next_id(),
          "org_id" int8 NOT NULL REFERENCES orgs(org_id),
          "user_id" int8 NOT NULL REFERENCES users(user_id),
          timestamp int8 NOT NULL,
          chat_id int8 NOT NULL,
          message_id int8 NOT NULL,
          sender_id int8 NOT NULL,
          group_name TEXT NOT NULL,
          message_text TEXT NOT NULL,
          metadata JSONB,
          active bool NOT NULL DEFAULT true,
          UNIQUE(chat_id, message_id)
);

ALTER TABLE "public"."ai_incoming_telegram_msgs" ADD CONSTRAINT "ai_incoming_telegram_msgs_pk" PRIMARY KEY ("telegram_msg_id");
CREATE INDEX idx_gn ON public.ai_incoming_telegram_msgs("group_name");
CREATE INDEX idx_ts ON public.ai_incoming_telegram_msgs("timestamp");
CREATE INDEX idx_ci ON public.ai_incoming_telegram_msgs("chat_id");
CREATE INDEX idx_mi ON public.ai_incoming_telegram_msgs("message_id");
CREATE INDEX idx_oi ON public.ai_incoming_telegram_msgs("org_id");
CREATE INDEX idx_ui ON public.ai_incoming_telegram_msgs("user_id");

-- Create a composite index for the subject and contents columns to facilitate full text search
ALTER TABLE public.ai_incoming_telegram_msgs ADD COLUMN message_text_tsvector tsvector GENERATED ALWAYS AS (to_tsvector('english', message_text)) STORED;
CREATE INDEX telegram_message_text_idx ON public.ai_incoming_telegram_msgs USING GIN (to_tsvector('english', message_text));

CREATE INDEX metadata_idx ON public.ai_incoming_telegram_msgs USING GIN (metadata);
CREATE INDEX idx_group_name_trgm ON public.ai_incoming_telegram_msgs USING GIN (group_name gin_trgm_ops);
CREATE INDEX ai_incoming_telegram_msgs_active_idx ON public.ai_incoming_telegram_msgs(active);

