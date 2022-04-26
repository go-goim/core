-- create database
create database if not exists goim;

-- define user table based on go structure User in current directory
DROP TABLE IF EXISTS goim.user;

create table if not exists goim.user (
	`id` bigint not null auto_increment,
	`uid` varchar(64) not null, -- 22 bytes of uuid
	`name` varchar(32) not null,
	`password` varchar(32) not null,
	`email` varchar(32) not null,
	`phone` varchar(32) not null,
	`avatar` varchar(128) not null,
	`status` tinyint not null,
	`create_at` int not null,
	`update_at` int not null,
	primary key (`id`),
	unique key (`uid`),
    key (`email`),
    key (`phone`)
) auto_increment = 10000 engine = innodb charset = utf8mb4;

-- mock data
insert into goim.user (`id`, `uid`, `name`, `password`, `email`, `phone`, `avatar`, `status`, `create_at`, `update_at`)
values
    (10000, '4F8DSQByUsEUMoETzTCabh', 'user1', 'e10adc3949ba59abbe56e057f20f883e', 'user0@example.com', ' ', ' ', 1, 1528894200, 1528894200),
    (10001, 'C6CtUjpC6h5e5SW9tBFNVX', 'user2', 'e10adc3949ba59abbe56e057f20f883e', 'user1@example.com', ' ', ' ', 0, 1528894200, 1528894200),
    (10002, '7mRZLYedtK1EwxzC5X1Lxf', 'user3', 'e10adc3949ba59abbe56e057f20f883e', 'user2@example.com', ' ', ' ', 0, 1528894200, 1528894200),
    (10003, 'WmbtshDDMUgb3KWFisWZ4E', 'user4', 'e10adc3949ba59abbe56e057f20f883e', 'user3@example.com', ' ', ' ', 0, 1528894200, 1528894200);