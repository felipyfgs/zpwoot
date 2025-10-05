-- Migration: initial_schema (rollback)
-- Drop zpwoot database schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_zp_message_updated_at ON "zpMessage";
DROP TRIGGER IF EXISTS update_zp_chatwoot_updated_at ON "zpChatwoot";
DROP TRIGGER IF EXISTS update_zp_webhooks_updated_at ON "zpWebhooks";
DROP TRIGGER IF EXISTS update_zp_sessions_updated_at ON "zpSessions";

-- Drop tables (in reverse order due to foreign key constraints)
DROP TABLE IF EXISTS "zpMessage";
DROP TABLE IF EXISTS "zpChatwoot";
DROP TABLE IF EXISTS "zpWebhooks";
DROP TABLE IF EXISTS "zpSessions";

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();
