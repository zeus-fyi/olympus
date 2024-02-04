CREATE TABLE ext_cluster_configs (
    ext_config_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    cloud_provider    TEXT NOT NULL CHECK ( cloud_provider IN ('aws', 'gcp', 'ovh', 'do', 'custom') ),
    region            TEXT NOT NULL,
    context           TEXT NOT NULL,
    context_alias     TEXT NOT NULL,
    env               TEXT NOT NULL DEFAULT 'test'
);

ALTER TABLE public.ext_cluster_configs
    ADD CONSTRAINT cloud_ctx_ns_org_uniq UNIQUE  ("org_id", "cloud_provider","region", "context");