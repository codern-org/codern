CREATE TABLE IF NOT EXISTS `user` (
  `id` VARCHAR(64) PRIMARY KEY,
  `email` VARCHAR(64) NOT NULL,
  `display_name` VARCHAR(64) NOT NULL,
  `profile_url` VARCHAR(128) NOT NULL,
  `account_type` ENUM('FREE', 'PRO') NOT NULL DEFAULT 'FREE',
  `provider` ENUM('SELF', 'GOOGLE') NOT NULL,
  `password` VARCHAR(128),
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS `session` (
  `id` VARCHAR(128) PRIMARY KEY,
  `user_id` VARCHAR(64) NOT NULL,
  `ip_address` VARCHAR(15) NOT NULL,
  `user_agent` VARCHAR(256) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `expired_at` DATETIME NOT NULL,
  FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
);
CREATE TABLE IF NOT EXISTS `organization` (
  `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
  `display_name` VARCHAR(64) NOT NULL,
  `owner_id` VARCHAR(191) NOT NULL
);

-- Seeding

INSERT INTO `user` VALUES (
  '62b870d7a68388007ba0f8ba292686c70dcb06b8',
  'admin@codern.app',
  'Codern Admin',
  '',
  'FREE',
  'SELF',
  '$2a$10$9Wz5oM5LlUvFDtgZh3xlzeDtyg1PATgFVg12cafi0cAFDzy8SCUbm',
  '2023-08-20 09:30:00'
);
