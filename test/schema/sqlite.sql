DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `homes`;
DROP TABLE IF EXISTS `locations`;
DROP TABLE IF EXISTS `employees`;
DROP TABLE IF EXISTS `projects`;
DROP TABLE IF EXISTS `tasks`;
DROP TABLE IF EXISTS `domains`;
DROP TABLE IF EXISTS `specialties`;
DROP TABLE IF EXISTS `employees_projects`;


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

CREATE TABLE IF NOT EXISTS `employees` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(64) NULL,
    `gender` INTEGER NULL,
    `age` INTEGER NULL,
    `phone` VARCHAR(30) NULL,
    `salary` REAL NULL,
    `created_at` DATE NULL,
    `updated_at` DATE NULL
);

CREATE TABLE IF NOT EXISTS `projects` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(64) NULL,
    `deadline` DATE NULL
);

CREATE TABLE IF NOT EXISTS `employees_projects` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `project_id` INTEGER,
    `employee_id` INTEGER,
    FOREIGN KEY(project_id) REFERENCES projects(id),
    FOREIGN KEY(employee_id) REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS `tasks` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(64) NULL,
    `project_id` INTEGER,
    `deadline` DATE NULL,
    FOREIGN KEY(project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS `domains` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(64) NULL
);

CREATE TABLE IF NOT EXISTS `specialties` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(64) NULL,
    `owner_id` INTEGER,
    `domain_id` INTEGER,
    FOREIGN KEY(owner_id) REFERENCES employees(id),
    FOREIGN KEY(domain_id) REFERENCES domains(id)
);
