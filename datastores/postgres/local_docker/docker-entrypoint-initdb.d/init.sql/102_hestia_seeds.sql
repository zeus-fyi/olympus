BEGIN;
INSERT INTO "public"."key_types" VALUES (0, 'age');
INSERT INTO "public"."key_types" VALUES (1, 'gpg');
INSERT INTO "public"."key_types" VALUES (2, 'pgp');
INSERT INTO "public"."key_types" VALUES (3, 'ecdsa');
INSERT INTO "public"."key_types" VALUES (4, 'bls');
INSERT INTO "public"."key_types" VALUES (5, 'bearer');
INSERT INTO "public"."key_types" VALUES (6, 'jwt');
INSERT INTO "public"."key_types" VALUES (7, 'password');
INSERT INTO "public"."key_types" VALUES (8, 'session');
INSERT INTO "public"."key_types" VALUES (9, 'verifyEmail');
INSERT INTO "public"."key_types" VALUES (10, 'resetPassword');
INSERT INTO "public"."key_types" VALUES (11, 'stripeCustomerID');
COMMIT;

BEGIN;
INSERT INTO "public"."services" VALUES (1677100016195486976, 'zeus');
INSERT INTO "public"."services" VALUES (1677096782693758000, 'ethereumEphemeryValidators');
INSERT INTO "public"."services" VALUES (1677096791420465000, 'ethereumMainnetValidators');

INSERT INTO "public"."services" VALUES (10, 'iris');
INSERT INTO "public"."services" VALUES (11, 'quickNodeMarketPlace');

COMMIT;