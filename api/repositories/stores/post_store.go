package stores

import (
	. "github.com/Banyango/gifoody_server/api/model"
	. "github.com/Banyango/gifoody_server/api/repositories/util"
	"github.com/jmoiron/sqlx"
)

type PostSQLStore struct {
	db *sqlx.DB
}

func NewPostSQLStore(session *sqlx.DB) *PostSQLStore {
	self := new(PostSQLStore)

	self.db = session

	return self
}

func (self *PostSQLStore) FindPosts(query PostQuery) StoreChannel {
	var storeChan = make(StoreChannel, 1)
	go func() {
		var results []Post
		rows, err := self.db.Query(`SELECT p.Id, p.name, p.url, p.user_id, p.post_date  FROM posts p LIMIT ? OFFSET ?`, query.Limit, query.Offset);
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}
		defer rows.Close()

		for rows.Next() {
			post := Post{}
			err := rows.Scan(&post.ID, &post.Name, &post.URL, &post.UserID, &post.PostDate)
			if err != nil {
				storeChan <- StoreResult{Data: nil, Err: err}
				return
			}
			results = append(results, post)
		}

		var count int
		row := self.db.QueryRow("SELECT COUNT(*) FROM posts")
		err = row.Scan(&count)
		if err != nil {
			storeChan <- StoreResult{Data: nil, Err: err}
			return
		}

		storeChan <- StoreResult{Data: results, Total:count, Err: nil}
	}()
	return storeChan
}

