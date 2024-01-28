package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/dpomian/gobind/db/mock"
	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/dpomian/gobind/token"
	"github.com/dpomian/gobind/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListNotebookByIdHandler(t *testing.T) {
	user := randomUser(t)
	notebook := randomNotebook(user.ID)

	testCases := []struct {
		name          string
		id            string
		title         string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(storage *mockdb.MockStorage)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   notebook.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authTypeBearer, user.ID.String(), time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				// build stubs
				storage.EXPECT().
					GetNotebook(gomock.Any(), gomock.Eq(db.GetNotebookParams{ID: notebook.ID, UserID: user.ID})).
					Times(1).
					Return(notebook, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchNotebook(t, recorder.Body, notebook)
			},
		},
		{
			name: "NotFound",
			id:   notebook.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authTypeBearer, user.ID.String(), time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				// build stubs
				storage.EXPECT().
					GetNotebook(gomock.Any(), gomock.Eq(db.GetNotebookParams{ID: notebook.ID, UserID: user.ID})).
					Times(1).
					Return(db.Notebook{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidId",
			id:   "this-is-not-a-valid-uuid",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authTypeBearer, user.ID.String(), time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				// build stubs
				storage.EXPECT().
					GetNotebook(gomock.Any(), gomock.Eq(notebook.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   notebook.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authTypeBearer, user.ID.String(), time.Minute)
			},
			buildStubs: func(storage *mockdb.MockStorage) {
				// build stubs
				storage.EXPECT().
					GetNotebook(gomock.Any(), gomock.Eq(db.GetNotebookParams{ID: notebook.ID, UserID: user.ID})).
					Times(1).
					Return(db.Notebook{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		fmt.Println("running test:", tc.name)

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			storage := mockdb.NewMockStorage(ctrl)
			tc.buildStubs(storage)

			// start server and send http request
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/notebooks/%s", tc.id)
			rq, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server := newTestServer(t, storage)
			require.NoError(t, err)

			tc.setupAuth(t, rq, server.tokenMaker)
			server.router.ServeHTTP(recorder, rq)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomNotebook(userId uuid.UUID) db.Notebook {
	return db.Notebook{
		ID:           uuid.New(),
		Title:        utils.RandomString(10),
		Topic:        utils.RandomString(5),
		Content:      utils.RandomString(50),
		UserID:       userId,
		Deleted:      false,
		LastModified: time.Now(),
		CreatedAt:    time.Now(),
	}
}

func requireBodyMatchNotebook(t *testing.T, body *bytes.Buffer, notebook db.Notebook) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotNotebook db.Notebook
	err = json.Unmarshal(data, &gotNotebook)
	require.NoError(t, err)
	require.Equal(t, notebook.ID, gotNotebook.ID)
	require.Equal(t, notebook.Title, gotNotebook.Title)
	require.Equal(t, notebook.Topic, gotNotebook.Topic)
	require.Equal(t, notebook.Content, gotNotebook.Content)
	require.Equal(t, notebook.Deleted, gotNotebook.Deleted)
}
