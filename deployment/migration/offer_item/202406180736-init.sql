-- +migrate Up
CREATE TABLE `offer_item` (
  `id` char(22) NOT NULL,
  `name` varchar(256) NOT NULL,
  `item_id` varchar(32) NOT NULL,
  `df_item_id` varchar(32) DEFAULT NULL,
  `coupon_banner_id` varchar(255) DEFAULT NULL,
  `special_rate` double NOT NULL,
  `special_amount` int(11) NOT NULL,
  `has_sample` tinyint(1) NOT NULL,
  `needs_preliminary_review` tinyint(1) NOT NULL,
  `needs_after_review` tinyint(1) NOT NULL,
  `needs_pr_mark` tinyint(1) NOT NULL DEFAULT '1',
  `post_required` tinyint(1) NOT NULL,
  `post_target` int(10) unsigned NOT NULL,
  `has_coupon` tinyint(1) NOT NULL,
  `has_special_commission` tinyint(1) NOT NULL,
  `has_lottery` tinyint(1) NOT NULL,
  `product_features` text NOT NULL,
  `cautionary_points` text NOT NULL,
  `reference_info` text NOT NULL,
  `other_info` text NOT NULL,
  `is_invitation_mail_sent` tinyint(1) NOT NULL,
  `is_offer_detail_mail_sent` tinyint(1) NOT NULL,
  `is_passed_preliminary_review_mail_sent` tinyint(1) NOT NULL,
  `is_failed_preliminary_review_mail_sent` tinyint(1) NOT NULL,
  `is_article_post_mail_sent` tinyint(1) NOT NULL,
  `is_passed_after_review_mail_sent` tinyint(1) NOT NULL,
  `is_failed_after_review_mail_sent` tinyint(1) NOT NULL,
  `is_closed` tinyint(1) NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(64) NOT NULL,
  `updated_at` datetime NOT NULL,
  `updated_by` varchar(64) NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`),
  KEY `idx_item_id` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `assignee` (
  `id` char(22) NOT NULL,
  `offer_item_id` char(22) NOT NULL,
  `ameba_id` varchar(256) NOT NULL,
  `stage` int(10) unsigned NOT NULL,
  `writing_fee` int(11) NOT NULL,
  `decline_reason` text,
  `created_at` datetime NOT NULL,
  `created_by` varchar(64) NOT NULL,
  `updated_at` datetime NOT NULL,
  `updated_by` varchar(64) NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `offer_item_id` (`offer_item_id`,`ameba_id`),
  KEY `idx_ameba_id` (`ameba_id`),
  KEY `idx_stage` (`stage`),
  CONSTRAINT `assign_ibfk_1` FOREIGN KEY (`offer_item_id`) REFERENCES `offer_item` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `drafted_item_info` (
  `offer_item_id` char(22) NOT NULL,
  `name` text NOT NULL COMMENT '商品情報の名前。',
  `content_name` text NOT NULL COMMENT '会社名。',
  `image_url` text NOT NULL COMMENT '商品画像URL。',
  `url` text NOT NULL COMMENT '商品詳細のURL。',
  `min_commission` double NOT NULL COMMENT '最小の報酬額',
  `min_commission_type` int(11) NOT NULL COMMENT 'commissionのタイプ',
  `max_commission` double NOT NULL COMMENT '最大の報酬額。',
  `max_commission_type` int(11) NOT NULL COMMENT 'commissionのタイプ。',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`offer_item_id`),
  CONSTRAINT `drafted_item_info_ibfk_1` FOREIGN KEY (`offer_item_id`) REFERENCES `offer_item` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `examination` (
  `id` char(23) NOT NULL,
  `offer_item_id` char(22) NOT NULL,
  `assignee_id` char(22) NOT NULL,
  `entry_id` varchar(256) DEFAULT NULL,
  `sns_user_id` varchar(256) DEFAULT NULL,
  `sns_screenshot_url` mediumblob,
  `reason` text,
  `examiner_name` varchar(255) DEFAULT NULL,
  `entry_type` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `offer_item_id` (`offer_item_id`),
  KEY `assignee_id` (`assignee_id`),
  CONSTRAINT `examination_ibfk_1` FOREIGN KEY (`offer_item_id`) REFERENCES `offer_item` (`id`) ON DELETE CASCADE,
  CONSTRAINT `examination_ibfk_2` FOREIGN KEY (`assignee_id`) REFERENCES `assignee` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `schedule` (
  `id` char(22) NOT NULL,
  `offer_item_id` char(22) NOT NULL,
  `schedule_type` int(10) unsigned NOT NULL,
  `start_date` datetime DEFAULT NULL,
  `end_date` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(64) NOT NULL,
  `updated_at` datetime NOT NULL,
  `updated_by` varchar(64) NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `offer_item_id` (`offer_item_id`,`schedule_type`),
  KEY `idx_end_date` (`end_date`),
  CONSTRAINT `schedule_ibfk_1` FOREIGN KEY (`offer_item_id`) REFERENCES `offer_item` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `questionnaire` (
  `offer_item_id` char(22) NOT NULL,
  `description` text NOT NULL,
  PRIMARY KEY (`offer_item_id`),
  CONSTRAINT `questionnaire_ibfk_1` FOREIGN KEY (`offer_item_id`) REFERENCES `offer_item` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `questionnaire_question` (
  `id` char(22) NOT NULL,
  `offer_item_id` char(22) NOT NULL,
  `title` text NOT NULL,
  `type` int(11) NOT NULL,
  `image` mediumtext NOT NULL,
  `answer_options` json DEFAULT NULL,
  `priority` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `offer_item_id` (`offer_item_id`),
  CONSTRAINT `questionnaire_question_ibfk_1` FOREIGN KEY (`offer_item_id`) REFERENCES `offer_item` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `questionnaire_question_answer` (
  `assignee_id` char(22) NOT NULL,
  `questionnaire_question_id` char(22) NOT NULL,
  `offer_item_id` char(22) NOT NULL,
  `answer` text NOT NULL,
  PRIMARY KEY (`assignee_id`,`questionnaire_question_id`),
  KEY `questionnaire_question_answer_ibfk_1` (`offer_item_id`),
  CONSTRAINT `questionnaire_question_answer_ibfk_1` FOREIGN KEY (`offer_item_id`) REFERENCES `offer_item` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE IF EXISTS `questionnaire_question_answer`;
DROP TABLE IF EXISTS  `questionnaire_question`;
DROP TABLE IF EXISTS  `questionnaire`;
DROP TABLE IF EXISTS examination;
DROP TABLE IF EXISTS schedule;
DROP TABLE IF EXISTS examination;
DROP TABLE IF EXISTS assignee;
DROP TABLE IF EXISTS offer_item;
DROP TABLE IF EXISTS  `draft_Item_Info`;
