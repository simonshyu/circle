-- ---
-- Globals
-- ---

-- SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
-- SET FOREIGN_KEY_CHECKS=0;

-- ---
-- Table 'Config'
-- 流程定义表
-- ---

DROP TABLE IF EXISTS `Config`;
		
CREATE TABLE `Config` (
  `id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`id`)
) COMMENT '流程定义表';

-- ---
-- Table 'Repo'
-- 代码仓库相关信息表
-- ---

DROP TABLE IF EXISTS `Repo`;
		
CREATE TABLE `Repo` (
  `id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`id`)
) COMMENT '代码仓库相关信息表';

-- ---
-- Table 'Sercet'
-- 保存代码仓库，镜像仓库等的口令信息。
-- ---

DROP TABLE IF EXISTS `Sercet`;
		
CREATE TABLE `Sercet` (
  `id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`id`)
) COMMENT '保存代码仓库，镜像仓库等的口令信息。';

-- ---
-- Table 'Build'
-- 代表一次构建相关信息。
-- ---

DROP TABLE IF EXISTS `Build`;
		
CREATE TABLE `Build` (
  `id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`id`)
) COMMENT '代表一次构建相关信息。';

-- ---
-- Table 'Task'
-- 代表一次构建调度信息。
-- ---

DROP TABLE IF EXISTS `Task`;
		
CREATE TABLE `Task` (
  `id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`id`)
) COMMENT '代表一次构建调度信息。';

-- ---
-- Table 'File'
-- 代表一次构建产生的artifact(如日志文件)
-- ---

DROP TABLE IF EXISTS `File`;
		
CREATE TABLE `File` (
  `id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`id`)
) COMMENT '代表一次构建产生的artifact(如日志文件)';

-- ---
-- Table 'Logs'
-- 代表一次构建的 Logs
-- ---

DROP TABLE IF EXISTS `Logs`;
		
CREATE TABLE `Logs` (
  `id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`id`)
) COMMENT '代表一次构建的 Logs';

-- ---
-- Foreign Keys 
-- ---


-- ---
-- Table Properties
-- ---

-- ALTER TABLE `Config` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `Repo` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `Sercet` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `Build` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `Task` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `File` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `Logs` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ---
-- Test Data
-- ---

-- INSERT INTO `Config` (`id`) VALUES
-- ('');
-- INSERT INTO `Repo` (`id`) VALUES
-- ('');
-- INSERT INTO `Sercet` (`id`) VALUES
-- ('');
-- INSERT INTO `Build` (`id`) VALUES
-- ('');
-- INSERT INTO `Task` (`id`) VALUES
-- ('');
-- INSERT INTO `File` (`id`) VALUES
-- ('');
-- INSERT INTO `Logs` (`id`) VALUES
-- ('');