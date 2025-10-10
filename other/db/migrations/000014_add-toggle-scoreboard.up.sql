ALTER TABLE `workspace`
ADD `is_scoreboard_enabled` TINYINT(1) NOT NULL DEFAULT '0' AFTER `profile_url`;
