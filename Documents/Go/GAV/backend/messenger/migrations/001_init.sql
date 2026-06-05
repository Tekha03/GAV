CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS "chats" (
    "id" uuid DEFAULT gen_random_uuid(),
    "is_group" boolean NOT NULL DEFAULT false,
    "title" text,
    "photo_url" text,
    "created_at" timestamptz,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "chat_members" (
    "chat_id" uuid,
    "user_id" uuid,
    "role" varchar(20) NOT NULL DEFAULT 'member',
    "joined_at" timestamptz,
    "left_at" timestamptz,
    "muted" boolean NOT NULL DEFAULT false,
    "last_read_message_id" uuid,
    PRIMARY KEY ("chat_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "messages" (
    "id" uuid,
    "chat_id" uuid NOT NULL,
    "sender_id" uuid NOT NULL,
    "text" text,
    "reply_to_id" uuid,
    "created_at" timestamptz,
    "edited_at" timestamptz,
    "deleted_at" timestamptz,
    "read_at" timestamptz,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "idx_messages_chat_id" ON "messages" ("chat_id");

CREATE TABLE IF NOT EXISTS "attachments" (
    "id" uuid DEFAULT gen_random_uuid(),
    "message_id" uuid NOT NULL,
    "url" text NOT NULL,
    "type" varchar(20) NOT NULL,
    "file_name" text,
    "file_size" bigint NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "idx_attachments_message_id" ON "attachments" ("message_id");

CREATE TABLE IF NOT EXISTS "reactions" (
    "id" uuid DEFAULT gen_random_uuid(),
    "message_id" uuid,
    "user_id" uuid,
    "emoji" varchar(10) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "idx_reactions_user_id" ON "reactions" ("user_id");
CREATE INDEX IF NOT EXISTS "idx_reactions_message_id" ON "reactions" ("message_id");

CREATE TABLE IF NOT EXISTS "pinned_messages" (
    "chat_id" uuid,
    "message_id" uuid,
    "pinned_at" timestamptz,
    PRIMARY KEY ("chat_id", "message_id")
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_messages_attachments'
    ) THEN
        ALTER TABLE "attachments"
        ADD CONSTRAINT "fk_messages_attachments"
        FOREIGN KEY ("message_id") REFERENCES "messages"("id");
    END IF;
END $$;
