DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `homes`;
DROP TABLE IF EXISTS `locations`;

CREATE TABLE IF NOT EXISTS `users` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` VARCHAR(64) NULL,
    `gender` INTEGER NULL,
    `age` INTEGER NULL,
    `phone` VARCHAR(30) NULL,
    `created_at` DATE NULL,
    `updated_at` DATE NULL
);

CREATE TABLE IF NOT EXISTS `homes` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `host_id` INTEGER,
    `location_id`INTEGER,
    `created_at` DATE NULL,
    FOREIGN KEY(host_id) REFERENCES user(id),
    FOREIGN KEY(location_id) REFERENCES locations(id)
);

CREATE TABLE IF NOT EXISTS `locations` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `address` VARCHAR(1000) NULL,
    `created_at` DATE NULL
);