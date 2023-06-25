CREATE TABLE `movies` (
    `id` bigint NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `title` char(150) NOT NULL,
    `description` text NOT NULL,
    `duration` varchar(250) NOT NULL,
    `artist` char(50) NOT NULL,
    `genre` char(50) NOT NULL,
    `video_url` text NOT NULL,
    `created_at` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL
);