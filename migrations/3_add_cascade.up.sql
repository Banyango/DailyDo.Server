alter table tasks
  drop FOREIGN KEY tasks_ibfk_2;

alter table tasks
  add constraint FK_TaskSubTasks FOREIGN KEY (task_id) references tasks (id) ON DELETE CASCADE;

alter table days
  drop FOREIGN KEY days_ibfk_2;

alter table days
  add constraint FK_DaysParentTask FOREIGN KEY (parent_task_id) references tasks (id) ON DELETE CASCADE;