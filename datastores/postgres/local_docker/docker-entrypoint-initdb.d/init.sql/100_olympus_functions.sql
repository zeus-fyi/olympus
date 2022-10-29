
-- since unix microseconds add 24 additional bits, this is good enough for distributed id generation, is sortable,
-- and has the best performance with ordering at the tail end since it doesn't rely on entropy to generate
-- the last bits, highly unlikely to collide
CREATE OR REPLACE FUNCTION next_id(OUT result int8) AS $$
DECLARE
    unix_time_utc bit(64) := EXTRACT('epoch' from NOW() at TIME ZONE ('UTC'))::int8::bit(64);
    unix_utc_now_microseconds bit(64) := FLOOR(EXTRACT('microseconds' FROM clock_timestamp()))::int8::bit(64);
    left_shifted_utc bit(64) := unix_time_utc << 40;
BEGIN
    result := (left_shifted_utc | unix_utc_now_microseconds)::int8;
END;
$$ LANGUAGE PLPGSQL;

CREATE EXTENSION hstore;
