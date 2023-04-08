CREATE TABLE nodes (
   node_id int8 NOT NULL DEFAULT next_id(),
   description text NOT NULL,
   slug text NOT NULL,
   memory int NOT NULL,
   vcpus int NOT NULL,
   disk int NOT NULL,
   price_monthly float8 NOT NULL,
   price_hourly float8 NOT NULL,
   region text NOT NULL,
   cloud_provider text NOT NULL,
   PRIMARY KEY (node_id)
);

CREATE INDEX nodes_region_idx ON nodes (region);
CREATE INDEX nodes_vcpus_idx ON nodes (vcpus);
CREATE INDEX nodes_memory_idx ON nodes (memory);
CREATE INDEX nodes_cloud_provider_idx ON nodes (cloud_provider);
CREATE INDEX nodes_price_monthly_idx ON nodes (price_monthly);
CREATE INDEX nodes_price_hourly_idx ON nodes (price_hourly);

CREATE TABLE disks (
   disk_id int8 NOT NULL DEFAULT next_id(),
   description text NOT NULL,
   type text NOT NULL,
   units text NOT NULL,
   size int NOT NULL,
   price_monthly float8 NOT NULL,
   price_hourly float8 NOT NULL,
   region text NOT NULL,
   cloud_provider text NOT NULL,
   PRIMARY KEY (disk_id)
);
