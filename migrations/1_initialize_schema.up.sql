create table tasks
(
  id         char(36)  NOT NULL primary key,
  user_id    char(36)  NOT NULL,
  task_id    char(36),
  text       varchar(512),
  completed  bool      NOT NULL,
  task_order integer
);

create table users
(
  id            char(36)     not null primary key,
  first_name    varchar(100),
  last_name     varchar(100),
  email         varchar(512) NOT NULL UNIQUE,
  username      varchar(256) NOT NULL UNIQUE,
  password      varchar(512) NOT NULL,
  confirm_token varchar(512),
  verified      bool,
  reset         bool
);

create table users_forgot_password
(
  id      char(36)     not null primary key,
  token   varchar(512) not null,
  created date,
  foreign Key (id) references users (id)
);

# Add user id foreign key to tasks.
alter table tasks
  add FOREIGN KEY (user_id) references users (id);
alter table tasks
  add FOREIGN KEY (task_id) references tasks (id);

