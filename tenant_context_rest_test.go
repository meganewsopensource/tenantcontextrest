package tenantcontextrest

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestTenantContextRest_ChangeContext(t *testing.T) {
	var key = " X-Tenant-ID"
	var schema = "empresaleve"

	db, mock, err := Create(schema)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.Use(newTenantContextRest(db, key).ChangeContext())
	r.GET("/", func(context *gin.Context) {

		context.JSON(200, gin.H{
			key: context.Request.Header[key],
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add(key, schema)

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func Create(schema string) (*gorm.DB, sqlmock.Sqlmock, error) {

	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf("SET search_path TO %s", schema))).
		WillReturnResult(sqlmock.NewResult(0, 0))

	return gormDB, mock, nil
}
