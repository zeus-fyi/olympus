BEGIN;
INSERT INTO "public"."key_types" VALUES (0, 'age');
INSERT INTO "public"."key_types" VALUES (1, 'gpg');
INSERT INTO "public"."key_types" VALUES (2, 'pgp');
INSERT INTO "public"."key_types" VALUES (3, 'ecdsa');
INSERT INTO "public"."key_types" VALUES (4, 'bls');
INSERT INTO "public"."key_types" VALUES (5, 'bearer');
INSERT INTO "public"."key_types" VALUES (6, 'jwt');
COMMIT;