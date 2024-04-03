package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *DB) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	clientID := 1

	client, err := db.selectClient(clientID)
	require.NoError(t, err)

	expectedID := clientID
	actualID := client.ID
	assert.Equal(t, expectedID, actualID)

	assert.NotEmptyf(t, client.FIO, client.Login, client.Birthday, client.Email)
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *DB) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	clientID := -1

	client, err := db.selectClient(clientID)
	require.Error(t, err)

	expectedErr := sql.ErrNoRows.Error()
	actualErr := err
	require.ErrorContains(t, actualErr, expectedErr)

	assert.Emptyf(t, client.ID, client.FIO, client.Login, client.Birthday, client.Email)
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *DB) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	client := Client{ //nolint:exhaustruct // ID gets DB
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	client.ID, err = db.insertClient(client)
	require.True(t,
		assert.NoError(t, err),
		assert.NotEmpty(t, client.ID))

	clientDataVerification, err := db.selectClient(client.ID)
	require.NoError(t, err)

	assert.EqualValues(t, clientDataVerification, client)

	err = db.deleteClient(client.ID)
	require.NoError(t, err)
}

func Test_InsertClient_ThenUpdateAndSelectAndCheck(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *DB) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	client := Client{ //nolint:exhaustruct // ID gets DB
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	client.ID, err = db.insertClient(client)
	require.True(t,
		assert.NoError(t, err),
		assert.NotEmpty(t, client.ID))

	clientDataUpdate := Client{
		ID:       client.ID,
		FIO:      "TestNew",
		Login:    "TestNew",
		Birthday: "19700102",
		Email:    "test@mail.com",
	}

	err = db.updateClient(clientDataUpdate.ID, clientDataUpdate)
	require.NoError(t, err)

	clientDataVerification, err := db.selectClient(client.ID)
	require.NoError(t, err)

	assert.EqualValues(t, clientDataVerification, clientDataUpdate)

	err = db.deleteClient(client.ID)
	require.NoError(t, err)
}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *DB) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	client := Client{ //nolint:exhaustruct // ID gets DB
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	client.ID, err = db.insertClient(client)
	require.True(t,
		assert.NoError(t, err),
		assert.NotEmpty(t, client.ID))

	_, err = db.selectClient(client.ID)
	require.NoError(t, err)

	err = db.deleteClient(client.ID)
	require.NoError(t, err)

	_, err = db.selectClient(client.ID)
	require.Error(t, err)

	expectedErr := sql.ErrNoRows.Error()
	actualErr := err
	assert.ErrorContains(t, actualErr, expectedErr)
}
