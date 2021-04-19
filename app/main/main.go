package main

import (
	ForumRepository "../../internal/forum/repository"
	ForumUsecase "../../internal/forum/usecase"
	ServiceDelivery "../../internal/service/delivery"
	ServiceRepository "../../internal/service/repository"
	ServiceUsecase "../../internal/service/usecase"
	ThreadDelivery "../../internal/thread/delivery"
	ThreadRepository "../../internal/thread/repository"
	ThreadUsecase "../../internal/thread/usecase"
	UserDelivery "../../internal/user/delivery"
	UserRepository "../../internal/user/repository"
	UserUsecase "../../internal/user/usecase"
	VoteRepository "../../internal/vote/repository"
	VoteUsecase "../../internal/vote/usecase"
	"database/sql"
	"fmt"
	"github.com/NikitaLobaev/BMSTU-DB/config" //TODO: import github...; везде так
	ForumDelivery "github.com/NikitaLobaev/BMSTU-DB/internal/forum/delivery"
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

	forumRepository := ForumRepository.NewForumRepository(dbConnection)
	serviceRepository := ServiceRepository.NewServiceRepository(dbConnection)
	threadRepository := ThreadRepository.NewThreadRepository(dbConnection)
	userRepository := UserRepository.NewUserRepository(dbConnection)
	voteRepository := VoteRepository.NewVoteRepository(dbConnection)

	forumUsecase := ForumUsecase.NewForumUsecase(forumRepository)
	serviceUsecase := ServiceUsecase.NewServiceUsecase(serviceRepository)
	voteUsecase := VoteUsecase.NewVoteUsecase(voteRepository)
	threadUsecase := ThreadUsecase.NewThreadUsecase(threadRepository, voteUsecase)
	userUsecase := UserUsecase.NewUserUsecase(userRepository)

	forumDelivery := ForumDelivery.NewForumHandler(forumUsecase)
	serviceDelivery := ServiceDelivery.NewServiceHandler(serviceUsecase)
	threadDelivery := ThreadDelivery.NewThreadHandler(threadUsecase)
	userDelivery := UserDelivery.NewUserHandler(userUsecase)

	echoWS := echo.New()

	forumDelivery.Configure(echoWS)
	serviceDelivery.Configure(echoWS)
	threadDelivery.Configure(echoWS)
	userDelivery.Configure(echoWS)

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
