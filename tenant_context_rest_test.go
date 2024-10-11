package tenantcontextrest

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	con := newTenantContext(dbs, "X-Tenant-ID")

	router := gin.Default()

	router.Use(con.ChangeContextRest())

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

	con := newTenantContext(dbs, "X-Tenant-ID")

	router := gin.Default()

	router.Use(con.ChangeContextRest())

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

func TestChangeContextGrpcErro(t *testing.T) {

	create, err := Create("tenant1")
	if err != nil {
		t.Fatal(err)
	}

	create2, err := Create("tenant2")
	if err != nil {
		t.Fatal(err)
	}

	dbs := map[string]*gorm.DB{
		"tenant1": create,
		"tenant2": create2,
	}

	con := newTenantContext(dbs, "X-Tenant-ID")

	interceptor := con.ChangeContextGrpc()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {

		db, exists := ctx.Value("db").(*gorm.DB)

		assert.True(t, exists)

		if create == db {
			return "OK", nil
		}

		return nil, status.Errorf(codes.NotFound, "tenant não encontrado")
	}

	md := metadata.Pairs("X-Tenant-ID", "tenant5")

	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)

	assert.Error(t, err)

	s, ok := status.FromError(err)

	assert.True(t, ok)

	assert.Equal(t, codes.NotFound, s.Code())
}

func TestChangeContextGrpcSucesso(t *testing.T) {

	create, err := Create("tenant1")
	if err != nil {
		t.Fatal(err)
	}

	create2, err := Create("tenant2")
	if err != nil {
		t.Fatal(err)
	}

	dbs := map[string]*gorm.DB{
		"tenant1": create,
		"tenant2": create2,
	}

	con := newTenantContext(dbs, "X-Tenant-ID")

	interceptor := con.ChangeContextGrpc()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		db, exists := ctx.Value("db").(*gorm.DB)
		assert.True(t, exists)

		if create == db {
			return "OK", nil
		}

		return nil, nil
	}

	md := metadata.Pairs("X-Tenant-ID", "tenant1")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)

	assert.NoError(t, err)
	assert.Equal(t, "OK", resp)
}

func TestChangeContextGrpcErroNotFoundKey(t *testing.T) {

	create, err := Create("tenant1")
	if err != nil {
		t.Fatal(err)
	}

	create2, err := Create("tenant2")
	if err != nil {
		t.Fatal(err)
	}

	dbs := map[string]*gorm.DB{
		"tenant1": create,
		"tenant2": create2,
	}

	con := newTenantContext(dbs, "X-Tenant-ID")

	interceptor := con.ChangeContextGrpc()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {

		db, exists := ctx.Value("db").(*gorm.DB)
		assert.True(t, exists)

		if create == db {
			return "OK", nil
		}

		return nil, status.Errorf(codes.NotFound, "tenant não encontrado")
	}

	md := metadata.Pairs("X-Tenant-", "tenant5")

	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)

	assert.Error(t, err)

	s, ok := status.FromError(err)

	assert.True(t, ok)

	assert.Equal(t, codes.NotFound, s.Code())
}

func TestChangeContextGrpcErroNotDataLoss(t *testing.T) {

	create, err := Create("tenant1")
	if err != nil {
		t.Fatal(err)
	}

	create2, err := Create("tenant2")
	if err != nil {
		t.Fatal(err)
	}

	dbs := map[string]*gorm.DB{
		"tenant1": create,
		"tenant2": create2,
	}

	con := newTenantContext(dbs, "X-Tenant-ID")

	interceptor := con.ChangeContextGrpc()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		db, exists := ctx.Value("db").(*gorm.DB)
		assert.True(t, exists)

		if create == db {
			return "OK", nil
		}

		return nil, status.Errorf(codes.NotFound, "tenant não encontrado")
	}

	_, err = interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	assert.Error(t, err)

	s, ok := status.FromError(err)

	assert.True(t, ok)

	assert.Equal(t, codes.DataLoss, s.Code())
}
