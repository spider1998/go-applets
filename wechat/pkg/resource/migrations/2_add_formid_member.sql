-- +migrate Up
ALTER TABLE `member` ADD `form_id` VARCHAR(64) NOT NULL COMMENT 'form_id' AFTER `person_id`;

-- +migrate Down
ALTER TABLE `member` DROP `form_id`;