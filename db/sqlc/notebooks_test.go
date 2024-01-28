package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateNotebook(t *testing.T) {
	user := createRandomUser(t)
	testUUID := uuid.MustParse("3560f394-1996-434d-8e0b-755a0f3601b5")
	arg := CreateNotebookParams{
		ID:        testUUID,
		Title:     "Title1",
		Topic:     "Topic1",
		Content:   "Content 1",
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}

	fmt.Println(arg)

	notebook, err := testQueries.CreateNotebook(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, notebook)

	require.Equal(t, arg.ID, notebook.ID)
	require.NotZero(t, notebook.CreatedAt)
	require.NotZero(t, notebook.LastModified)
}

func TestGetNotebook(t *testing.T) {
	testUUID := uuid.MustParse("3560f394-1996-434d-8e0b-755a0f3601b5")

	arg := GetNotebookParams{
		ID:     testUUID,
		UserID: testData.userId1,
	}
	notebook, err := testQueries.GetNotebook(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, notebook)
}

func TestUpdateNotebook(t *testing.T) {
	testUUID := uuid.MustParse("3560f394-1996-434d-8e0b-755a0f3601b5")

	getNotebookParams := GetNotebookParams{
		ID:     testUUID,
		UserID: testData.userId1,
	}
	notebook, err := testQueries.GetNotebook(context.Background(), getNotebookParams)
	require.NoError(t, err)
	require.NotEmpty(t, notebook)

	// update content
	notebook.Content = "Content 2"

	arg := UpdateNotebookParams{
		ID:           testUUID,
		UserID:       testData.userId1,
		Title:        notebook.Title,
		Content:      notebook.Content,
		Topic:        notebook.Topic,
		LastModified: time.Now(),
	}

	updatedNotebook, err := testQueries.UpdateNotebook(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedNotebook)
	require.Equal(t, notebook.Content, updatedNotebook.Content)
}

func TestListNotebooks(t *testing.T) {
	arg := ListNotebooksParams{
		UserID: testData.userId1,
		Limit:  10,
		Offset: 0,
	}
	notebooks, err := testQueries.ListNotebooks(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, notebooks)
	require.Equal(t, 1, len(notebooks))
}

func TestDeleteNotebook(t *testing.T) {
	testUUID := uuid.MustParse("3560f394-1996-434d-8e0b-755a0f3601b5")
	arg := DeleteNotebookParams{
		ID:           testUUID,
		UserID:       testData.userId1,
		LastModified: time.Now(),
	}
	deletedNotebook, err := testQueries.DeleteNotebook(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, deletedNotebook)
}
