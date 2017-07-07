-- ---
-- Globals
-- ---

-- SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
-- SET FOREIGN_KEY_CHECKS=0;

-- ---
-- Table 'config'
-- 流程定义表
-- ---

DROP TABLE IF EXISTS `config`;
		
CREATE TABLE `config` (
  `config_repo_id` INTEGER NULL DEFAULT NULL,
  `config_data` BLOB NULL DEFAULT NULL COMMENT 'yaml 格式 pipeline 定义的文件',
  `config_hash` VARCHAR(250) NULL DEFAULT NULL COMMENT 'yaml 文件的 hash 值',
  `config_id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`config_id`)
) COMMENT '流程定义表';

-- ---
-- Table 'repos'
-- 代码仓库相关信息表
-- ---

DROP TABLE IF EXISTS `repos`;
		
CREATE TABLE `repos` (
  `repo_scm_id` INTEGER NULL DEFAULT NULL COMMENT '代码仓库类型',
  `repo_clone` VARCHAR NULL DEFAULT NULL COMMENT '代码仓库下载地址',
  `repo_branch` VARCHAR NULL DEFAULT NULL,
  `repo_owner` VARCHAR NULL DEFAULT NULL,
  `repo_name` VARCHAR NULL DEFAULT NULL,
  `repo_allow_push` VARCHAR NULL DEFAULT NULL COMMENT '是否允许 push 触发构建',
  `repo_id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`repo_id`)
) COMMENT '代码仓库相关信息表';

-- ---
-- Table 'sercets'
-- 默认拷贝于 scm_account 保存代码仓库，镜像仓库等的口令信息。
-- ---

DROP TABLE IF EXISTS `sercets`;
		
CREATE TABLE `sercets` (
  `secret_repo_id` INTEGER NULL DEFAULT NULL,
  `secret_name` VARCHAR NULL DEFAULT NULL,
  `secret_value` BLOB NULL DEFAULT NULL,
  `secret_is_default` TINYINT NULL DEFAULT 对于 repo 是否是默认 secret,
  `secret_id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`secret_id`)
) COMMENT '默认拷贝于 scm_account 保存代码仓库，镜像仓库等的口令信息。';

-- ---
-- Table 'builds'
-- 代表一次构建相关信息。
-- ---

DROP TABLE IF EXISTS `builds`;
		
CREATE TABLE `builds` (
  `build_repo_id` INTEGER NULL DEFAULT NULL,
  `build_config_id` INTEGER NULL DEFAULT NULL,
  `build_created` DATETIME NULL DEFAULT NULL,
  `build_started` DATETIME NULL DEFAULT NULL,
  `build_finished` DATETIME NULL DEFAULT NULL,
  `build_number` INTEGER NULL DEFAULT NULL,
  `build_id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`build_id`)
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
-- Table 'scm_account'
-- 
-- ---

DROP TABLE IF EXISTS `scm_account`;
		
CREATE TABLE `scm_account` (
  `scm_host` VARCHAR NULL DEFAULT NULL,
  `scm_login` VARCHAR NULL DEFAULT NULL,
  `scm_password` VARCHAR NULL DEFAULT NULL,
  `scm_type` VARCHAR NULL DEFAULT NULL,
  `scm_id` INTEGER NULL AUTO_INCREMENT DEFAULT NULL,
  PRIMARY KEY (`scm_id`)
);

-- ---
-- Foreign Keys 
-- ---

ALTER TABLE `config` ADD FOREIGN KEY (config_repo_id) REFERENCES `repos` (`repo_id`);
ALTER TABLE `repos` ADD FOREIGN KEY (repo_scm_id) REFERENCES `scm_account` (`scm_id`);
ALTER TABLE `sercets` ADD FOREIGN KEY (secret_repo_id) REFERENCES `repos` (`repo_id`);
ALTER TABLE `builds` ADD FOREIGN KEY (build_repo_id) REFERENCES `repos` (`repo_id`);
ALTER TABLE `builds` ADD FOREIGN KEY (build_config_id) REFERENCES `config` (`config_id`);

-- ---
-- Table Properties
-- ---

-- ALTER TABLE `config` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `repos` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `sercets` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `builds` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `Task` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `File` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `Logs` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
-- ALTER TABLE `scm_account` ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ---
-- Test Data
-- ---

-- INSERT INTO `config` (`config_repo_id`,`config_data`,`config_hash`,`config_id`) VALUES
-- ('','','','');
-- INSERT INTO `repos` (`repo_scm_id`,`repo_clone`,`repo_branch`,`repo_owner`,`repo_name`,`repo_allow_push`,`repo_id`) VALUES
-- ('','','','','','','');
-- INSERT INTO `sercets` (`secret_repo_id`,`secret_name`,`secret_value`,`secret_is_default`,`secret_id`) VALUES
-- ('','','','','');
-- INSERT INTO `builds` (`build_repo_id`,`build_config_id`,`build_created`,`build_started`,`build_finished`,`build_number`,`build_id`) VALUES
-- ('','','','','','','');
-- INSERT INTO `Task` (`id`) VALUES
-- ('');
-- INSERT INTO `File` (`id`) VALUES
-- ('');
-- INSERT INTO `Logs` (`id`) VALUES
-- ('');
-- INSERT INTO `scm_account` (`scm_host`,`scm_login`,`scm_password`,`scm_type`,`scm_id`) VALUES
-- ('','','','','');