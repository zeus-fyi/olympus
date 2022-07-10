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
CREATE TYPE validator_status AS ENUM ('pending', 'active', 'exited', 'withdrawal', 'unknown');
CREATE TYPE validator_substatus AS ENUM ('pending_initialized', 'pending_queued', 'active_ongoing', 'active_exiting', 'active_slashed','exited_unslashed', 'exited_slashed', 'withdrawal_possible', 'withdrawal_done', 'unknown');

-- ----------------------------
-- Table structure for validators
-- ----------------------------
CREATE TABLE "public"."validators" (
   "index" int4 NOT NULL,
   "balance" int8,
   "effective_balance" int8,
   "activation_eligibility_epoch" numeric CHECK (activation_eligibility_epoch <= activation_epoch),
   "activation_epoch" numeric CHECK (activation_epoch >= activation_eligibility_epoch),
   "exit_epoch" numeric,
   "withdrawable_epoch" numeric,
   "slashed" bool,
   "status" validator_status,
   "substatus" validator_substatus,
   "network" network NOT NULL DEFAULT 'mainnet',
   "pubkey" char(98) COLLATE "pg_catalog"."default" NOT NULL,
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


