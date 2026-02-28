create table if not exists users (
 id    serial       primary key,
 name  varchar(255) not null,
 email varchar(255) not null default '',
 age   int          not null default 0,
 phone varchar(50)  not null default ''
);
insert into users (name, email, age, phone) values ('John Doe', 'john@example.com', 25, '+1234567890');
