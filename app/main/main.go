package main

import (
	"database/sql"
	"fmt"
	"github.com/NikitaLobaev/BMSTU-DB/config"
	ForumDelivery "github.com/NikitaLobaev/BMSTU-DB/internal/forum/delivery"
	ForumRepository "github.com/NikitaLobaev/BMSTU-DB/internal/forum/repository"
	ForumRepositoryV1 "github.com/NikitaLobaev/BMSTU-DB/internal/forum/repository/v1"
	ForumRepositoryV2 "github.com/NikitaLobaev/BMSTU-DB/internal/forum/repository/v2"
	ForumUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/forum/usecase"
	PostDelivery "github.com/NikitaLobaev/BMSTU-DB/internal/post/delivery"
	PostRepository "github.com/NikitaLobaev/BMSTU-DB/internal/post/repository"
	PostRepositoryV1 "github.com/NikitaLobaev/BMSTU-DB/internal/post/repository/v1"
	PostRepositoryV2 "github.com/NikitaLobaev/BMSTU-DB/internal/post/repository/v2"
	PostUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/post/usecase"
	ServiceDelivery "github.com/NikitaLobaev/BMSTU-DB/internal/service/delivery"
	ServiceRepository "github.com/NikitaLobaev/BMSTU-DB/internal/service/repository"
	ServiceRepositoryV1 "github.com/NikitaLobaev/BMSTU-DB/internal/service/repository/v1"
	ServiceRepositoryV2 "github.com/NikitaLobaev/BMSTU-DB/internal/service/repository/v2"
	ServiceUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/service/usecase"
	ThreadDelivery "github.com/NikitaLobaev/BMSTU-DB/internal/thread/delivery"
	ThreadRepository "github.com/NikitaLobaev/BMSTU-DB/internal/thread/repository"
	ThreadRepositoryV1 "github.com/NikitaLobaev/BMSTU-DB/internal/thread/repository/v1"
	ThreadRepositoryV2 "github.com/NikitaLobaev/BMSTU-DB/internal/thread/repository/v2"
	ThreadUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/thread/usecase"
	UserDelivery "github.com/NikitaLobaev/BMSTU-DB/internal/user/delivery"
	UserRepository "github.com/NikitaLobaev/BMSTU-DB/internal/user/repository"
	UserRepositoryV1 "github.com/NikitaLobaev/BMSTU-DB/internal/user/repository/v1"
	UserRepositoryV2 "github.com/NikitaLobaev/BMSTU-DB/internal/user/repository/v2"
	UserUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/user/usecase"
	VoteRepository "github.com/NikitaLobaev/BMSTU-DB/internal/vote/repository"
	VoteRepositoryV1 "github.com/NikitaLobaev/BMSTU-DB/internal/vote/repository/v1"
	VoteRepositoryV2 "github.com/NikitaLobaev/BMSTU-DB/internal/vote/repository/v2"
	VoteUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/vote/usecase"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	argsCount = 2
	argsUsage = "Usage: go run main.go <config file>"
)

func main() {
	if len(os.Args) != argsCount {
		fmt.Println(argsUsage)
		return
	}

	config_, err := config.NewConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	dbConfigString := genDBConfigString(config_)
	dbConnection, err := sql.Open("postgres", dbConfigString)
	defer func() {
		if err := dbConnection.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	if err := dbConnection.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database config:", dbConfigString)

	var forumRepository ForumRepository.ForumRepository
	var serviceRepository ServiceRepository.ServiceRepository
	var threadRepository ThreadRepository.ThreadRepository
	var userRepository UserRepository.UserRepository
	var voteRepository VoteRepository.VoteRepository
	var postRepository PostRepository.PostRepository
	switch config_.Database.DBName {
	case "forums_2nf":
		forumRepository = ForumRepositoryV2.NewForumRepositoryV2(dbConnection)
		postRepository = PostRepositoryV2.NewPostRepositoryV2(dbConnection)
		serviceRepository = ServiceRepositoryV2.NewServiceRepositoryV2(dbConnection)
		threadRepository = ThreadRepositoryV2.NewThreadRepositoryV2(dbConnection)
		userRepository = UserRepositoryV2.NewUserRepositoryV2(dbConnection)
		voteRepository = VoteRepositoryV2.NewVoteRepositoryV2(dbConnection)
		break
	case "forums_3nf":
		break
	default: //"forums_1nf"
		forumRepository = ForumRepositoryV1.NewForumRepositoryV1(dbConnection)
		postRepository = PostRepositoryV1.NewPostRepositoryV1(dbConnection)
		serviceRepository = ServiceRepositoryV1.NewServiceRepositoryV1(dbConnection)
		threadRepository = ThreadRepositoryV1.NewThreadRepositoryV1(dbConnection)
		userRepository = UserRepositoryV1.NewUserRepositoryV1(dbConnection)
		voteRepository = VoteRepositoryV1.NewVoteRepositoryV1(dbConnection)
		break
	}

	serviceUsecase := ServiceUsecase.NewServiceUsecase(serviceRepository)
	voteUsecase := VoteUsecase.NewVoteUsecase(voteRepository)
	threadUsecase := ThreadUsecase.NewThreadUsecase(threadRepository, voteUsecase)
	userUsecase := UserUsecase.NewUserUsecase(userRepository)
	forumUsecase := ForumUsecase.NewForumUsecase(forumRepository, threadUsecase, userUsecase)
	postUsecase := PostUsecase.NewPostUsecase(postRepository, forumUsecase, threadUsecase, userUsecase)

	forumDelivery := ForumDelivery.NewForumHandler(forumUsecase)
	serviceDelivery := ServiceDelivery.NewServiceHandler(serviceUsecase)
	threadDelivery := ThreadDelivery.NewThreadHandler(threadUsecase)
	userDelivery := UserDelivery.NewUserHandler(userUsecase)
	postDelivery := PostDelivery.NewPostHandler(postUsecase)

	echoWS := echo.New()

	forumDelivery.Configure(echoWS)
	serviceDelivery.Configure(echoWS)
	threadDelivery.Configure(echoWS)
	userDelivery.Configure(echoWS)
	postDelivery.Configure(echoWS)

	if err := echoWS.Start(genWSConfigString(config_)); err != nil {
		log.Fatal(err)
	}
}

func genDBConfigString(config *config.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host,
		config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)
}

func genWSConfigString(config *config.Config) string {
	return fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
}
