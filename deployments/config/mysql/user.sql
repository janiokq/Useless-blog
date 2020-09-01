

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `id` int(11) NOT NULL AUTO_INCREMENT,
                        `user_name` varchar(30) DEFAULT '' COMMENT '用户名',
                        `phone` char(20) DEFAULT '' COMMENT '手机号',
                        `password` char(32) DEFAULT '' COMMENT '密码',
                        `avatar_url` char(32) DEFAULT '' COMMENT '头像',
                        `create_at` datetime DEFAULT NULL COMMENT '创建时间',
                        `update_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='用户表';

SET FOREIGN_KEY_CHECKS = 1;
