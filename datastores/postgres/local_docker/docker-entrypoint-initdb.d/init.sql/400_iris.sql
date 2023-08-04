CREATE TABLE "public"."org_routes" (
    "route_id" int8 NOT NULL DEFAULT next_id(),
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "route_path" text NOT NULL
);
ALTER TABLE "public"."org_routes" ADD CONSTRAINT "org_routes_pk" PRIMARY KEY ("route_id");
CREATE INDEX org_routes_path_ind ON org_routes ("org_id", "route_path");
ALTER TABLE "public"."org_routes" ADD CONSTRAINT "org_routes_uniq" UNIQUE ("org_id", "route_path");

CREATE TABLE "public"."org_route_groups" (
    "route_group_id" int8 NOT NULL DEFAULT next_id(),
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "auto_generated" bool NOT NULL DEFAULT false,
    "route_group_name" text NOT NULL
);
ALTER TABLE "public"."org_route_groups" ADD CONSTRAINT "org_route_groups_pk" PRIMARY KEY ("route_group_id");
CREATE INDEX org_route_groups_path_ind ON org_route_groups ("org_id", "route_group_name");

CREATE TABLE "public"."org_routes_groups" (
    "route_group_id" int8  NOT NULL REFERENCES org_route_groups(route_group_id),
    "route_id" int8  NOT NULL REFERENCES org_routes(route_id)
);
ALTER TABLE "public"."org_routes_groups" ADD CONSTRAINT "org_routes_groups_pk" PRIMARY KEY ("route_group_id", "route_id");

