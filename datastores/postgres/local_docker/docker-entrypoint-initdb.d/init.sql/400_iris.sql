CREATE TABLE "public"."org_routes" (
    "route_id" int8 NOT NULL DEFAULT next_id(),
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "route_path" text NOT NULL
);
ALTER TABLE "public"."org_routes" ADD CONSTRAINT "org_routes_pk" PRIMARY KEY ("route_id");
CREATE INDEX org_routes_path_ind ON org_routes ("org_id", "route_path");
ALTER TABLE "public"."org_routes" ADD CONSTRAINT "org_routes_uniq" UNIQUE ("org_id", "route_path");
