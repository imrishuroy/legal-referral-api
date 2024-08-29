-- Drop constraints first if they exist
ALTER TABLE IF EXISTS attachments DROP CONSTRAINT IF EXISTS fk_attachment_message;
ALTER TABLE IF exists messages DROP CONSTRAINT IF EXISTS fk_message_sender;
ALTER TABLE IF exists messages DROP CONSTRAINT IF EXISTS fk_message_recipient;
ALTER TABLE IF exists messages DROP CONSTRAINT IF EXISTS unique_room_message;

-- Drop indexes if necessary
DROP INDEX IF EXISTS idx_sender_id;
DROP INDEX IF EXISTS idx_recipient_id;
DROP INDEX IF EXISTS idx_room_id;
DROP INDEX IF EXISTS idx_parent_message_id;
DROP INDEX IF EXISTS idx_attachment_id;
DROP INDEX IF EXISTS idx_message_id;

-- Drop tables
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS messages;
