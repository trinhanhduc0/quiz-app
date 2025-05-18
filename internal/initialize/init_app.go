package initialize

import (
	"bytes"
	"encoding/json"
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

func SeedSampleQuestions(author string) []*entity.Question {
	now := time.Now()

	return []*entity.Question{
		// 1. SINGLE CHOICE
		{
			ID:   primitive.NewObjectID(),
			Type: "multiple_choice_single",
			QuestionContent: entity.QuestionContent{
				Text: "What is the capital of France?",
			},
			Options: []entity.Option{
				{ID: primitive.NewObjectID(), Text: "Berlin"},
				{ID: primitive.NewObjectID(), Text: "Madrid"},
				{ID: primitive.NewObjectID(), Text: "Paris", IsCorrect: true},
				{ID: primitive.NewObjectID(), Text: "Rome"},
			},
			Metadata:   entity.Metadata{Author: author},
			Tags:       []string{"geography", "easy"},
			Suggestion: []string{"It's a European country."},
			Score:      1,
			Created_At: now,
			Updated_At: now,
		},

		// 2. MULTIPLE CHOICE
		{
			ID:   primitive.NewObjectID(),
			Type: "multiple_choice_multiple",
			QuestionContent: entity.QuestionContent{
				Text: "Which of the following are programming languages?",
			},
			Options: []entity.Option{
				{ID: primitive.NewObjectID(), Text: "Python", IsCorrect: true},
				{ID: primitive.NewObjectID(), Text: "HTML"},
				{ID: primitive.NewObjectID(), Text: "Go", IsCorrect: true},
				{ID: primitive.NewObjectID(), Text: "CSS"},
			},
			Metadata:   entity.Metadata{Author: author},
			Tags:       []string{"programming", "logic"},
			Suggestion: []string{"Focus on languages used to build logic, not just styling."},
			Score:      1.5,
			Created_At: now,
			Updated_At: now,
		},

		// 3. FILL IN THE BLANK
		{
			ID:   primitive.NewObjectID(),
			Type: "fill_in_the_blank",
			QuestionContent: entity.QuestionContent{
				Text: "Complete the sentence correctly.",
			},
			FillInTheBlanks: []entity.FillInTheBlank{
				{
					ID:            primitive.NewObjectID(),
					TextBefore:    "The sun rises in the",
					Blank:         "___",
					CorrectAnswer: "east",
					TextAfter:     ".",
				},
			},
			Metadata:   entity.Metadata{Author: author},
			Tags:       []string{"science", "basic"},
			Suggestion: []string{"Think of direction."},
			Score:      1,
			Created_At: now,
			Updated_At: now,
		},

		// 4. ORDERING QUESTION
		{
			ID:   primitive.NewObjectID(),
			Type: "order_question",
			QuestionContent: entity.QuestionContent{
				Text: "Arrange the steps to make a cup of tea.",
			},
			OrderItems: []entity.OrderItem{
				{ID: primitive.NewObjectID(), Text: "Boil water", Order: 1},
				{ID: primitive.NewObjectID(), Text: "Put tea bag in cup", Order: 2},
				{ID: primitive.NewObjectID(), Text: "Pour water", Order: 3},
				{ID: primitive.NewObjectID(), Text: "Add sugar", Order: 4},
			},
			Metadata:   entity.Metadata{Author: author},
			Tags:       []string{"daily life", "sequence"},
			Suggestion: []string{"Start with heating water."},
			Score:      2,
			Created_At: now,
			Updated_At: now,
		},

		// 5. MATCH CHOICE
		{
			ID:   primitive.NewObjectID(),
			Type: "match_choice_question",
			QuestionContent: entity.QuestionContent{
				Text: "Match the animal with the sound it makes.",
			},
			MatchItems: []entity.MatchItem{
				{ID: primitive.NewObjectID(), Text: "Dog"},
				{ID: primitive.NewObjectID(), Text: "Cat"},
				{ID: primitive.NewObjectID(), Text: "Cow"},
			},
			MatchOptions: []entity.MatchOption{
				{ID: primitive.NewObjectID(), Text: "Bark", MatchId: "Dog"},
				{ID: primitive.NewObjectID(), Text: "Meow", MatchId: "Cat"},
				{ID: primitive.NewObjectID(), Text: "Moo", MatchId: "Cow"},
			},
			Metadata:   entity.Metadata{Author: author},
			Tags:       []string{"animals", "kids"},
			Suggestion: []string{"Think of farm animals."},
			Score:      2,
			Created_At: now,
			Updated_At: now,
		},
	}
}
func createTestData() {
	apiURL := "http://localhost:8080/questions" // 🔁 Thay đổi thành endpoint của bạn
	authorID := "ZG1mwdlEFtRwxXRezbmgf6Ctij13"
	fmt.Println(authorID)
	questions := SeedSampleQuestions("qnce02@gmail.com")
	// Marshal & POST
	for i, value := range questions {
		body, _ := json.Marshal(value)
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("❌ Failed to post:", err)
			break
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

}

// Hàm xử lý cho route "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the home page!")
	createTestData()
}
