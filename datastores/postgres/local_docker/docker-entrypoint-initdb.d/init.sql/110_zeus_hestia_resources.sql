CREATE TABLE resources (
   resource_id int8 NOT NULL DEFAULT next_id(),
   type text NOT NULL,
   PRIMARY KEY (resource_id)
);

CREATE TABLE nodes (
   resource_id int8 NOT NULL REFERENCES resources(resource_id),
   description text NOT NULL,
   slug text NOT NULL,
   memory int NOT NULL,
   memory_units text NOT NULL,
   vcpus int NOT NULL,
   disk int NOT NULL,
   disk_units text NOT NULL,
   price_monthly float8 NOT NULL,
   price_hourly float8 NOT NULL,
   region text NOT NULL,
   cloud_provider text NOT NULL,
   PRIMARY KEY (resource_id)
);

CREATE INDEX nodes_region_idx ON nodes (region);
CREATE INDEX nodes_vcpus_idx ON nodes (vcpus);
CREATE INDEX nodes_memory_idx ON nodes (memory);
CREATE INDEX nodes_cloud_provider_idx ON nodes (cloud_provider);
CREATE INDEX nodes_price_monthly_idx ON nodes (price_monthly);
CREATE INDEX nodes_price_hourly_idx ON nodes (price_hourly);

CREATE TABLE disks (
   resource_id int8 NOT NULL REFERENCES resources(resource_id),
   description text NOT NULL,
   type text NOT NULL,
   disk_units text NOT NULL,
   disk_size int NOT NULL,
   price_monthly float8 NOT NULL,
   price_hourly float8 NOT NULL,
   region text NOT NULL,
   cloud_provider text NOT NULL,
   PRIMARY KEY (resource_id)
);

CREATE INDEX disks_cloud_provider_idx ON disks (region);
CREATE INDEX disks_region_idx ON disks (cloud_provider);

CREATE TABLE org_resources (
   org_resource_id int8 NOT NULL DEFAULT next_id(),
   resource_id int8 NOT NULL REFERENCES resources(resource_id),
   org_id int8 NOT NULL REFERENCES orgs(org_id),
   quantity float8 NOT NULL,
   begin_service timestamptz NOT NULL DEFAULT NOW(),
   end_service timestamptz DEFAULT NULL,
   free_trial boolean NOT NULL DEFAULT FALSE,
   PRIMARY KEY (org_resource_id)
);

CREATE INDEX begin_resource_idx ON org_resources (begin_service);
CREATE INDEX end_resource_idx ON org_resources (end_service);

CREATE TABLE digitalocean_node_pools (
    org_resource_id int8 NOT NULL REFERENCES org_resources(org_resource_id),
    resource_id int8 NOT NULL REFERENCES resources(resource_id),
    node_pool_id text NOT NULL,
    node_context_id text NOT NULL,
    PRIMARY KEY (org_resource_id)
);

CREATE TABLE org_resources_cloud_ctx (
   org_resource_id int8 NOT NULL REFERENCES org_resources(org_resource_id),
   cloud_ctx_ns_id int8 NOT NULL REFERENCES topologies_org_cloud_ctx_ns(cloud_ctx_ns_id),
   PRIMARY KEY (org_resource_id)
);
CREATE INDEX org_resources_cloud_ctx_id_idx ON org_resources_cloud_ctx (cloud_ctx_ns_id);
