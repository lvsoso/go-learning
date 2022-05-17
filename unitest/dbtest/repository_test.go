package dbtest

import (
	"database/sql"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAll(t *testing.T) {
	var repository *Repository
	var mock sqlmock.Sqlmock
	var db *sql.DB
	var err error

	db, mock, err = sqlmock.New()
	assert.Nil(t, err)
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	repository = &Repository{db: gdb}

	// list all
	const sqlSelectAll = "SELECT * FROM `blogs"

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
		WillReturnRows(sqlmock.NewRows(nil))

	l, err := repository.ListAll()
	assert.Nil(t, err)
	assert.Equal(t, []*Blog{}, l)

	// select one
	blog := &Blog{
		ID:        1,
		Title:     "post",
		Content:   "hello",
		Tags:      pq.StringArray{"go", "golang"},
		CreatedAt: time.Now(),
	}

	rows := sqlmock.
		NewRows([]string{"id", "title", "content", "tags", "created_at"}).
		AddRow(blog.ID, blog.Title, blog.Content, blog.Tags, blog.CreatedAt)

	const sqlSelectOne = "SELECT * FROM `blogs` WHERE id = ? ORDER BY `blogs`.`id` LIMIT 1"

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
		WithArgs(blog.ID).
		WillReturnRows(rows)

	dbBlog, err := repository.Load(blog.ID)
	assert.Nil(t, err)
	assert.Equal(t, blog, dbBlog)

	// list
	rows = sqlmock.
		NewRows([]string{"id", "title", "content", "tags", "created_at"}).
		AddRow(1, "post 1", "hello 1", nil, time.Now()).
		AddRow(2, "post 2", "hello 2", pq.StringArray{"go"}, time.Now())
	const sqlSelectFirstTen = "SELECT * FROM `blogs` LIMIT 10"
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectFirstTen)).WillReturnRows(rows)

	l, err = repository.List(0, 10)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(l))
	assert.Equal(t, pq.StringArray(pq.StringArray(nil)), l[0].Tags)
	assert.Equal(t, pq.StringArray(pq.StringArray{"go"}), l[1].Tags)
	assert.Equal(t, uint(2), l[1].ID)

	// update
	blog = &Blog{
		Title:     "post",
		Content:   "hello",
		Tags:      pq.StringArray{"a", "b"},
		CreatedAt: time.Now(),
	}

	const sqlUpdate = "UPDATE `blogs` SET `title`=?,`content`=?,`tags`=?,`created_at`=? WHERE `id` = ?"
	const sqlSelectOne2 = "SELECT * FROM `blogs` WHERE `id` = ? LIMIT 1"

	blog.ID = 1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(blog.Title, blog.Content, blog.Tags, blog.CreatedAt, blog.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne2)).
		WithArgs(blog.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(blog.ID))

	err = repository.Save(blog)
	assert.Nil(t, err)

	// insert
	blog = &Blog{
		Title:     "post",
		Content:   "hello",
		Tags:      pq.StringArray{"a", "b"},
		CreatedAt: time.Now(),
	}
	const sqlInsert = "INSERT INTO `blogs` (`title`,`content`,`tags`,`created_at`) VALUES (?,?,?,?)"
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sqlInsert)).
		WithArgs(blog.Title, blog.Content, blog.Tags, blog.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repository.Save(blog)
	assert.Nil(t, err)

	// search
	rows = sqlmock.
		NewRows([]string{"id", "title", "content", "tags", "created_at"}).
		AddRow(1, "post 1", "hello 1", nil, time.Now())

	// limit/offset is not parameter
	const sqlSearch = "SELECT * FROM `blogs` WHERE title like ? LIMIT 10"
	const q = "os"

	mock.ExpectQuery(regexp.QuoteMeta(sqlSearch)).
		WithArgs("%" + q + "%").
		WillReturnRows(rows)

	l, err = repository.SearchByTitle(q, 0, 10)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(l))
	assert.Equal(t, true, strings.Contains(l[0].Title, q))
}
