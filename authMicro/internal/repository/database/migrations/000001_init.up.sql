CREATE TABLE registration_session
(
    id           UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    code         TEXT      NOT NULL,
    email        TEXT      NOT NULL UNIQUE,
    code_expires TIMESTAMP NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE login_session
(
    id           UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    account_id   UUID,
    email        TEXT      NOT NULL UNIQUE,
    code         TEXT      NOT NULL,
    code_expires TIMESTAMP NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_token_session
(
    id            UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    account_id    UUID      NOT NULL,
    refresh_token TEXT      NOT NULL,
    expires_at    TIMESTAMP NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW()
);
