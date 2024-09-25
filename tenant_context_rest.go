package tenantcontextrest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type tenantContextRest struct {
	db  *gorm.DB
	key string
}

func newTenantContextRest(db *gorm.DB, key string) tenantContextRest {

	return tenantContextRest{
		db:  db,
		key: key,
	}
}

func (t tenantContextRest) ChangeContext() gin.HandlerFunc {

	return func(c *gin.Context) {

		schema := c.GetHeader(t.key)

		err := t.changeSchema(schema)

		if err != nil {
			panic("An error occurred while switching to the scheme " + schema)
		}
	}
}

func (t tenantContextRest) changeSchema(schema string) error {

	log.Println("Trying to change schema")

	tx := t.db.Exec(fmt.Sprintf("SET search_path TO %s", schema))

	return tx.Error
}
