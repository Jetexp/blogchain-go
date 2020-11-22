package content

import (
	"database/sql"
	"github.com/zikwall/blogchain/src/models"
	"github.com/zikwall/blogchain/src/models/tag"
	"github.com/zikwall/blogchain/src/models/user"
	"github.com/zikwall/blogchain/src/service"
)

type (
	ContentModel struct {
		models.BlogchainModel
	}
	Content struct {
		Id         int64          `db:"id"`
		Uuid       string         `db:"uuid"`
		UserId     int64          `db:"user_id"`
		Title      string         `db:"title"`
		Annotation string         `db:"annotation"`
		Content    string         `db:"content"`
		CreatedAt  sql.NullInt64  `db:"created_at"`
		UpdatedAt  sql.NullInt64  `db:"updated_at"`
		Image      sql.NullString `db:"image"`

		User user.User `db:"user"`
		Tags []tag.Tag
	}
	PublicContent struct {
		Id         int64  `json:"id"`
		Uuid       string `json:"uuid"`
		Title      string `json:"title"`
		Annotation string `json:"annotation"`
		Content    string `json:"content"`
		CreatedAt  int64  `json:"created_at"`
		UpdatedAt  int64  `json:"updated_at"`
		Image      string `json:"image"`

		Related Related `json:"related"`
	}
	Related struct {
		Publisher user.PublicUser `json:"publisher"`
		Tags      []tag.Tag       `json:"tags"`
	}
)

func NewContentModel(conn *service.BlogchainDatabaseInstance) ContentModel {
	return ContentModel{struct {
		Connection *service.BlogchainDatabaseInstance
	}{
		Connection: conn,
	},
	}
}

func (c *Content) Response() PublicContent {
	return PublicContent{
		Id:         c.Id,
		Uuid:       c.Uuid,
		Title:      c.Title,
		Annotation: c.Annotation,
		Content:    c.Content,
		CreatedAt:  c.CreatedAt.Int64,
		UpdatedAt:  c.UpdatedAt.Int64,
		Image:      c.Image.String,
		Related: Related{
			Publisher: c.User.Properties(),
			Tags:      c.Tags,
		},
	}
}
