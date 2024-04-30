ALTER TABLE `assignment`
ADD `publish_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `due_date`;

UPDATE `assignment` SET publish_date = created_at;
