DROP TABLE IF EXISTS `users`;

CREATE TABLE IF NOT EXISTS `users` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(64) NULL,
    `address` INTEGER NULL,
    `gender` INTEGER NULL,
    `phone` VARCHAR(30) NULL,
    `created_at` DATE NULL,
    `updated_at` DATE NULL
);