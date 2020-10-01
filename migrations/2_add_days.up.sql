create table days
(
  id             char(36) NOT NULL primary key,
  date           DATE     NOT NULL,
  summary        varchar(512),
  user_id        char(36) NOT NULL,
  parent_task_id char(36) NOT NULL
);

alter table days
  ADD FOREIGN KEY (user_id) references users (id);
alter table days
  ADD FOREIGN KEY (parent_task_id) references tasks (id);

### Forgot some tasks stuff.
alter table tasks
  ADD COLUMN discriminator ENUM ('DayParent', 'Task','SubTask', 'Summary') not null;
alter table tasks alter completed set default false;
alter table tasks alter task_order set default 0;

