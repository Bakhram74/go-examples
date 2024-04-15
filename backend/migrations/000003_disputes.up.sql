CREATE SCHEMA disputes;

CREATE TABLE disputes.dispute
(
    dispute_id                   UUID                                                                      NOT NULL
        CONSTRAINT pk_dispute
            PRIMARY KEY,
    shortage_id                  UUID REFERENCES shortages.shortage(shortage_id)                           NOT NULL,
    organization_id              UUID REFERENCES subjects.organization(organization_id)                    NOT NULL, -- офис, где рассматривается спор
    is_shortage_canceled         BOOLEAN                                                                   NOT NULL, -- анулировано ли списание
    is_dispute_reopened          BOOLEAN                                                                   NOT NULL,
    status                       VARCHAR                                                                   NOT NULL, -- opened, in_work, closed
    created_at                   TIMESTAMP WITHOUT TIME ZONE                                               NOT NULL,
    closed_at                    TIMESTAMP WITHOUT TIME ZONE                                               NULL,
    reopened_at                  TIMESTAMP WITHOUT TIME ZONE                                               NULL,
    is_arbitr_invited            BOOLEAN                                                                   NOT NULL
);

CREATE INDEX ix_dispute_shortage_id
    ON disputes.dispute(shortage_id);

CREATE INDEX ix_dispute_organization_id
    ON disputes.dispute(organization_id);


CREATE TABLE disputes.dispute_role
(
    user_id      UUID REFERENCES subjects.user(user_id)                         NOT NULL,
    dispute_id   UUID REFERENCES disputes.dispute(dispute_id) ON DELETE CASCADE NOT NULL,
    dispute_role VARCHAR                                                        NOT NULL, -- complainant, responsible_person, guilty_responsible_person, guilty_worker
    created_at   TIMESTAMP WITHOUT TIME ZONE                                    NOT NULL
);

CREATE INDEX ix_dispute_role_dispute_id
    ON disputes.dispute_role(dispute_id);


CREATE TABLE disputes.chat
(
    message_id      UUID                                                           NOT NULL
        CONSTRAINT pk_chat
            PRIMARY KEY,
    dispute_id      UUID REFERENCES disputes.dispute(dispute_id) ON DELETE CASCADE NOT NULL,
    sender_id       UUID REFERENCES subjects.user(user_id)                         NOT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE                                    NOT NULL,
    message_body    VARCHAR(400)                                                   NULL,
    attachment_path VARCHAR(90)                                                    NULL
);

CREATE INDEX ix_chat_dispute_id
    ON disputes.chat(dispute_id);


CREATE TABLE disputes.arbitrinvitation
(
    arbitr_invitation_id UUID                                                           NOT NULL
        CONSTRAINT pk_arbitrinvitation
            PRIMARY KEY,
    organization_id      UUID REFERENCES subjects.organization(organization_id)         NOT NULL,
    dispute_id           UUID REFERENCES disputes.dispute(dispute_id) ON DELETE CASCADE NOT NULL,
    status               VARCHAR                                                        NOT NULL, -- opened, in_work, closed
    created_at           TIMESTAMP WITHOUT TIME ZONE                                    NOT NULL,
    in_work_at           TIMESTAMP WITHOUT TIME ZONE                                    NULL,
    closed_at            TIMESTAMP WITHOUT TIME ZONE                                    NULL,
    arbitr_id            UUID REFERENCES subjects.user(user_id)                         NULL
);

CREATE UNIQUE INDEX uq_arbitrinvitation_dispute_id ON disputes.arbitrinvitation (dispute_id);

CREATE INDEX ix_arbitrinvitation_organization_id
    ON disputes.arbitrinvitation(organization_id);

CREATE INDEX ix_arbitrinvitation_arbitr_id
    ON disputes.arbitrinvitation(arbitr_id);


CREATE TABLE disputes.revision
(
    revision_id           UUID                                                           NOT NULL
        CONSTRAINT pk_revision
            PRIMARY KEY,
    organization_id      UUID REFERENCES subjects.organization(organization_id)         NOT NULL,
    dispute_id           UUID REFERENCES disputes.dispute(dispute_id) ON DELETE CASCADE NOT NULL,
    status               VARCHAR                                                        NOT NULL, -- opened, in_work, closed
    created_at           TIMESTAMP WITHOUT TIME ZONE                                    NOT NULL,
    in_work_at           TIMESTAMP WITHOUT TIME ZONE                                    NULL,
    closed_at            TIMESTAMP WITHOUT TIME ZONE                                    NULL,
    worker_id            UUID REFERENCES subjects.user(user_id)                         NULL -- кто взял в работу запрос
);

CREATE INDEX ix_revision_organization_id
    ON disputes.revision(organization_id);

CREATE INDEX ix_revision_worker_id
    ON disputes.revision(worker_id);

CREATE INDEX ix_revision_dispute_id
    ON disputes.revision(dispute_id);


CREATE TABLE disputes.correspondence
(
    correspondence_id UUID                                                             NOT NULL
        CONSTRAINT pk_correspondence
            PRIMARY KEY,
    revision_id       UUID REFERENCES disputes.revision(revision_id) ON DELETE CASCADE NOT NULL,
    sender_id         UUID REFERENCES subjects.user(user_id)                           NOT NULL,
    created_at        TIMESTAMP WITHOUT TIME ZONE                                      NOT NULL,
    message_body      VARCHAR(400)                                                     NULL,
    attachment_path   VARCHAR(90)                                                      NULL
);

CREATE INDEX ix_correspondence_revision_id
    ON disputes.correspondence(revision_id);
