ALTER TABLE `workspace`
ADD `is_archived` TINYINT(1) NOT NULL DEFAULT '0' AFTER `is_deleted`;
