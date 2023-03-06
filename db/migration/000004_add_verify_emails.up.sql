CREATE TABLE "verify_emails"
(
    "id"                  bigserial PRIMARY KEY,
    "username"            varchar        NOT NULL,
    "email"               varchar        NOT NULL,
    "secret_code"         varchar UNIQUE NOT NULL,
    "is_used"             bool           NOT NULL DEFAULT false,
    "password_changed_at" timestamptz    NOT NULL DEFAULT '0001-01-01',
    "created_at"          timestamptz    NOT NULL DEFAULT (now()),
    "expired_at"          timestamptz    NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "users" ADD COLUMN "is_email_verified" bool NOT NULL DEFAULT FALSE;