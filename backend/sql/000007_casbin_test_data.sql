-- Добавление доступов согласно атрибутам в casbin
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'admin', 'disputes', '*');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'responsible_person', 'disputes', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'arbitr', 'disputes', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'complainant', 'disputes', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'responsible_person', 'revisions', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'admin', 'revisions', '*');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'admin', 'organizations', '*');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'arbitr', 'organizations', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'responsible_person', 'organizations', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'complainant', 'organizations', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'admin', 'revisions', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'admin', 'correspondences', 'get');
INSERT INTO common.casbinrule(ptype, v0, v1, v2) VALUES('p', 'responsible_person', 'correspondences', 'get');