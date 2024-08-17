package driver

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v3"
	"github.com/mdelclaro/gobrax/src/api/helpers"
	database "github.com/mdelclaro/gobrax/src/db"
	"github.com/mdelclaro/gobrax/src/repository/entities"
	"github.com/stretchr/testify/assert"
)

var (
	app *fiber.App
	db  *sql.DB
	now       = time.Time{}
	id  int32 = 1
)

func TestMain(m *testing.M) {
	app = fiber.New()
	api := app.Group("/api")

	SetupDriverRoutes(api)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestDriverHandlers(t *testing.T) {
	tests := []struct {
		name string

		route  string
		method string
		body   any

		expectedCode int
		expectedBody any

		mock func()
	}{
		{
			name:         "[Success] - Test Get All Drivers",
			route:        "/api/driver",
			method:       "GET",
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": []entities.Driver{
					{
						GormModel: entities.GormModel{
							ID: id,
						},
						Name:          "name",
						LicenseNumber: "123",
						IsActive:      true,
					},
				},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				drivers := sqlmock.NewRows([]string{
					"id", "name", "license_number", "is_active",
				}).
					AddRow(id, "name", "123", true)

				expectedSQL := "SELECT (.+) FROM \"drivers\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(drivers)
			},
		},
		{
			name:         "[Success] - Test Get Driver By Id",
			route:        fmt.Sprintf("/api/driver/%d", id),
			method:       "GET",
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": entities.Driver{
					GormModel: entities.GormModel{
						ID: id,
					},
					Name:          "name",
					LicenseNumber: "123",
					IsActive:      true,
				},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				driver := sqlmock.NewRows([]string{
					"id", "name", "license_number", "is_active",
				}).
					AddRow(id, "name", "123", true)

				expectedSQL := "SELECT (.+) FROM \"drivers\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(driver)
			},
		},
		{
			name:         "[Invalid] - Test Get Driver By Id With Invalid Id",
			route:        fmt.Sprintf("/api/driver/%s", "INVALID"),
			method:       "GET",
			expectedCode: 400,
			expectedBody: helpers.BuildError(fmt.Errorf("invalid id provided: %s", errors.New("strconv.Atoi: parsing \"INVALID\": invalid syntax"))),
			mock:         func() {},
		},
		{
			name:   "[Success] - Test Add Driver",
			route:  "/api/driver",
			method: "POST",
			body: entities.Driver{
				GormModel: entities.GormModel{
					ID: id,
				},
				Name:          "name",
				LicenseNumber: "123",
				IsActive:      true,
			},
			expectedCode: 201,
			expectedBody: map[string]any{
				"data": entities.Driver{
					GormModel: entities.GormModel{
						ID: id,
					},
					Name:          "name",
					LicenseNumber: "123",
					IsActive:      true,
				},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				row := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "name", "license_number", "is_active",
				}).
					AddRow(id, now, now, "name", "123", true)

				expectedSQL := "INSERT INTO \"drivers\" (.+) VALUES (.+)"
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
				mock.ExpectCommit()
			},
		},
		{
			name:   "[Success] - Test Update Driver",
			route:  "/api/driver",
			method: "PUT",
			body: entities.Driver{
				GormModel: entities.GormModel{
					ID: id,
				},
				Name: "name_edited",
			},
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": entities.Driver{
					GormModel: entities.GormModel{
						ID: id,
					},
					Name:          "name_edited",
					LicenseNumber: "123",
					IsActive:      true,
				},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				expectedSQL := "UPDATE \"drivers\" SET .+"
				row := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "name", "license_number", "is_active",
				}).
					AddRow(id, now, now, "name_edited", "123", true)
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
				mock.ExpectCommit()

				driver := sqlmock.NewRows([]string{
					"id", "name", "license_number", "is_active",
				}).
					AddRow(id, "name", "123", true)

				expectedSQL = "SELECT (.+) FROM \"drivers\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(driver)
			},
		},
		{
			name:         "[Success] - Test Delete Driver",
			route:        fmt.Sprintf("/api/driver/%d", id),
			method:       "DELETE",
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": "",
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				expectedSQL := "DELETE FROM \"drivers\" .+"

				mock.ExpectBegin()
				mock.ExpectExec(expectedSQL).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			defer db.Close()

			reqBody, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			bodyReader := bytes.NewReader(reqBody)

			req, _ := http.NewRequest(
				tt.method,
				tt.route,
				bodyReader,
			)

			res, err := app.Test(req, -1)
			assert.NoError(t, err)

			body, _ := io.ReadAll(res.Body)
			parsedBody, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			assert.Equal(t, string(parsedBody), string(body))
			assert.Equal(t, tt.expectedCode, res.StatusCode)
		})
	}
}
