package tests

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
	"github.com/mdelclaro/gobrax/src/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	app *fiber.App
	db  *sql.DB
	now       = time.Time{}
	id  int32 = 1
)

func TestMain(m *testing.M) {
	app = utils.SetupApp()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestTruckHandlers(t *testing.T) {
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
			name:         "[Success] - Test Get All Trucks",
			route:        "/api/truck",
			method:       "GET",
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": []entities.Truck{{
					GormModel: entities.GormModel{
						ID: 1,
					},
					LicensePlate:     "123",
					FuelUsed:         decimal.NewFromInt(0),
					DistanceTraveled: decimal.NewFromInt(0),
					DriverID:         &id,
					Driver: &entities.Driver{
						GormModel: entities.GormModel{
							ID: 1,
						},
						Name:          "driver",
						LicenseNumber: "123",
						IsActive:      true,
					},
				}},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				trucks := sqlmock.NewRows([]string{
					"id", "license_plate", "fuel_used", "distance_traveled", "driver_id", "Driver__id", "Driver__name", "Driver__license_number", "Driver__is_active",
				}).
					AddRow(1, "123", "0", "0", 1, 1, "driver", "123", true)

				expectedSQL := "SELECT (.+) FROM \"trucks\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(trucks)
			},
		},
		{
			name:         "[Success] - Test Get Truck By Id",
			route:        fmt.Sprintf("/api/truck/%d", id),
			method:       "GET",
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": entities.Truck{
					GormModel: entities.GormModel{
						ID: id,
					},
					LicensePlate:     "123",
					FuelUsed:         decimal.NewFromInt(0),
					DistanceTraveled: decimal.NewFromInt(0),
					DriverID:         &id,
					Driver: &entities.Driver{
						GormModel: entities.GormModel{
							ID: 1,
						},
						Name:          "driver",
						LicenseNumber: "123",
						IsActive:      true,
					},
				},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				trucks := sqlmock.NewRows([]string{
					"id", "license_plate", "fuel_used", "distance_traveled", "driver_id", "Driver__id", "Driver__name", "Driver__license_number", "Driver__is_active",
				}).
					AddRow(id, "123", "0", "0", 1, 1, "driver", "123", true)

				expectedSQL := "SELECT (.+) FROM \"trucks\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(trucks)
			},
		},
		{
			name:         "[Invalid] - Test Get Truck By Id With Invalid Id",
			route:        fmt.Sprintf("/api/truck/%s", "INVALID"),
			method:       "GET",
			expectedCode: 400,
			expectedBody: helpers.BuildError(fmt.Errorf("invalid id provided: %s", errors.New("strconv.Atoi: parsing \"INVALID\": invalid syntax"))),
			mock:         func() {},
		},
		{
			name:   "[Success] - Test Add Truck",
			route:  "/api/truck",
			method: "POST",
			body: entities.Truck{
				GormModel: entities.GormModel{
					ID: id,
				},
				LicensePlate:     "123",
				FuelUsed:         decimal.NewFromInt(0),
				DistanceTraveled: decimal.NewFromInt(0),
			},
			expectedCode: 201,
			expectedBody: map[string]any{
				"data": entities.Truck{
					GormModel: entities.GormModel{
						ID:        id,
						CreatedAt: now,
						UpdatedAt: now,
					},
					LicensePlate:     "123",
					FuelUsed:         decimal.NewFromInt(0),
					DistanceTraveled: decimal.NewFromInt(0),
					DriverID:         nil,
				},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				row := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "license_plate", "fuel_used", "distance_traveled",
				}).
					AddRow(id, now, now, "123", "0", "0")

				expectedSQL := "INSERT INTO \"trucks\" (.+) VALUES (.+)"
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
				mock.ExpectCommit()
			},
		},
		{
			name:   "[Success] - Test Update Truck",
			route:  "/api/truck",
			method: "PUT",
			body: entities.Truck{
				GormModel: entities.GormModel{
					ID: id,
				},
				LicensePlate:     "456",
				FuelUsed:         decimal.NewFromInt(0),
				DistanceTraveled: decimal.NewFromInt(0),
				DriverID:         &id,
			},
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": entities.Truck{
					GormModel: entities.GormModel{
						ID:        id,
						CreatedAt: now,
						UpdatedAt: now,
					},
					LicensePlate:     "456",
					FuelUsed:         decimal.NewFromInt(0),
					DistanceTraveled: decimal.NewFromInt(0),
					DriverID:         &id,
					Driver: &entities.Driver{
						GormModel: entities.GormModel{
							ID: 1,
						},
						Name:          "driver",
						LicenseNumber: "123",
						IsActive:      true,
					},
				},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				expectedSQL := "UPDATE \"trucks\" SET .+"
				row := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "license_plate", "fuel_used", "distance_traveled",
				}).
					AddRow(id, now, now, "123", "0", "0")
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
				mock.ExpectCommit()

				trucks := sqlmock.NewRows([]string{
					"id", "license_plate", "fuel_used", "distance_traveled", "driver_id", "Driver__id", "Driver__name", "Driver__license_number", "Driver__is_active",
				}).
					AddRow(id, "456", "0", "0", 1, 1, "driver", "123", true)

				expectedSQL = "SELECT (.+) FROM \"trucks\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(trucks)
			},
		},
		{
			name:         "[Success] - Test Delete Truck",
			route:        fmt.Sprintf("/api/truck/%d", id),
			method:       "DELETE",
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": "",
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				expectedSQL := "DELETE FROM \"trucks\" .+"

				mock.ExpectBegin()
				mock.ExpectExec(expectedSQL).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:   "[Success] - Test Update Truck Driver",
			route:  fmt.Sprintf("/api/truck/update-driver/%v?driverId=%v", id, id),
			method: "POST",
			body: entities.Truck{
				GormModel: entities.GormModel{
					ID: id,
				},
				LicensePlate:     "456",
				FuelUsed:         decimal.NewFromInt(0),
				DistanceTraveled: decimal.NewFromInt(0),
				DriverID:         &id,
			},
			expectedCode: 200,
			expectedBody: map[string]any{
				"data": entities.Truck{
					GormModel: entities.GormModel{
						ID:        id,
						CreatedAt: now,
						UpdatedAt: now,
					},
					LicensePlate:     "456",
					FuelUsed:         decimal.NewFromInt(0),
					DistanceTraveled: decimal.NewFromInt(0),
					DriverID:         &id,
					Driver: &entities.Driver{
						GormModel: entities.GormModel{
							ID: 1,
						},
						Name:          "driver",
						LicenseNumber: "123",
						IsActive:      true,
					},
				},
			},
			mock: func() {
				dbConn, _, mock := database.StartDbMock(t)
				db = dbConn

				// find truck by id
				truck := sqlmock.NewRows([]string{
					"id", "license_plate", "fuel_used", "distance_traveled", "driver_id", "Driver__id", "Driver__name", "Driver__license_number", "Driver__is_active",
				}).
					AddRow(id, "123", "0", "0", 1, 1, "driver", "123", true)

				expectedSQL := "SELECT (.+) FROM \"trucks\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(truck)

				// find driver by id
				driver := sqlmock.NewRows([]string{
					"id", "name", "license_number", "is_active",
				}).
					AddRow(id, "driver", "123", true)

				expectedSQL = "SELECT (.+) FROM \"drivers\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(driver)

				// update truck driver
				expectedSQL = "UPDATE \"trucks\" SET .+"
				row := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "license_plate", "fuel_used", "distance_traveled",
				}).
					AddRow(id, now, now, "123", "0", "0")
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
				mock.ExpectCommit()

				// find updated truck
				trucks := sqlmock.NewRows([]string{
					"id", "license_plate", "fuel_used", "distance_traveled", "driver_id", "Driver__id", "Driver__name", "Driver__license_number", "Driver__is_active",
				}).
					AddRow(id, "456", "0", "0", 1, 1, "driver", "123", true)

				expectedSQL = "SELECT (.+) FROM \"trucks\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(trucks)
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
