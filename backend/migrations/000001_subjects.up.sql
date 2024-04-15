CREATE SCHEMA subjects;

CREATE TABLE subjects.organization
(
    organization_id    UUID       NOT NULL
        CONSTRAINT pk_organization
            PRIMARY KEY,
    organization_title VARCHAR    NOT NULL,
    organization_code  VARCHAR(3) NOT NULL, -- WBS (WB Supplier), WBO (WB Office)
    inner_id           BIGINT     NOT NULL, -- office_id or supplier_id
    is_del             BOOLEAN    NOT NULL
);

CREATE UNIQUE INDEX uq_organization_inner_id ON subjects.organization (inner_id, organization_code);


CREATE TABLE subjects.user
(
    user_id               UUID                                                   NOT NULL
        CONSTRAINT pk_user
            PRIMARY KEY,
    employee_id           INTEGER                                                NULL,
    organization_id       UUID REFERENCES subjects.organization(organization_id) NOT NULL,
    fullname              VARCHAR                                                NOT NULL,
    phone                 VARCHAR                                                NOT NULL,
    role                  VARCHAR                                                NOT NULL, -- admin, responsible_person, arbitr, complainant
    created_at            TIMESTAMP WITHOUT TIME ZONE                            NOT NULL,
    avatar_url            VARCHAR                                                NULL,
    is_del                BOOLEAN                                                NOT NULL
);

CREATE UNIQUE INDEX uq_user_employee_id ON subjects.user (employee_id);
CREATE UNIQUE INDEX uq_user_phone ON subjects.user (phone);
