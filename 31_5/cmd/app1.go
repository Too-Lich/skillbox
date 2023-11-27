package main

import (
	"31_5/pkg/color"
	"31_5/pkg/handler"
	"31_5/pkg/user"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

/*
HTTP-сервис, который принимает входящие соединения с JSON-данными и обрабатывает их следующим образом:

1. Cоздание пользователя с полями: имя, возраст и массив друзей. Пример запроса: {"name":"some name","age":"24","friends":[]} localhost:9000/create

2. Создание друзей для пользователя. Обработчик, который делает друзей из двух пользователей. Пример запроса: {"source_id": "Игорь", "target_id": "Иван"}

3. Удаление пользователя. Принимает ID пользователя и удаляет его из хранилища, а также стирает его из массива friends у всех его друзей. DELETE -d '{"target_id": "Петр"}' localhost:8080/user

4. Возвращает всех друзей пользователя. После /friends/ указывается id пользователя, друзей которого мы хотим увидеть. GET http://localhost:8080/friends/Игорь

5. Обновление возраста пользователя. Пример запроса {"new age": 28} localhost:9000/Игорь
*/
func main() {
	fmt.Println(color.Green() + "<<-- App1. HTTP-сервис, принимающий данные JSON -->>" + color.Rst())

	srv := handler.Service{
		Store:   make(map[string]*user.User),
		AppName: "app_1: ",
	}

	var r = chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", handler.Greeting)                 //Приветствие
	r.Post("/create", srv.CreateUser)             //Создание пользователя
	r.Post("/makefriend", srv.MakeFriendHandler)  //Подружить пользователей
	r.Delete("/user", srv.DeleteUserHandler)      //Удалить пользователя
	r.Get("/friends/{id}", srv.GetFriendsHandler) //Посмотреть список друзей пользователя
	r.Put("/{id}", srv.NewAgeHandler)             //Обновить возраст пользователя

	http.ListenAndServe("localhost:8080", r)
}
