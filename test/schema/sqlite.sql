CREATE TABLE `user` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` VARCHAR(64) NULL,
    `gender` INTEGER NULL,
    `age` INTEGER NULL,
    `phone` VARCHAR(30) NULL,
    `created` DATE NULL
);

CREATE TABLE `location` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `host_id` INTEGER
    `loc` VARCHAR(64) NULL,
    `lat` FLOAT NULL,
    `lon` FLOAT NULL,
    `created` DATE NULL
    FOREIGN KEY(host_id) REFERENCES user(id)
);