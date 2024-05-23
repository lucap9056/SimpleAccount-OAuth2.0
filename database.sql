CREATE TABLE ThirdApp (
    id              CHAR(36),
    name            VARCHAR(255)     NOT NULL,
    client          INT             NOT NULL,
    salt            VARCHAR(64)     NOT NULL          COMMENT 'secret salt',
    hash            VARCHAR(64)     NOT NULL          COMMENT 'secret hash',
    callback        VARCHAR(255)    NOT NULL,
    description     TEXT            DEFAULT ''        COMMENT 'app description',
    permissions     INT             NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (client) REFERENCES User(id)
);