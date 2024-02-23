ALTER TABLE `workspace`
ADD `is_deleted` TINYINT(1) NOT NULL DEFAULT '0' AFTER `is_open_scoreboard`;
