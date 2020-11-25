package stores

import (
	"context"
	"fmt"
	. "github.com/Banyango/gifoody_server/api/model"
	. "github.com/Banyango/gifoody_server/api/repositories/util"
	"github.com/jmoiron/sqlx"
)

type TaskSQLStore struct {
	SqlStore
}

func NewTaskSQLStore(session *sqlx.DB) *TaskSQLStore {
	self := new(TaskSQLStore)

	self.Db = session

	return self
}

func (self *TaskSQLStore) GetTaskAsync(query TaskQuery, ctx context.Context) StoreChannel {
	tx := ctx.Value(TransactionContextKey).(*sqlx.Tx)
	var storeChan = make(StoreChannel, 1)
	go func() {
		var results []Task
		rows, err := tx.Query(`SELECT p.id, p.task_id ,p.text, p.completed, p.user_id, p.task_order FROM tasks p LIMIT ? OFFSET ?`, query.Limit, query.Offset)
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
		row := self.Db.QueryRow("SELECT COUNT(*) FROM tasks")
		err = row.Scan(&count)
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}

		storeChan <- StoreResult{Data: results, Total: count, Err: nil}
	}()
	return storeChan
}

func (self *TaskSQLStore) GetTaskByIdAsync(id string, ctx context.Context) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	tx := ctx.Value(TransactionContextKey).(*sqlx.Tx)
	go func() {
		user := Task{}
		err := tx.Get(&user, "SELECT * from tasks t WHERE t.id = ?", id)
		storeChan <- StoreResult{
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *TaskSQLStore) GetTaskByOrderIdAsync(parentId string, orderId string, ctx context.Context) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	tx := ctx.Value(TransactionContextKey).(*sqlx.Tx)
	go func() {
		user := Task{}
		err := tx.Get(&user, "SELECT * from tasks t WHERE t.task_id = ? and t.task_order = ?", parentId, orderId)
		storeChan <- StoreResult {
			Data:  user,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *TaskSQLStore) GetTasksByParentAsync(id string, ctx context.Context) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	tx := ctx.Value(TransactionContextKey).(*sqlx.Tx)
	go func() {
		var results []Task
		rows, err := tx.Query(`WITH RECURSIVE s AS
                   ( SELECT * FROM tasks
                     WHERE task_order = ?
                     UNION
                     SELECT f.*
                     FROM tasks AS f, s AS a
                     WHERE f.task_order = a.id and f.task_id = ?)
					 SELECT s.id, s.task_id, s.text, s.completed, s.task_order FROM s where s.discriminator = 'Task'`, id, id )
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}
		defer rows.Close()

		for rows.Next() {
			task := Task{
				Type:"Task",
			}
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

func (self *TaskSQLStore) GetChildrenByTaskIdAsync(id string, ctx context.Context) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	tx := ctx.Value(TransactionContextKey).(*sqlx.Tx)
	go func() {
		var results []Task
		rows, err := tx.Query(`WITH RECURSIVE s AS
                   ( SELECT * FROM tasks
                     WHERE task_order = ?
                     UNION
                     SELECT f.*
                     FROM tasks AS f, s AS a
                     WHERE f.task_order = a.id and f.task_id = ?)
					 SELECT s.id, s.discriminator, s.task_id, s.text, s.completed, s.task_order FROM s where s.discriminator in ('SubTask', 'Summary')`, id, id)
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

func (self *TaskSQLStore) Save(task Task, ctx context.Context) StoreResult {
	tx := ctx.Value(TransactionContextKey).(*sqlx.Tx)
	_, err := tx.Exec("INSERT INTO tasks "+
		"(id, discriminator, task_id, text, completed, task_order, user_id) "+
		"values (?, ?, ?, ?, ?, ?, ?)",
		task.ID, task.Type, task.TaskID, task.Text, task.Completed, task.Order, task.UserID)
	return StoreResult{
		Data:  task,
		Total: 1,
		Err:   err,
	}
}

func (self *TaskSQLStore) UpdateAsync(task Task, ctx context.Context) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	tx := ctx.Value(TransactionContextKey).(*sqlx.Tx)
	go func() {
		_, err := tx.NamedExec("UPDATE tasks SET `text`=:text, task_order=:task_order, completed=:completed WHERE id=:id", &task)
		storeChan <- StoreResult{
			Data:  task,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *TaskSQLStore) Delete(id string, ctx context.Context) StoreResult {
	if id == "" {
		return StoreResult{
			Data:  nil,
			Total: 0,
			Err:   fmt.Errorf("null id"),
		}
	}

	tx := ctx.Value(TransactionContextKey).(*sqlx.Tx)
	_, err := tx.Exec("DELETE from tasks where id = ?", id)

	return StoreResult{
		Data:  nil,
		Total: 0,
		Err:   err,
	}
}