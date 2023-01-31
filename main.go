package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

var (
	errCreateTestDb = errors.New("Failed to create db driver")
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/app", handler)
	if err := http.ListenAndServe(":5005", mux); err != nil {
		log.Fatalf("Error to start server")
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := RedisExecution(); err != nil {
		log.Fatal("error in redis execution", err)
	}
	if err := MysqlExecution(); err != nil {
		log.Fatal("error in mysql execution", err)
	}
	fmt.Fprintf(w, "Hello World")
	fmt.Println("In app")
}

func RedisExecution() error {

	client := redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Failed to ping redis server")
		return err
	}

	fmt.Println(pong, err)
	return nil
}

func MysqlExecution() error {

	db, err := sql.Open("mysql", "root:whatever@tcp(mysql-server:3306)/test")
	if err != nil {
		fmt.Print(err)
		return fmt.Errorf("error in opening sql driver: %w", errCreateTestDb)

	}

	_, err = db.Query("Create table employee(id int, value varchar(40))")
	if err != nil {
		return fmt.Errorf("error in create table employee: %w", err)
	}

	return nil
}
