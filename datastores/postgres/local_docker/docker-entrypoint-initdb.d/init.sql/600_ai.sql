CREATE TABLE "public"."ai_incoming_email_tasks" (
    "email_id"  int8 NOT NULL DEFAULT next_id(),
    "msg_id" int8 NOT NULL,
    "from" text NOT NULL,
    "subject" text NOT NULL DEFAULT 'empty',
    "contents" text NOT NULL DEFAULT 'empty'
);

ALTER TABLE "public"."ai_incoming_email_tasks" ADD CONSTRAINT "ai_incoming_email_tasks_pk" PRIMARY KEY ("email_id");

-- Create a unique index on the msg_id column
CREATE UNIQUE INDEX idx_msg_id ON public.ai_incoming_email_tasks(msg_id);

-- Create an index on the from column
CREATE INDEX idx_from ON public.ai_incoming_email_tasks("from");

-- Create a composite index for the subject and contents columns to facilitate full text search
CREATE INDEX idx_subject_contents ON public.ai_incoming_email_tasks USING GIN (to_tsvector('english', subject || ' ' || contents));

CREATE TABLE "public"."ai_outgoing_email_tasks" (
    "response_id" int8 NOT NULL DEFAULT next_id(),
    "email_id" int8 NOT NULL,
    "msg_id" int8 NOT NULL,
    "from" text NOT NULL,
    "subject" text NOT NULL DEFAULT 'empty',
    "contents" text NOT NULL DEFAULT 'empty'
);

ALTER TABLE "public"."ai_outgoing_email_tasks" ADD CONSTRAINT "ai_outgoing_email_tasks_pk" PRIMARY KEY ("response_id");

-- Add a foreign key constraint to email_id
ALTER TABLE "public"."ai_outgoing_email_tasks" ADD CONSTRAINT "ai_outgoing_email_tasks_email_id_fk" FOREIGN KEY ("email_id")
    REFERENCES "public"."ai_incoming_email_tasks" ("email_id");

-- Add a foreign key constraint to msg_id, assuming msg_id in ai_incoming_email_tasks is unique and of type text
ALTER TABLE "public"."ai_outgoing_email_tasks" ADD CONSTRAINT "ai_outgoing_email_tasks_msg_id_fk" FOREIGN KEY ("msg_id")
    REFERENCES "public"."ai_incoming_email_tasks" ("msg_id");

-- Create a GIN index for the subject and contents columns in the ai_outgoing_email_tasks table for full text search
CREATE INDEX idx_subject_contents_out ON public.ai_outgoing_email_tasks USING GIN (to_tsvector('english', subject || ' ' || contents));