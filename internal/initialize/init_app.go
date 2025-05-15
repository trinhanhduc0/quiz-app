package initialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	entity "quiz-app/internal/domain/entities"
	"quiz-app/internal/domain/service"
	"quiz-app/internal/infrastructure/persistence/aws"
	persistence "quiz-app/internal/infrastructure/persistence/mongodb"
	redisdb "quiz-app/internal/infrastructure/persistence/redis"
	routes "quiz-app/internal/infrastructure/router"
	"time"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	dbSevice = "mongodb"
	dbName   = "dbapp"
)

func InitApp() {
	InitDB()
	//Init router
	InitRouter()

}

func createTestData() {
	types := []string{"multiple_choice_question", "match_choice_question", "fill_in_the_blank", "single_choice_question", "order_question"}
	apiURL := "http://localhost:8080/questions" // 🔁 Thay đổi thành endpoint của bạn
	authorID := "ZG1mwdlEFtRwxXRezbmgf6Ctij13"
	fmt.Println(apiURL)
	for i := 0; i < 1000; i++ {
		qType := types[rand.Intn(len(types))]
		var question entity.Question
		question.Type = qType
		question.Metadata = entity.Metadata{Author: authorID}
		question.Created_At = time.Now()
		question.Updated_At = time.Now()
		question.Score = float32(rand.Intn(10) + 1)
		question.Tags = []string{"tag1", "tag2"}
		question.QuestionContent = entity.QuestionContent{
			Text:     fmt.Sprintf("Câu hỏi %d (%s)", i+1, qType),
			ImageURL: "",
		}

		switch qType {
		case "multiple_choice_question", "single_choice_question":
			for j := 0; j < 3; j++ {
				question.Options = append(question.Options, entity.Option{
					ID:        primitive.NewObjectID(),
					Text:      fmt.Sprintf("Lựa chọn %d", j+1),
					IsCorrect: j == 0,
					ImageURL:  "",
				})
			}
		case "match_choice_question":
			for j := 0; j < 3; j++ {
				question.MatchItems = append(question.MatchItems, entity.MatchItem{
					ID:   primitive.NewObjectID(),
					Text: fmt.Sprintf("Ghép %d", j+1),
				})
				question.MatchOptions = append(question.MatchOptions, entity.MatchOption{
					ID:   primitive.NewObjectID(),
					Text: fmt.Sprintf("Đáp án %d", j+1),
				})
			}
		case "fill_in_the_blank":
			question.FillInTheBlanks = []entity.FillInTheBlank{
				{
					ID:            primitive.NewObjectID(),
					TextBefore:    "Tôi thích",
					Blank:         "___",
					CorrectAnswer: "bơi",
					TextAfter:     "vào cuối tuần.",
				},
			}
		case "order_question":
			for j := 0; j < 3; j++ {
				question.OrderItems = append(question.OrderItems, entity.OrderItem{
					ID:    primitive.NewObjectID(),
					Text:  fmt.Sprintf("Bước %d", j+1),
					Order: j + 1,
				})
			}
		}

		// Marshal & POST
		body, _ := json.Marshal(question)
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("❌ Failed to post:", err)
			continue
		}
		defer resp.Body.Close()
		fmt.Printf("✅ (%d) %s\n", i+1, resp.Status)
	}
}

func InitDB() {
	//Init mongo
	switch dbSevice {
	case "mongodb":
		{
			persistence.ConnectMongoDB(os.Getenv(persistence.MongoConnectionString))
		}
	default:

	}
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
	//authService, err := service.NewAuthService("firebase-config.json", ([]byte)("jwt-secret"))

	authService, err := service.NewAuthService("firebase-config.json", ([]byte)(os.Getenv("JWT_SECRET")))
	fmt.Println(authService)
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

	go routes.NewRoutesAuth(router, *authService, *userUseCase, *redisUseCase, *classUseCase).SetLoginRoute()
	routes.NewRouterTest(*testUseCase, *classUseCase, *questionUseCase, *answerUseCase, *redisUseCase, *authHandler).GetTestRouter(router)
	routes.NewRouterQuestion(*questionUseCase, *authHandler).GetQuestionRouter(router)
	routes.NewRouterClass(*classUseCase, *redisUseCase, *authHandler).GetClassRouter(router)
	routes.NewRoutesFile(fileUseCase, awsS3UseCase, authHandler).GetRoutesFile(router)
	// routes.NewRouterAnswer(answerUseCase, testUseCase, questionUseCase, redisUseCase, authHandler).GetAnswerRouter(router)

	// Apply CORS handler
	handler := c.Handler(router)
	router.HandleFunc("/", homeHandler).Methods("GET")
	// Start the server
	port := "8080"
	log.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
	createTestData()

}

// Hàm xử lý cho route "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the home page!")
}
