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
		rows, err := self.db.Query(`SELECT p.id, p.task_id ,p.text, p.completed, p.user_id, p.task_order FROM tasks p LIMIT ? OFFSET ?`, query.Limit, query.Offset)
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
		err := self.db.Get(&user, "SELECT * from tasks t WHERE t.id = ?", id)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *TaskSQLStore) GetTasksByParentAsync(id string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		var results []Task
		rows, err := self.db.Query(`SELECT s.id, s.task_id, s.text, s.completed, s.task_order FROM tasks s WHERE s.task_id = ? AND s.discriminator = 'Task'`, id)
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

		storeChan <- StoreResult{Data: results, Total: len(results), Err: nil}
	}()
	return storeChan
}

func (self *TaskSQLStore) GetChildrenByTaskIdAsync(id string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		var results []Task
		rows, err := self.db.Query(`SELECT s.id, s.discriminator, s.task_id, s.text, s.completed, s.task_order FROM tasks s WHERE s.task_id = ? AND ( s.discriminator = 'SubTask' or s.discriminator = 'Summary')`, id, )
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}
		defer rows.Close()

		for rows.Next() {
			task := Task{}
			err := rows.Scan(&task.ID, &task.Type, &task.TaskID, &task.Text, &task.Completed, &task.Order)
			if err != nil {
				storeChan <- StoreResult{Data: nil, Err: err}
				return
			}
			results = append(results, task)
		}

		storeChan <- StoreResult{Data: results, Total: len(results), Err: nil}
	}()
	return storeChan
}

func (self *TaskSQLStore) Save(task Task) StoreResult {
	_, err := self.db.Exec("INSERT INTO tasks "+
		"(id, discriminator, task_id, text, completed, task_order, user_id) "+
		"values (?, ?, ?, ?, ?, ?, ?)",
		task.ID, task.Type, task.TaskID, task.Text, task.Completed, task.Completed, task.UserID)
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
