create table user (
    `id` bigint(20) not null auto_increment,
    `user_id` bigint(20) not null default 0,
    `username` varchar(64) not null default '',
    `password` varchar(64) not null default '',
    `email` varchar(64) not null default '',
    `gender` tinyint(4) not null default 0,
    `create_time` bigint(20) not null default 0,
    `update_time` bigint(20) not null default 0,
    primary key (`id`),
    unique key `idx_username` (`username`) using btree,
    unique key `idx_user_id` (`user_id`) using btree
)engine=innodb;

CREATE table `community` (
	`id` int(11) not null AUTO_INCREMENT,
	`community_id` int(10) UNSIGNED not null,
	`community_name` VARCHAR(128) not null default '',
	`introduction` VARCHAR(256) not null default '', 
    `create_time` bigint(20) not null default 0,
    `update_time` bigint(20) not null default 0,
    primary key (`id`),
    unique key `idx_community_id` (`community_id`),
    unique key `idx_community_name` (`community_name`)
)engine=innodb;

create table post (
	`id` bigint(20) not null auto_increment,
	`post_id` bigint(20) not null default 0 comment '贴子id',
	`title` varchar(128) not null default '' comment '标题',
	`content` varchar(1024) not null default '' comment '内容',
	`author_id` bigint(20) not null DEFAULT 0 comment '作者',
	`community_id` bigint(20) not null default 0 comment '所属社区',
	`status` tinyint(4) not null DEFAULT 0 comment '贴子状态',
    `create_time` bigint(20) NOT NULL DEFAULT '0',
    `update_time` bigint(20) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_id` (`post_id`),
    KEY `idx_author_id` (`author_id`),
	KEY `idx_community_id` (`community_id`)
) ENGINE=InnoDB;