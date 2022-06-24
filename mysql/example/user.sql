CREATE DATABASE test CHARSET utf8mb4;

USE test;

CREATE TABLE `user_example` (
    `id`  int(10) UNSIGNED NOT NULL AUTO_INCREMENT ,
    `created_at`  timestamp NULL DEFAULT NULL ,
    `updated_at`  timestamp NULL DEFAULT NULL ,
    `deleted_at`  timestamp NULL DEFAULT NULL ,
    `name`  varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL ,
    `age`  int(11) NOT NULL ,
    `gender`  varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL ,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uix_user_name` (`name`) USING BTREE ,
    INDEX `idx_user_deleted_at` (`deleted_at`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci AUTO_INCREMENT=1 ROW_FORMAT=DYNAMIC;


INSERT INTO `user_example` VALUES (1, '2019-10-14 11:11:51', '2019-10-14 11:11:51', NULL, '刘备', 25, '男');
INSERT INTO `user_example` VALUES (2, '2019-10-14 10:16:03', '2019-10-14 10:16:03', NULL, '关羽', 22, '男');
INSERT INTO `user_example` VALUES (3, '2019-10-14 13:26:32', '2019-10-14 13:26:32', NULL, '张飞', 23, '男');
INSERT INTO `user_example` VALUES (4, '2019-10-14 13:47:18', '2019-10-14 13:47:18', NULL, '诸葛亮', 18, '男');
INSERT INTO `user_example` VALUES (5, '2019-10-14 14:15:44', '2019-10-14 14:15:44', NULL, '赵云', 21, '男');
INSERT INTO `user_example` VALUES (6, '2019-10-14 13:58:41', '2019-10-14 13:58:41', NULL, '曹操', 29, '男');
INSERT INTO `user_example` VALUES (7, '2019-10-14 14:27:36', '2019-10-14 14:27:36', NULL, '司马懿', 23, '男');
INSERT INTO `user_example` VALUES (8, '2019-10-14 14:05:35', '2019-10-14 14:05:35', NULL, '孙权', 27, '男');
INSERT INTO `user_example` VALUES (9, '2019-10-14 14:25:31', '2019-10-14 14:25:31', NULL, '周瑜', 20, '男');
INSERT INTO `user_example` VALUES (10, '2019-10-14 18:10:01', '2019-10-14 18:10:01', NULL, '大乔', 16, '女');
INSERT INTO `user_example` VALUES (11, '2019-10-14 19:40:21', '2019-10-14 19:40:21', NULL, '小乔', 15, '女');
INSERT INTO `user_example` VALUES (12, '2019-10-14 14:30:47', '2019-10-14 14:30:47', NULL, '董卓', 33, '男');
INSERT INTO `user_example` VALUES (13, '2019-10-14 13:50:49', '2019-10-14 13:50:49', NULL, '吕布', 23, '男');
INSERT INTO `user_example` VALUES (14, '2019-10-14 14:33:19', '2019-10-14 14:33:19', NULL, '貂蝉', 16, '女');

