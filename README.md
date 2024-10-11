## Tenant Context : Middleware para Gerenciamento de Múltiplos Locatários com Schemas no PostgreSQL (Baseado em Headers)

Tenant Context  é um middleware projetado  com o objetivo de simplificar a gestão de múltiplos locatários em aplicações que utilizam o banco de dados PostgreSQL. A biblioteca ajuda a implementar a abordagem de isolamento de dados por schema, onde cada locatário possui um schema exclusivo no banco de dados. A troca de schema é realizada com base em um header específico presente na requisição HTTP.

### Funcionalidade Principal:
* Troca Automática de Schema Baseada em Header: O middleware intercepta as requisições HTTP e busca um header predefinido (por exemplo, X-Tenant-ID) para identificar o locatário associado. Antes que o controlador da sua aplicação acesse o banco de dados, o middleware realiza a troca do schema ativo na conexão, garantindo que todas as operações subsequentes sejam executadas no contexto correto do locatário.

### Exemplo gin 


~~~ go 

    dbs := map[string]*gorm.DB{
        "tenant1": dbTenant1,
        "tenant2": dbTenant1,
        }

    context := newTenantContext(dbs, "X-Tenant-ID")
    
    router := gin.Default()

    router.Use(context.ChangeContextRest())	
	
	
    router.GET("/test", func(c *gin.Context) {

        db, exists := c.MustGet("db").(*gorm.DB)

        var employee []Employee
 	
		if exists{
	 
		    db.Find(&employee)     	
		}
		
        c.JSON(http.StatusOK, employee)
	
    })

~~~

### Exemplo GRPc

~~~ go 

    dbs := map[string]*gorm.DB{
        "tenant1": dbTenant1,
        "tenant2": dbTenant1,
        }

    context := newTenantContext(dbs, "X-Tenant-ID")
    
    l, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("run on port: 8080")

    server := grpc.NewServer(grpc.UnaryInterceptor(ChangeContextGrpc()))

    pb.RegisterMegaMDFeServer(server, NewMyController())
    
    err = server.Serve(l)
    if err != nil {
    	log.Fatal(err)
    }

   /////////
   
   func (M MyController) resource(ctx context.Context, c *pb.mymessage) (*emptypb.Empty, error) {

        db := ctx.Value("db").(*gorm.db)

        db.ping()
        
        return &emptypb.Empty{}, nil
}


~~~