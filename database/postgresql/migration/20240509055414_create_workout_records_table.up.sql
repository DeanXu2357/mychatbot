CREATE TABLE records
(
    id         BIGSERIAL PRIMARY KEY,
    reps       SMALLINT    NOT NULL DEFAULT 0,
    weight     SMALLINT    NOT NULL DEFAULT 0,
    event_id   BIGINT      NOT NULL,
    user_id    BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
