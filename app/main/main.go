package main

import (
	"../../config" //TODO: import github...; везде так
	forumDelivery "../../internal/forum/delivery"
	forumRepository "../../internal/forum/repository"
	forumUsecase "../../internal/forum/usecase"
	userDelivery "../../internal/user/delivery"
	userRepository "../../internal/user/repository"
	userUsecase "../../internal/user/usecase"
	"database/sql"
	"fmt"
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

	forumRepository_ := forumRepository.NewForumRepository(dbConnection)
	userRepository_ := userRepository.NewUserRepository(dbConnection)

	forumUsecase_ := forumUsecase.NewForumUsecase(forumRepository_)
	userUsecase_ := userUsecase.NewUserUsecase(userRepository_)

	forumDelivery_ := forumDelivery.NewForumHandler(forumUsecase_)
	userDelivery_ := userDelivery.NewUserHandler(userUsecase_)

	echoWS := echo.New()

	forumDelivery_.Configure(echoWS)
	userDelivery_.Configure(echoWS)

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
