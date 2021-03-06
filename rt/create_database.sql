
drop database if exists pondgw;
create database pondgw;

use pondgw;

-- --------------------------------------------------------------------
-- EMAIL-related tables
-- --------------------------------------------------------------------

create table email (
    id int not null auto_increment,
    ts timestamp default current_timestamp on update current_timestamp,
    status int default 0,
    addr varchar(255) not null,
    pubkey mediumblob not null,
    token varchar(32) not null,
    primary key(id)
) engine=MyISAM;

create index mail_idx on email(addr);

-- --------------------------------------------------------------------
-- POND-related tables
-- --------------------------------------------------------------------

create table pond (
    id int not null auto_increment,
    ts timestamp default current_timestamp on update current_timestamp,
    status int default 0,
    peer varchar(16) not null,
    primary key(id)
) engine=MyISAM;

create index pond_idx on pond(peer);

-- --------------------------------------------------------------------
-- statistics-related tables
-- --------------------------------------------------------------------

create table stats (
	name varchar(8),
	val int default 0,
    primary key(name)
) engine=MyISAM;

insert into stats(name) values('last');
insert into stats(name) values('daily');
insert into stats(name) values('weekly');
insert into stats(name) values('monthly');
insert into stats(name) values('yearly');

-- --------------------------------------------------------------------
-- user management
-- --------------------------------------------------------------------

-- drop user 'pondgw'@'localhost';
-- create user 'pondgw'@'localhost' identified by 'pondgw';
-- grant select,delete,update,insert on pondgw.* to 'pondgw'@'localhost';
