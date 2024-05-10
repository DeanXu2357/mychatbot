CREATE TABLE events
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    user_id    BIGINT      NOT NULL,
    tags       VARCHAR(20)[10],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
