/*
 Navicat PostgreSQL Data Transfer

 Source Server         : zeus-do
 Source Server Type    : PostgreSQL
 Source Server Version : 140004
 Source Host           : private-db-postgresql-sfo3-zeus-do-user-9587342-0.b.db.ondigitalocean.com:25060
 Source Catalog        : zeus
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140004
 File Encoding         : 65001

 Date: 07/07/2022 20:27:04
*/

CREATE TYPE network AS ENUM ('mainnet', 'prater');

-- ----------------------------
-- Table structure for validators
-- ----------------------------
CREATE TABLE "public"."validators" (
   "index" int4 NOT NULL,
   "balance" int8,
   "effective_balance" int8,
   "activation_eligibility_epoch" int8 CHECK (activation_eligibility_epoch <= activation_epoch) NOT NULL DEFAULT (9223372036854775807),
   "activation_epoch" int8 CHECK (activation_epoch >= activation_eligibility_epoch) NOT NULL DEFAULT (9223372036854775807),
   "exit_epoch" int8 NOT NULL DEFAULT (9223372036854775807),
   "withdrawable_epoch" int8 NOT NULL DEFAULT (9223372036854775807),
   "updated_at" timestamptz  NOT NULL DEFAULT NOW(),
   "slashed" bool NOT NULL DEFAULT false,
   "pubkey" text NOT NULL CHECK(LENGTH(pubkey)=98),
   "status" text CHECK (status IN ('pending', 'active', 'exited', 'withdrawal', 'unknown')) NOT NULL DEFAULT 'unknown',
   "substatus" text CHECK (substatus IN ('pending_initialized', 'pending_queued', 'active_ongoing', 'active_exiting', 'active_slashed','exited_unslashed', 'exited_slashed', 'withdrawal_possible', 'withdrawal_done', 'unknown')) DEFAULT 'unknown',
   "network" network NOT NULL DEFAULT 'mainnet',
   "withdrawal_credentials" text
)
;
ALTER TABLE "public"."validators" OWNER TO "doadmin";

-- ----------------------------
-- Indexes structure for table validators
-- ----------------------------
CREATE UNIQUE INDEX "pubkey_index" ON "public"."validators" USING btree (
    "pubkey" COLLATE "pg_catalog"."default" "pg_catalog"."bpchar_ops" ASC NULLS LAST
    );

-- ----------------------------
-- Primary Key structure for table validators
-- ----------------------------
ALTER TABLE "public"."validators" ADD CONSTRAINT "validators_pkey" PRIMARY KEY ("index");
CREATE INDEX "last_updated_at_index" ON "public"."validators" (updated_at ASC);

