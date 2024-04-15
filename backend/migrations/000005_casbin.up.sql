CREATE SCHEMA common;

CREATE TABLE common.casbinrule
(
   id    SERIAL       NOT NULL
      CONSTRAINT pk_casbinrule
            PRIMARY KEY,            -- request or policy definition
   ptype VARCHAR(255) NULL,         -- user attribute if p, user_id if d
   v0    VARCHAR(255) NULL,         -- resourse_name if p, user atribute if d
   v1    VARCHAR(255) NULL,         -- action name if p
   v2    VARCHAR(255) NULL,         -- reserved, need for framework
   v3    VARCHAR(255) NULL,         -- reserved, need for framework
   v4    VARCHAR(255) NULL,         -- reserved, need for framework
   v5    VARCHAR(255) NULL          -- reserved, need for framework
);