ALTER TABLE `workspace`
ADD `is_open_scoreboard` TINYINT(1) NOT NULL DEFAULT '0' AFTER `profile_url`;
