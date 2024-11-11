package initialize

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	//Init router
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
	redisRepo, err := redisdb.GetRedis()
	if err != nil {
		fmt.Println(err)
		fmt.Println("NOT RUN REDIS")
	}
	redisUseCase := service.NewRedisUseCase(redisRepo)
	os.Getenv("JWT_SECRET")
	//authService, err := service.NewAuthService("firebase-config.json", ([]byte)("jwt-secret"))

	authService, err := service.NewAuthService("firebase-config.json", ([]byte)(os.Getenv("JWT_SECRET")))

	authUseCase := service.NewAuthUseCase(authService, *redisUseCase)

	authHandler := service.NewAuthHandler(authUseCase)

	// Initialize repositories
	userRepo := persistence.NewUserMongoRepository()
	classRepo := persistence.NewClassMongoRepository()
	testRepo := persistence.NewTestMongoRepository()
	questionRepo := persistence.NewQuestionMongoRepository()
	fileRepo := persistence.NewFileMongoRepository()
	answerRepo := persistence.NewAnswerMongoRepository()

	// Initialize use cases
	userUseCase := service.NewUserUseCase(userRepo)
	classUseCase := service.NewClassUseCase(classRepo, testRepo)
	questionUseCase := service.NewQuestionUseCase(questionRepo)
	testUseCase := service.NewTestUseCase(testRepo)
	fileUseCase := service.NewFileUseCase(fileRepo)
	answerUseCase := service.NewAnswerUseCase(answerRepo)

	awsS3UseCase := aws.NewFileAWSRepository("quiz-app-image-storage", "ap-southeast-2")

	routes.NewRoutesAuth(router, *authService, *userUseCase, *redisUseCase).SetLoginRoute()
	routes.NewRouterTest(*testUseCase, *classUseCase, *questionUseCase, *answerUseCase, *redisUseCase, *authHandler).GetTestRouter(router)
	routes.NewRouterQuestion(*questionUseCase, *authHandler).GetQuestionRouter(router)
	routes.NewRouterClass(*classUseCase, *redisUseCase, *authHandler).GetClassRouter(router)
	routes.NewRoutesFile(fileUseCase, awsS3UseCase, authHandler).GetRoutesFile(router)
	routes.NewRouterAnswer(answerUseCase, testUseCase, questionUseCase, redisUseCase, authHandler).GetAnswerRouter(router)

	// Apply CORS handler
	handler := c.Handler(router)
	router.HandleFunc("/", homeHandler).Methods("GET")
	// Start the server
	port := "8080"
	log.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))

}

// Hàm xử lý cho route "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the home page!")
}
