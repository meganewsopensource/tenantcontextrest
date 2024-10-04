package tenantcontextrest

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
)

type repo interface {
	ChangeContextRest()
	ChangeContextGrpc()
}

type tenantContext struct {
	dbs map[string]*gorm.DB
	key string
}

func newTenantContext(dbs map[string]*gorm.DB, key string) tenantContext {

	return tenantContext{
		dbs: dbs,
		key: key,
	}
}

func (t tenantContext) ChangeContextRest() gin.HandlerFunc {

	return func(c *gin.Context) {

		schema := c.GetHeader(t.key)

		db := t.dbs[schema]

		c.Set("db", db)

	}
}

func (t tenantContext) ChangeContextGrpc() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.DataLoss, "Metadados ausentes")
		}

		schemas := md.Get(t.key)
		if len(schemas) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "%s de autorização ausente", t.key)
		}

		schema := schemas[0]

		log.Printf("Schema: %s", schema)

		db := t.dbs[schema]

		newCtx := context.WithValue(ctx, "db", db)

		return handler(newCtx, req)

	}
}
