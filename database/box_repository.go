package database

import (
	"database/sql"

	"github.com/wakieu/drtbox/entity"
)

type BoxRepository struct {
	Db *sql.DB
}

func NewBoxRepository(db *sql.DB) *BoxRepository {
	return &BoxRepository{
		Db: db,
	}
}

func (r *BoxRepository) Save(box *entity.Box) error {
	if box.IsEmpty() {
		return nil
	}
	_, err := r.Db.Exec("INSERT INTO box(boxpath, text) VALUES(?, ?);",
		box.BoxPath, box.Text)
	if err != nil {
		return err
	}
	return nil
}

func (r *BoxRepository) GetContent(path string) (entity.Box, error) {
	var box entity.Box
	err := r.Db.QueryRow("SELECT * FROM box WHERE boxpath = ?;", path).Scan(
		&box.BoxPath, &box.Text,
	)
	box.BoxPath = path
	if err != nil {
		if err == sql.ErrNoRows {
			return box, nil
		} else {
			return box, err
		}
	}
	return box, nil
}

func (r *BoxRepository) Exists(path string) (bool, error) {
	var box entity.Box
	err := r.Db.QueryRow("SELECT * FROM box WHERE boxpath = ?;", path).Scan(
		&box.BoxPath, &box.Text,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (r *BoxRepository) Write(path string, text string) error {
	_, err := r.Db.Exec("UPDATE box SET text = ? WHERE boxpath = ?;", text, path)
	if err != nil {
		return err
	}
	return nil
}

func (r *BoxRepository) Delete(path string) error {
	_, err := r.Db.Exec("DELETE from box WHERE boxpath = ?;", path)
	if err != nil {
		return err
	}
	return nil
}

func (r *BoxRepository) GetChildren(path string) ([]string, error) {
	var children []string
	match := path + "/%"
	rows, err := r.Db.Query("SELECT boxpath FROM box WHERE boxpath LIKE ?;", match)
	if err != nil {
		return children, err
	}

	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			return children, err
		}
		children = append(children, row)
	}
	if err = rows.Err(); err != nil {
		return children, err
	}

	return children, nil
}
