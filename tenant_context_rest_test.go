package tenantcontextrest

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChangeContextRest(t *testing.T) {

	create, err := Create("tenant1")
	if err != nil {
		t.Fatal(err)
	}

	create2, err := Create("tenant1")
	if err != nil {
		t.Fatal(err)
	}

	dbs := map[string]*gorm.DB{
		"tenant1": create,
		"tenant2": create2,
	}

	context := newTenantContext(dbs, "X-Tenant-ID")

	router := gin.Default()

	router.Use(context.ChangeContextRest())

	router.GET("/test", func(c *gin.Context) {

		db, exists := c.MustGet("db").(*gorm.DB)

		assert.True(t, exists)

		if create == db {

			c.String(http.StatusOK, "OK")

		} else {

			c.String(http.StatusNotFound, "not found")

		}

	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	req.Header.Set("X-Tenant-ID", "tenant1")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestChangeContextRestErro(t *testing.T) {

	create, err := Create("tenant1")
	if err != nil {
		t.Fatal(err)
	}

	create2, err := Create("tenant1")
	if err != nil {
		t.Fatal(err)
	}

	dbs := map[string]*gorm.DB{
		"tenant1": create,
		"tenant2": create2,
	}

	context := newTenantContext(dbs, "X-Tenant-ID")

	router := gin.Default()

	router.Use(context.ChangeContextRest())

	router.GET("/test", func(c *gin.Context) {

		db, exists := c.MustGet("db").(*gorm.DB)

		assert.True(t, exists)

		if create == db {

			c.String(http.StatusOK, "OK")

		}

	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	req.Header.Set("X-Tenant-ID", "tenant5")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

}

func Create(schema string) (*gorm.DB, error) {

	db, _, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  schema,
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return gormDB, nil
}
