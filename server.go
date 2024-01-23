package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/cmd"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/config"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/directives"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/graph"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/middlewares"
	"github.com/myanmarmarathon/mkitchen-distribution-backend-graphql/models"
	"github.com/ravilushqa/otelgqlgen"
	"go.opentelemetry.io/otel"
)

const defaultPort = "8000"

var tracer = otel.Tracer("mkitchen-distribution")

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file

	c := graph.Config{Resolvers: &graph.Resolver{
		Tracer: tracer,
	}}
	c.Directives.Auth = directives.Auth

	h := handler.NewDefaultServer(graph.NewExecutableSchema(c))
	h.Use(otelgqlgen.Middleware())
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func customNotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
}

func main() {

	cmd.Execute()
	
	port := os.Getenv("API_PORT")
	if port == "" {
		port = defaultPort
	}

	// Connect to Database
	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	models.MigrateTable()

	

	// Setup UpTrace Logging with OpenTelemetry
	// upTraceDsn := os.Getenv("UPTRACE_DSN")
	// ctx := context.Background()
	// uptrace.ConfigureOpentelemetry(
	// 	uptrace.WithDSN(upTraceDsn),
	// )
	// defer uptrace.Shutdown(ctx)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	r.Use(cors.New(config))
	r.Use(middlewares.AuthMiddleware())
	r.Use(middlewares.LoaderMiddleware())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.NoRoute(customNotFoundHandler)
	r.Run(":" + port)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
}
