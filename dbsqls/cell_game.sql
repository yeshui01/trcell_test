/*
 Navicat Premium Data Transfer

 Source Server         : wsl_mysql5.7
 Source Server Type    : MySQL
 Source Server Version : 50739
 Source Host           : localhost:3306
 Source Schema         : cell_game

 Target Server Type    : MySQL
 Target Server Version : 50739
 File Encoding         : 65001

 Date: 29/09/2022 10:04:42
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for role_base
-- ----------------------------
DROP TABLE IF EXISTS `role_base`;
CREATE TABLE `role_base`  (
  `role_id` bigint(19) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(19) NOT NULL DEFAULT 0,
  `role_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `create_time` bigint(19) NOT NULL DEFAULT 0,
  `level` int(11) NOT NULL DEFAULT 0,
  `login_time` bigint(19) NOT NULL DEFAULT 0,
  `offline_time` bigint(19) NOT NULL DEFAULT 0,
  `upd_time` bigint(19) NOT NULL DEFAULT 0,
  PRIMARY KEY (`role_id`) USING BTREE,
  UNIQUE INDEX `uname_idx`(`role_name`) USING BTREE COMMENT '昵称'
) ENGINE = InnoDB AUTO_INCREMENT = 1000002 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for role_coin
-- ----------------------------
DROP TABLE IF EXISTS `role_coin`;
CREATE TABLE `role_coin`  (
  `role_id` bigint(19) NOT NULL COMMENT '角色id',
  `coin` json NOT NULL COMMENT '货币数据',
  `upd_time` bigint(19) NOT NULL DEFAULT 0,
  PRIMARY KEY (`role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for role_equip
-- ----------------------------
DROP TABLE IF EXISTS `role_equip`;
CREATE TABLE `role_equip`  (
  `role_id` bigint(20) NOT NULL,
  `equip` mediumblob NOT NULL COMMENT '装备数据',
  `upd_time` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for role_extra
-- ----------------------------
DROP TABLE IF EXISTS `role_extra`;
CREATE TABLE `role_extra`  (
  `role_id` bigint(19) NOT NULL,
  `extra_data` mediumblob NOT NULL COMMENT '角色额外数据',
  `upd_time` bigint(20) NOT NULL,
  PRIMARY KEY (`role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
