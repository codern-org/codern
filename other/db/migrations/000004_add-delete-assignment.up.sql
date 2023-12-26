-- add is_deleted column to assignment table

ALTER TABLE `assignment`

ADD `is_deleted` TINYINT(1) NOT NULL DEFAULT 0;
