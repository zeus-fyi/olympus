BEGIN;
INSERT INTO "public"."key_types" VALUES (0, 'age');
INSERT INTO "public"."key_types" VALUES (1, 'gpg');
INSERT INTO "public"."key_types" VALUES (2, 'pgp');
INSERT INTO "public"."key_types" VALUES (3, 'ecdsa');
INSERT INTO "public"."key_types" VALUES (4, 'bls');
INSERT INTO "public"."key_types" VALUES (5, 'bearer');
INSERT INTO "public"."key_types" VALUES (6, 'jwt');
COMMIT;

BEGIN;
INSERT INTO "public"."services" VALUES (1677100016195486976, 'zeus');
INSERT INTO "public"."services" VALUES (1677096782693758000, 'ethereumEphemeryValidators');
INSERT INTO "public"."services" VALUES (1677096791420465000, 'ethereumMainnetValidators');
COMMIT;