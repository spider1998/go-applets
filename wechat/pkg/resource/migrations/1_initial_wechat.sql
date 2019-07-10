-- +migrate Up
CREATE TABLE `member`
(
  `id`               VARCHAR(64)  NOT NULL,
  `person_id`        VARCHAR(64)  NOT NULL COMMENT '人员 ID',
  `mobile`           VARCHAR(32)  NOT NULL COMMENT '手机号',
  `open_id`          VARCHAR(32)  NOT NULL COMMENT 'OpenID',
  `state`            TINYINT      NOT NULL COMMENT '状态',
  `create_time`      DATETIME     NOT NULL,
  `update_time`      DATETIME     NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`person_id`)
)
  COLLATE = 'utf8mb4_general_ci'
  ENGINE = InnoDB COMMENT '小程序人员';

-- +migrate Down
DROP TABLE `member`;
