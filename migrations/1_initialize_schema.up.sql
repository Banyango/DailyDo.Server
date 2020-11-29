################################################
## Tables
################################################

create table tasks
(
  id         char(36)  not null primary key,
  user_id    char(36)  not null,
  task_id    char(36),
  text       varchar(512),
  day_id     char(36),
  completed  bool      not null,
  task_order integer
);

create table users
(
  id            char(36)     not null primary key,
  first_name    varchar(100),
  last_name     varchar(100),
  email         varchar(512) not null unique,
  username      varchar(256) not null unique,
  password      varchar(512) not null,
  confirm_token varchar(512),
  verified      bool,
  reset         bool
);

create table users_forgot_password
(
  id      char(36)     not null primary key,
  token   varchar(512) not null,
  created date,
  foreign key (id) references users (id)
);

create table days
(
    id             char(36) not null primary key,
    date           DATE     not null,
    summary        varchar(512),
    user_id        char(36) not null,
    parent_task_id char(36) not null
);

################################################
## Constraints
################################################

# Tasks
alter table tasks add column discriminator enum ('DayParent', 'Task','SubTask', 'Summary') not null;
alter table tasks alter completed set default false;
alter table tasks alter task_order set default 0;
alter table tasks
    add constraint FK_TaskUsers foreign key (user_id) references users (id);
alter table tasks
    add constraint FK_TaskSubTasks foreign key (task_id) references tasks (id);
alter table tasks
    add constraint FK_TaskDay foreign key (day_id) references days (id) on delete cascade;

# Days
alter table days
    add constraint FK_DayUser foreign key (user_id) references users (id);
alter table days
    add constraint FK_DaysParentTask foreign key (parent_task_id) references tasks (id);

