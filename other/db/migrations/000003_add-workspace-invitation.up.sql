CREATE TABLE IF NOT EXISTS `workspace_invitation` (
  `id` VARCHAR(128) PRIMARY KEY,
  `workspace_id` BIGINT UNSIGNED NOT NULL,
  `inviter_id` VARCHAR(64) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `valid_at` DATETIME NOT NULL,
  `valid_until` DATETIME NOT NULL,

  FOREIGN KEY (`workspace_id`) REFERENCES `workspace`(`id`),
  FOREIGN KEY (`inviter_id`) REFERENCES `user`(`id`)
)
