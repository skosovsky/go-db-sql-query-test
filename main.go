package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Client struct {
	ID       int
	FIO      string
	Login    string
	Birthday string
	Email    string
}

func (c Client) String() string {
	return fmt.Sprintf("ID: %d, FIO: %s, Login %s, Birthday: %s, Email %s",
		c.ID, c.FIO, c.Login, c.Birthday, c.Email)
}

type DB struct {
	*sql.DB
}

func NewDB() (DB, error) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		err = fmt.Errorf("db open error: %w", err)
		return DB{DB: nil}, err
	}

	return DB{DB: db}, nil
}

func (d DB) selectClient(id int) (Client, error) {
	var client Client

	row := d.QueryRow("SELECT id, fio, login, birthday, email FROM clients WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&client.ID, &client.FIO, &client.Login, &client.Birthday, &client.Email)
	if err != nil {
		err = fmt.Errorf("row scan error: %w", err)
		return Client{}, err
	}

	return client, nil
}

func (d DB) insertClient(client Client) (int, error) {
	res, err := d.Exec("INSERT INTO clients (fio, login, birthday, email) VALUES (:fio, :login, :birthday, :email)",
		sql.Named("fio", client.FIO),
		sql.Named("login", client.Login),
		sql.Named("birthday", client.Birthday),
		sql.Named("email", client.Email))
	if err != nil {
		err = fmt.Errorf("db exec error: %w", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("get last insert id error %w: ", err)
		return 0, err
	}

	return int(id), nil
}

func (d DB) updateClient(id int, client Client) error {
	_, err := d.Exec("UPDATE clients SET fio = :fio, login = :login, birthday = :birthday, email = :email  WHERE id = :id",
		sql.Named("id", id),
		sql.Named("fio", client.FIO),
		sql.Named("login", client.Login),
		sql.Named("birthday", client.Birthday),
		sql.Named("email", client.Email))
	if err != nil {
		err = fmt.Errorf("db exec error: %w", err)
		return err
	}

	return nil
}

func (d DB) deleteClient(id int) error {
	_, err := d.Exec("DELETE FROM clients WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		err = fmt.Errorf("db exec error: %w", err)
		return err
	}

	return nil
}

func main() {

}
