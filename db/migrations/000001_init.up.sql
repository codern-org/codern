CREATE TABLE IF NOT EXISTS `user` (
  `id` VARCHAR(64) PRIMARY KEY,
  `email` VARCHAR(64) NOT NULL,
  `display_name` VARCHAR(64) NOT NULL,
  `profile_url` VARCHAR(128) NOT NULL,
  `account_type` VARCHAR(32) NOT NULL,
  `provider` VARCHAR(32) NOT NULL,
  `password` VARCHAR(128),
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
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

-- Workspace

CREATE TABLE IF NOT EXISTS `workspace` (
  `id` BIGINT UNSIGNED PRIMARY KEY,
  `name` VARCHAR(64) NOT NULL,
  `profile_url` VARCHAR(128) NOT NULL,
  `owner_id` VARCHAR(64) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`owner_id`) REFERENCES `user`(`id`)
);

CREATE TABLE IF NOT EXISTS `workspace_participant` (
  `workspace_id` BIGINT UNSIGNED NOT NULL,
  `user_id` VARCHAR(64) NOT NULL,
  `joined_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `recently_visited_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`workspace_id`, `user_id`),
  FOREIGN KEY (`workspace_id`) REFERENCES `workspace`(`id`),
  FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
);

CREATE TABLE IF NOT EXISTS `assignment` (
  `id` BIGINT UNSIGNED PRIMARY KEY,
  `workspace_id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  `description` VARCHAR(64) NOT NULL,
  `detail_url` VARCHAR(128) NOT NULL,
  `memory_limit` INTEGER NOT NULL,
  `time_limit` INTEGER NOT NULL,
  `level` VARCHAR(32) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`workspace_id`) REFERENCES `workspace`(`id`)
);

CREATE TABLE IF NOT EXISTS `testcase` (
  `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `assignment_id` BIGINT UNSIGNED NOT NULL,
  `input_file_url` VARCHAR(128) NOT NULL,
  `output_file_url` VARCHAR(128) NOT NULL,
  FOREIGN KEY (`assignment_id`) REFERENCES `assignment`(`id`)
);

CREATE TABLE IF NOT EXISTS `submission` (
  `id` BIGINT UNSIGNED PRIMARY KEY,
  `assignment_id` BIGINT UNSIGNED NOT NULL,
  `user_id` VARCHAR(64) NOT NULL,
  `language` VARCHAR(64) NOT NULL,
  `file_url` VARCHAR(128) NOT NULL,
  `submitted_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `compilation_log` LONGTEXT,
  FOREIGN KEY (`assignment_id`) REFERENCES `assignment`(`id`),
  FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
);

CREATE TABLE IF NOT EXISTS `submission_result` (
  `submission_id` BIGINT UNSIGNED NOT NULL,
  `testcase_id` BIGINT UNSIGNED NOT NULL,
  `status` VARCHAR(32) NOT NULL,
  `status_detail` VARCHAR(32),
  `memory_usage` INTEGER,
  `time_usage` INTEGER,
  PRIMARY KEY (`submission_id`, `testcase_id`),
  FOREIGN KEY (`submission_id`) REFERENCES `submission`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`testcase_id`) REFERENCES `testcase`(`id`) ON DELETE CASCADE
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

INSERT INTO `workspace` VALUES (
  '1',
  'Codern Playground',
  '',
  '62b870d7a68388007ba0f8ba292686c70dcb06b8',
  '2023-08-20 09:30:00'
);

INSERT INTO `workspace_participant` VALUES (
  '1',
  '62b870d7a68388007ba0f8ba292686c70dcb06b8',
  '2023-08-20 09:30:00',
  DEFAULT
);

INSERT INTO `assignment` VALUES (
  '1',
  '1',
  'Keeratikorn Noodle',
  'The hardest algorithm question in the software world',
  '',
  '1024',
  '500',
  'EASY',
  '2023-08-20 09:30:00',
  '2023-08-20 09:30:00'
);

INSERT INTO `assignment` VALUES (
  '2',
  '1',
  'Porama Chicken',
  'The most chicken algorithm question in the software world',
  '',
  '1024',
  '500',
  'EASY',
  '2023-08-20 09:30:00',
  '2023-08-20 09:30:00'
);

INSERT INTO `testcase` VALUES (
  '1',
  '1',
  'file_url_1'
);

INSERT INTO `testcase` VALUES (
  '2',
  '1',
  'file_url_2'
);

INSERT INTO `testcase` VALUES (
  '3',
  '2',
  'file_url_1'
);

INSERT INTO `workspace` VALUES (
  '2',
  'Codern Playground 2',
  '',
  '62b870d7a68388007ba0f8ba292686c70dcb06b8',
  '2023-08-20 09:30:00'
);

INSERT INTO `workspace_participant` VALUES (
  '2',
  '62b870d7a68388007ba0f8ba292686c70dcb06b8',
  '2023-08-20 09:30:00',
  DEFAULT
);
