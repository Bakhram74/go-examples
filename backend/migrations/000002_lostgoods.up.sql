CREATE SCHEMA shortages;


CREATE TABLE shortages.lostreason
(
    lostreason_id  INTEGER NOT NULL
        CONSTRAINT pk_lostreason
            PRIMARY KEY,
    lostreason_val VARCHAR NOT NULL,
    is_del         BOOLEAN NOT NULL
);


CREATE TABLE shortages.shortage
(
    shortage_id     UUID                                                   NOT NULL
        CONSTRAINT pk_shortage
            PRIMARY KEY,
    user_id         UUID REFERENCES subjects.user(user_id)                 NOT NULL, -- на кого отнесено списание
    goods_id        BIGINT                                                 NOT NULL,
    tare_id         BIGINT                                                 NOT NULL,
    tare_type       VARCHAR(3)                                             NOT NULL,
    lostreason_id   INTEGER REFERENCES shortages.lostreason(lostreason_id) NOT NULL,
    currency_code   INTEGER                                                NOT NULL,
    lost_amount     NUMERIC(11,2)                                          NOT NULL,
    is_disputed     BOOLEAN                                                NOT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE                            NOT NULL
);

CREATE INDEX ix_shortage_goods_id
    ON shortages.shortage(goods_id);
