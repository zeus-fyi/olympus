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

CREATE TYPE network AS ENUM ('Mainnet', 'Prater');

-- ----------------------------
-- Table structure for validators
-- ----------------------------
CREATE TABLE "public"."validators" (
   "index" int4 NOT NULL,
   "pubkey" char(98) COLLATE "pg_catalog"."default" NOT NULL,
   "network" network NOT NULL DEFAULT 'Mainnet'
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
