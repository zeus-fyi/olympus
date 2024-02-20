CREATE TABLE public.user_entities (
    entity_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    nickname TEXT NOT NULL,
    platform TEXT NOT NULL,
    first_name TEXT NULL,
    last_name TEXT NULL
);
CREATE UNIQUE INDEX user_entities_uniq_platform_nickname_uniq ON public.user_entities(org_id,nickname, platform);
CREATE INDEX user_entities_fn_idx ON public.user_entities(first_name);
CREATE INDEX user_entities_ln_idx ON public.user_entities(last_name);
CREATE INDEX user_entities_org_id_idx ON public.user_entities(org_id);

CREATE TABLE public.user_entities_md (
    entity_metadata_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    entity_id BIGINT NOT NULL REFERENCES public.user_entities(entity_id),
    json_data JSONB,
    text_data TEXT
);

CREATE TABLE public.user_entities_md_labels (
    entity_metadata_label_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    entity_metadata_id BIGINT NOT NULL REFERENCES public.user_entities_md(entity_metadata_id),
    label TEXT NOT NULL
);

CREATE INDEX labels_idx ON public.user_entities_md_labels(label);
CREATE INDEX label_md_idx ON public.user_entities_md_labels(entity_metadata_id, label);
CREATE UNIQUE INDEX user_entities_md_labels_uniq ON public.user_entities_md_labels(entity_metadata_id, label);
