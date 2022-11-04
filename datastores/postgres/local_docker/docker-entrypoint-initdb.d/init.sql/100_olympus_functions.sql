
-- since unix microseconds add 24 additional bits, this is good enough for distributed id generation, is sortable,
-- and has the best performance with ordering at the tail end since it doesn't rely on entropy to generate
-- the last bits, highly unlikely to collide
CREATE OR REPLACE FUNCTION "public"."next_id"()
    RETURNS "pg_catalog"."int8" AS
$BODY$
DECLARE
    unix_utc_now bigint := (SELECT (EXTRACT('epoch' from NOW() at TIME ZONE ('UTC'))*1000000000));
BEGIN
    RETURN unix_utc_now;
END;
$BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

CREATE EXTENSION hstore;
