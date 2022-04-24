-- create database
create database if not exists goim;

-- define user table based on go structure User in current directory
create table if not exists goim.user (
	`id` bigint not null auto_increment,
	`uid` varchar(64) not null,
	`name` varchar(32) not null,
	`password` varchar(32) not null,
	`email` varchar(32) not null,
	`phone` varchar(32) not null,
	`avatar` varchar(128) not null,
	`status` tinyint not null,
	`create_at` int not null,
	`update_at` int not null,
	primary key (`id`),
	unique key (`uid`)
) auto_increment = 10000 engine = innodb charset = utf8mb4;