CREATE TABLE gcp_services (
  "name" text NOT NULL,
  "service_id" text NOT NULL,
  "display_name" text NOT NULL,
  "business_entity_name" text NOT NULL
);

ALTER TABLE "public"."gcp_services" ADD CONSTRAINT "gcp_services_pk" PRIMARY KEY ("service_id");

CREATE TABLE gcp_services_skus (
    "service_id" text NOT NULL REFERENCES gcp_services("service_id"),
    "name" text NOT NULL,
    "sku_id" text NOT NULL,
    "description" text,
    "service_display_name" text,
    "resource_family" text,
    "resource_group" text,
    "usage_type" text,
    "service_regions" jsonb,
    "pricing_info" jsonb,
    "service_provider_name" text,
    "geo_taxonomy" jsonb,
    PRIMARY KEY ("sku_id")
);

CREATE INDEX gcp_services_skus_name_idx ON gcp_services_skus ("name");
CREATE INDEX gcp_services_skus_sku_id_idx ON gcp_services_skus ("sku_id");
CREATE INDEX gcp_services_skus_description_idx ON gcp_services_skus ("description");
CREATE INDEX gcp_services_skus_service_provider_name_idx ON gcp_services_skus ("service_provider_name");
CREATE INDEX gcp_services_skus_resource_group_idx ON gcp_services_skus ("resource_group");
CREATE INDEX gcp_services_skus_resource_family_idx ON gcp_services_skus ("resource_family");
CREATE INDEX gcp_services_skus_usage_type_idx ON gcp_services_skus ("usage_type");
CREATE INDEX gcp_services_skus_service_regions_idx ON gcp_services_skus USING gin ("service_regions");
CREATE INDEX gcp_services_skus_pricing_info_idx ON gcp_services_skus USING gin ("pricing_info");
