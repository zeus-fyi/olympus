CREATE ROLE "reader" LOGIN PASSWORD 'postgres';
GRANT "reader" TO "pg_read_all_data";
GRANT "reader" TO "pg_read_all_settings";
GRANT "reader" TO "pg_read_all_stats";