alter table tasks
  drop FOREIGN KEY FK_TaskSubTasks;

alter table tasks
  add constraint FOREIGN KEY (task_id) references tasks (id);
