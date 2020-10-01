package stores

import (
	"context"
	"fmt"
	. "github.com/Banyango/gifoody_server/api/model"
	. "github.com/Banyango/gifoody_server/api/repositories/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DaySQLStore struct {
	db *sqlx.DB
}

func NewDaySQLStore(session *sqlx.DB) *DaySQLStore {
	self := new(DaySQLStore)

	self.db = session

	return self
}

func (self *DaySQLStore) GetDaysAsync(userId string, limit int, offset int) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		var results []Day
		rows, err := self.db.Query(`SELECT d.Id, d.date ,d.user_id, d.parent_task_id, d.summary FROM days d WHERE d.user_id = ? ORDER BY d.date DESC LIMIT ? OFFSET ?`, userId, limit, offset)
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}
		defer rows.Close()

		for rows.Next() {
			day := Day{}
			err := rows.Scan(&day.ID, &day.Date, &day.UserID, &day.ParentTaskID, &day.Summary)
			if err != nil {
				storeChan <- StoreResult{Data: nil, Err: err}
				return
			}
			results = append(results, day)
		}

		var count int
		row := self.db.QueryRow("SELECT COUNT(*) FROM days")
		err = row.Scan(&count)
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}

		storeChan <- StoreResult{Data: results, Total: count, Err: nil}
	}()
	return storeChan
}

func (self *DaySQLStore) GetDayByIdAsync(id string) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		day := Day{}
		err := self.db.Get(&day, "SELECT * from days WHERE id = ?", id)
		storeChan <- StoreResult{
			Data:  day,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *DaySQLStore) Save(day Day) StoreResult {

	ctx := context.Background()
	tx, err := self.db.BeginTx(ctx, nil)
	if err != nil {
		return StoreResult{Err: err}
	}

	parentTaskId := uuid.New().String()
	day.ParentTaskID = parentTaskId

	_, err = self.db.ExecContext(ctx, "INSERT INTO tasks "+
		"(id, discriminator, user_id) values (?, ?, ?)",
		parentTaskId, "DayParent", day.UserID)
	if err != nil {
		_ = tx.Rollback()
		return StoreResult{Err: err}
	}

	_, err = self.db.ExecContext(ctx, "INSERT INTO days "+
		"(id, `date`, user_id, parent_task_id, summary) "+
		"values (?, ?, ?, ?, ?)",
		day.ID, day.Date, day.UserID, day.ParentTaskID, day.Summary)
	if err != nil {
		_ = tx.Rollback()
		return StoreResult{Err: err}
	}

	err = tx.Commit()
	if err != nil {
		return StoreResult{Err: err}
	}

	return StoreResult{
		Data:  day,
		Total: 1,
		Err:   err,
	}
}

func (self *DaySQLStore) UpdateAsync(day Day) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		_, err := self.db.NamedExec("UPDATE days SET summary=:summary WHERE id=:id", &day)
		storeChan <- StoreResult{
			Data:  day,
			Total: 1,
			Err:   err,
		}
	}()
	return storeChan
}

func (self *DaySQLStore) Delete(id string) StoreResult {
	if id == "" {
		return StoreResult{
			Data:  nil,
			Total: 0,
			Err:   fmt.Errorf("null id"),
		}
	}

	_, err := self.db.Exec("DELETE from days where id = ?", id)

	return StoreResult{
		Data:  nil,
		Total: 0,
		Err:   err,
	}
}
