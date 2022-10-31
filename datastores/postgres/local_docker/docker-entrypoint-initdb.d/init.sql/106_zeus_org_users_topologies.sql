-- user registered topologies
CREATE TABLE "public"."org_users_topologies" (
    "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "user_id" int8 NOT NULL REFERENCES users(user_id)
);
