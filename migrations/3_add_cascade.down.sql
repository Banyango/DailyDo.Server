alter table tasks
  drop FOREIGN KEY FK_TaskSubTasks;

alter table tasks
  add constraint FOREIGN KEY (task_id) references tasks (id);

alter table days
  drop FOREIGN KEY FK_DaysParentTask;

alter table days
  add constraint FOREIGN KEY (parent_task_id) references tasks (id);