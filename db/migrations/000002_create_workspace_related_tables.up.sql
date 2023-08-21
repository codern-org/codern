CREATE TABLE IF NOT EXISTS `workspace` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,
  `profile_url` VARCHAR(128) NOT NULL,
  `owner_id` VARCHAR(64) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`owner_id`) REFERENCES `user`(`id`)
);
CREATE TABLE IF NOT EXISTS `workspace_participant` (
  `workspace_id` INTEGER NOT NULL,
  `user_id` VARCHAR(64) NOT NULL,
  `joined_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`workspace_id`, `user_id`),
  FOREIGN KEY (`workspace_id`) REFERENCES `workspace`(`id`),
  FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
);
