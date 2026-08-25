CREATE TABLE IF NOT EXISTS "outbox_events" (
    "id" uuid NOT NULL,
    "event_type" varchar(100) NOT NULL,
    "payload" jsonb NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "published_at" timestamptz,
    "attempts" integer NOT NULL DEFAULT 0,
    "next_attempt_at" timestamptz NOT NULL DEFAULT now(),
    "locked_until" timestamptz,
    "last_error" text,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "idx_outbox_events_pending"
ON "outbox_events" ("next_attempt_at", "created_at")
WHERE "published_at" IS NULL;
