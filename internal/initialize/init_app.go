package initialize

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	"quiz-app/internal/infrastructure/persistence/aws"
	persistence "quiz-app/internal/infrastructure/persistence/mongodb"
	redisdb "quiz-app/internal/infrastructure/persistence/redis"
	routes "quiz-app/internal/infrastructure/router"

	"github.com/rs/cors"
)

var (
	dbName = "dbapp"
)

func InitApp() {
	//Init mongo
	persistence.ConnectMongoDB(os.Getenv(persistence.MongoConnectionString))

	userRepo := persistence.NewUserMongoRepository()

	user, err := userRepo.GetUser(context.TODO(), &entity.User{
		EmailID: "ZG1mwdlEFtRwxXRezbmgf6Ctij13",
	})

	fmt.Println(user, "ERROR: ", err)

	InitRouter()
}
func InitRouter() {
	// Initialize CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	router := routes.NewRouter()
	routes.GetRoutesAuth(router).SetRouteLogin()
	// Initialize repositories
	classRepo := persistence.NewClassMongoRepository()
	testRepo := persistence.NewTestMongoRepository()
	questionRepo := persistence.NewQuestionMongoRepository()
	fileRepo := persistence.NewFileMongoRepository()
	answerRepo := persistence.NewAnswerMongoRepository()

	redisRepo, err := redisdb.GetRedis()

	if err != nil {
		fmt.Println(err)
		fmt.Println("NOT RUN REDIS")
	}

	// Initialize use cases
	classUseCase := service.NewClassUseCase(classRepo, testRepo)
	questionUseCase := service.NewQuestionUseCase(questionRepo)
	testUseCase := service.NewTestUseCase(testRepo)
	fileUseCase := service.NewFileUseCase(fileRepo)
	answerUseCase := service.NewAnswerUseCase(answerRepo)
	redisUseCase := service.NewRedisUseCase(redisRepo)

	awsS3UseCase := aws.NewFileAWSRepository("quiz-app-image-storage", "ap-southeast-2")

	routes.NewRouterTest(*testUseCase, *classUseCase, *questionUseCase, *answerUseCase, *redisUseCase).GetTestRouter(router)
	routes.NewRouterQuestion(*questionUseCase).GetQuestionRouter(router)
	routes.NewRouterClass(*classUseCase, *redisUseCase).GetClassRouter(router)
	routes.NewRoutesFile(fileUseCase, awsS3UseCase).GetRoutesFile(router)
	routes.NewRouterAnswer(answerUseCase, testUseCase, questionUseCase, redisUseCase).GetAnswerRouter(router)

	// Apply CORS handler
	handler := c.Handler(router)

	// Start the server
	port := "8080"
	log.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))

}
