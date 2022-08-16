-- User 用户表
CREATE TABLE IF NOT EXISTS `user` (
    `id` bigint unsigned PRIMARY KEY  NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `username` varchar(64) UNIQUE KEY NOT NULL DEFAULT '' COMMENT '用户的名称',
    `status` int NOT NULL DEFAULT 1 COMMENT '用户状态，1=正常，2=禁用',
    `email` varchar(255) UNIQUE KEY NOT NULL DEFAULT '' COMMENT '邮箱',
    `password` varchar(255) NOT NULL DEFAULT '' COMMENT '加密后的密码',
    `description` text NOT NULL COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX(`status`),
    INDEX(`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户表';