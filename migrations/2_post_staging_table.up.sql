create table post_staging(
    id char(36) not null,
    name varchar(512) NOT NULL,
    url varchar(512) NOT NULL,
    user_id char(36) NOT NULL,
    post_date date,
    imported bool NOT NULL,
    Quality int,
    PRIMARY KEY (id)
);

create table user(
    id char(36) not null,
    first_name varchar(100),
    last_name varchar(100),
    email varchar(512) NOT NULL UNIQUE,
    username varchar(256) NOT NULL UNIQUE,
    password varchar(512) NOT NULL,
    confirm_token varchar(512),
    verified bool,
    reset bool,
    PRIMARY KEY (id)
);

create table user_forgot_password (
    id char(36) not null,
    token varchar(512) not null,
    created date,
    foreign Key (id) references user(id)
);

# Add user id foreign key to posts.
alter table posts add FOREIGN KEY (user_id) references user(id);

# Forgot primary key on posts.
alter table posts add PRIMARY KEY (id);