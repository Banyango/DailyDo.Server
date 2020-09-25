package stores

import (
	"fmt"
	. "github.com/Banyango/gifoody_server/api/model"
	. "github.com/Banyango/gifoody_server/api/repositories/util"
	"github.com/jmoiron/sqlx"
)

type TaskSQLStore struct {
	db *sqlx.DB
}

func NewTaskSQLStore(session *sqlx.DB) *TaskSQLStore {
	self := new(TaskSQLStore)

	self.db = session

	return self
}

func (self *TaskSQLStore) GetTaskAsync(query TaskQuery) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		var results []Task
		rows, err := self.db.Query(`SELECT p.Id, p.taskID ,p.task, p.completed, p.user_id, p.order FROM tasks p LIMIT ? OFFSET ?`, query.Limit, query.Offset)
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}
		defer rows.Close()

		for rows.Next() {
			task := Task{}
			err := rows.Scan(&task.ID, &task.TaskID, &task.Text, &task.Completed, &task.UserID)
			if err != nil {
				storeChan <- StoreResult{Data: nil, Err: err}
				return
			}
			results = append(results, task)
		}

		var count int
		row := self.db.QueryRow("SELECT COUNT(*) FROM tasks")
		err = row.Scan(&count)
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}

		storeChan <- StoreResult{Data: results, Total: count, Err: nil}
	}()
	return storeChan
}

func (self *TaskSQLStore) GetTaskByIdAsync(id string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		user := Task{}
		err := self.db.Get(&user, "SELECT * from tasks WHERE id = ?", id)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *TaskSQLStore) GetChildrenByTaskIdAsync(id string, limit int, offset int) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		var results []Task
		rows, err := self.db.Query(`SELECT s.Id, s.task_id, s.text, s.completed, s.order  FROM tasks s WHERE s.task_id = ? AND (s.type = 'subtask' or s.type = 'summary') LIMIT ? OFFSET ?`, id,
			limit, offset);
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}
		defer rows.Close()

		for rows.Next() {
			task := Task{}
			err := rows.Scan(&task.ID, &task.TaskID, &task.Text, &task.Completed, &task.Order)
			if err != nil {
				storeChan <- StoreResult{Data: nil, Err: err}
				return
			}
			results = append(results, task)
		}

		var count int
		row := self.db.QueryRow("SELECT COUNT(*) FROM tasks")
		err = row.Scan(&count)
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}

		storeChan <- StoreResult{Data: results, Total: count, Err: nil}
	}()
	return storeChan
}

func (self *TaskSQLStore) Save(task Task) StoreResult {
	_, err := self.db.Exec("INSERT INTO tasks " +
		"(id, task_id, text, completed, task_order, user_id) " +
		"values (?, ?, ?, ?, ?, ?)",
		task.ID, task.TaskID, task.Text, task.Completed, task.Completed, task.UserID)
	return StoreResult{
		Data:  task,
		Total: 1,
		Err:   err,
	}
}

func (self *TaskSQLStore) UpdateAsync(task Task) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		_, err := self.db.NamedExec("UPDATE tasks SET `text`=:text, task_order=:task_order, completed=:completed WHERE id=:id", &task)
		storeChan <- StoreResult{
			Data:  task,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *TaskSQLStore) Delete(id string) StoreResult {
	if id == "" {
		return StoreResult{
			Data:  nil,
			Total: 0,
			Err:   fmt.Errorf("null id"),
		}
	}

	_, err := self.db.Exec("DELETE from tasks where id = ?", id)

	return StoreResult{
		Data:  nil,
		Total: 0,
		Err:   err,
	}
}