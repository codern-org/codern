ALTER TABLE `assignment`
ADD `is_auto_trim_enabled` BOOLEAN NOT NULL DEFAULT false AFTER `updated_at`;
