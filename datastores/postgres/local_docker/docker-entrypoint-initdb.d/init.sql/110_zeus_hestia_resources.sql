CREATE TABLE nodes (
   node_id int8 NOT NULL DEFAULT next_id(),
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