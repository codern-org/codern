CREATE TABLE IF NOT EXISTS `user` (
  `id` VARCHAR(64) PRIMARY KEY,
  `email` VARCHAR(64) NOT NULL,
  `display_name` VARCHAR(64) NOT NULL,
  `profile_url` VARCHAR(128) NOT NULL,
  `provider` ENUM('SELF', 'GOOGLE') NOT NULL,
  `created_at` INTEGER UNSIGNED NOT NULL
);
CREATE TABLE IF NOT EXISTS `session` (
  `id` VARCHAR(64) PRIMARY KEY,
  `user_id` VARCHAR(64) NOT NULL,
  `ip_address` VARCHAR(15) NOT NULL,
  `user_agent` VARCHAR(128) NOT NULL,
  `expiry_at` INTEGER UNSIGNED NOT NULL,
  `created_at` INTEGER UNSIGNED NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
);
CREATE TABLE IF NOT EXISTS `organization` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `display_name` VARCHAR(64) NOT NULL,
  `owner_id` VARCHAR(191) NOT NULL
);