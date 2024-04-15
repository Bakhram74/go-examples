CREATE SCHEMA tarificator;

CREATE TABLE tarificator.productivity
(
    productivity_id UUID                                                           NOT NULL
        CONSTRAINT pk_productivity
            PRIMARY KEY,
    user_id         UUID REFERENCES subjects.user(user_id)                         NOT NULL,
    dispute_id      UUID REFERENCES disputes.dispute(dispute_id) ON DELETE CASCADE NOT NULL,
    in_role         VARCHAR                                                        NOT NULL, -- arbitr / responsible_person
    created_at      TIMESTAMP WITHOUT TIME ZONE                                    NOT NULL
);

CREATE INDEX ix_productivity_user_id
    ON tarificator.productivity(user_id);

CREATE TABLE tarificator.penaltyerror
(
    penaltyerror_id  INTEGER NOT NULL
        CONSTRAINT pk_penaltyerror
            PRIMARY KEY,
    penaltyerror_val VARCHAR NOT NULL,
    is_del           BOOLEAN NOT NULL
);

CREATE TABLE tarificator.penalty
(
    penalty_id      UUID                                                           NOT NULL
        CONSTRAINT pk_penalty
            PRIMARY KEY,
    user_id         UUID REFERENCES subjects.user(user_id)                         NOT NULL,
    dispute_id      UUID REFERENCES disputes.dispute(dispute_id) ON DELETE CASCADE NOT NULL,
    penaltyerror_id INTEGER REFERENCES tarificator.penaltyerror(penaltyerror_id)   NOT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE                                    NOT NULL
);

CREATE INDEX ix_penalty_user_id
    ON tarificator.penalty(user_id);