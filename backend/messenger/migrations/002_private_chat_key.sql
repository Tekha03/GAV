ALTER TABLE "chats"
ADD COLUMN IF NOT EXISTS "private_key" text;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_chats_private_key"
ON "chats" ("private_key")
WHERE "private_key" IS NOT NULL;
