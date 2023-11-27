package handler

import (
	"31_5/pkg/color"         //Раскрашиваем текст
	"31_5/pkg/repo/userRepo" //Репозиторий для хранилища пользователей
	"31_5/pkg/user"          //Создание пользователей
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	msg     string
	appName string
	u       user.User
	storage userRepo.Storage
)

type Service struct {
	Store   map[string]*user.User
	AppName string //Имя приложения, для отслеживания отправителя
}

// Функция-приветствие. Отвечает на handler "/"
func Greeting(w http.ResponseWriter, r *http.Request) {
	//Ответ на запрос
	msg = "I'm alive!\n"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))

	/* Просмотр тела запроса
	content, _ := io.ReadAll(r.Body)
	fmt.Println(string(content))
	*/
}

// 1. Метод для типа service - создание пользователя. Парсит входящий запрос в формате JSON, создает пользователя с полями:
// name (имя, string), age (возраст, int), friends (список друзей, []string)
func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		s.Store[u.Name] = &u
		msg = s.AppName + "Пользователь " + color.Green() + u.Name + color.Rst() + " создан\n"
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(msg))

		//Парсинг json-файла
		rawDataIn, err := os.ReadFile("users.json")
		if err != nil {
			log.Fatal("Не могу открыть хранилище:", err)
		}

		err = json.Unmarshal(rawDataIn, &storage)
		if err != nil {
			log.Fatal("Не правильный формат файла:", err)
		}

		newUser := user.UpdateUser(u.Name, u.Age, u.Friends)
		storage.Users = append(storage.Users, newUser)

		/*
			for _, x := range storage.Users {
				fmt.Println(x.Name)
			}
		*/
		rawDataOut, err := json.MarshalIndent(&storage, "", "  ")
		if err != nil {
			log.Fatal("Парсинг JSON не удался:", err)
		}

		//Запись (добавление) созданного пользователя в файл
		err = os.WriteFile("users.json", rawDataOut, 0)
		if err != nil {
			log.Fatal("Не получилось записать в файл:", err)
		}
		return

	}
	w.WriteHeader(http.StatusBadRequest)
}

// 2. Создание друзей для пользователя. Обработчик, который делает друзей из двух пользователей.
// Например, если мы создали двух пользователей и нам вернулись их ID, то в запросе мы можем указать ID пользователя, который инициировал запрос на дружбу,
// и ID пользователя, который примет инициатора в друзья.
func (s *Service) MakeFriendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		type friendRequest struct {
			SourceId string `json:"source_id"`
			TargetId string `json:"target_id"`
		}

		var fr friendRequest

		if err := json.Unmarshal(content, &fr); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//Парсинг json-файла
		rawDataIn, err := os.ReadFile("users.json")
		if err != nil {
			log.Fatal("Не могу открыть хранилище:", err)
		}

		err = json.Unmarshal(rawDataIn, &storage)
		if err != nil {
			log.Fatal("Неправильный формат файла:", err)
		}

		var newStorage userRepo.Storage //Нужно перезаписать хранилище с пользователями, т.к. у кого-то новые друзья

		ok := 0
		//Пробежаться по списку пользователей и найти пользователя (доработать - если не один такой?)
		for _, users := range storage.Users {
			if users.Name == fr.SourceId {
				users.Friends = append(users.Friends, fr.TargetId)
				ok++
			}
			newStorage.Users = append(newStorage.Users, users)
		}
		if ok == 0 { //Если не нашли пользователя
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Errorf("Пользователь %s не найден\n", fr.SourceId).Error()))
			return
		} else { //Если есть такой пользователь, сообщаем о добавлении друга
			msg = s.AppName + "Пользователь " + color.Green() + fr.SourceId + color.Rst() + " имеет нового друга " + color.Yellow() + fr.TargetId + color.Rst() + "\n"
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(msg))
		}

		rawDataOut, err := json.MarshalIndent(&newStorage, "", "  ")
		if err != nil {
			log.Fatal("Парсинг в JSON не удался:", err)
		}

		//Запись (добавление) созданного пользователя в файл
		err = os.WriteFile("users.json", rawDataOut, 0)
		if err != nil {
			log.Fatal("Не получилось записать в файл:", err)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// 3. Обработчик, который удаляет пользователя. Принимает ID пользователя и удаляет его из хранилища, а также стирает его из массива friends у всех его друзей.
// Запрос возвращает 200 и имя удалённого пользователя
func (s *Service) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		//Парсинг json-файла
		rawDataIn, err := os.ReadFile("users.json")
		if err != nil {
			log.Fatal("Не могу открыть хранилище:", err)
		}

		err = json.Unmarshal(rawDataIn, &storage)
		if err != nil {
			log.Fatal("Неправильный формат файла:", err)
		}

		type DeleteRequest struct {
			TargetId string `json:"target_id"`
		}

		var dr DeleteRequest
		if err := json.Unmarshal(content, &dr); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		ok := 0
		var newStorage userRepo.Storage //Нужно перезаписать хранилище с пользователями, т.к. удалили пользователя
		for _, user := range storage.Users {
			if user.Name != dr.TargetId {
				newStorage.Users = append(newStorage.Users, user)
			} else {
				ok++
			}
			for j, friend := range user.Friends {
				if friend == dr.TargetId {
					user.Friends = append(user.Friends[:j], user.Friends[j+1:]...)
				}
			}
		}
		if ok == 0 { // если пользователя нет в хранилище - нужно об этом сказать
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Errorf("%s Пользователь %s не найден\n", s.AppName, dr.TargetId).Error()))
			return
		} else { // если есть ok > 0, сообщаем об удалении пользователя
			w.WriteHeader(http.StatusOK)
			msg = s.AppName + "Пользователь " + color.Green() + dr.TargetId + color.Rst() + color.Red() + " удален" + color.Rst() + "\n"
			w.Write([]byte(msg))
		}

		delete(s.Store, dr.TargetId)

		rawDataOut, err := json.MarshalIndent(&newStorage, "", "  ")
		if err != nil {
			log.Fatal("Парсинг в JSON не удался:", err)
		}

		//Запись (добавление) созданного пользователя в файл
		err = os.WriteFile("users.json", rawDataOut, 0)
		if err != nil {
			log.Fatal("Не получилось записать в файл:", err)
		}

		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// 4. Обработчик, который возвращает всех друзей пользователя.
// После /friends/ указывается id пользователя, друзей которого мы хотим увидеть.
func (s *Service) GetFriendsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		requestPath := r.URL.Path
		_, requestFile := path.Split(requestPath)
		userID := requestFile

		//Парсинг json-файла
		rawDataIn, err := os.ReadFile("users.json")
		if err != nil {
			log.Fatal("Не могу открыть хранилище:", err)
		}

		err = json.Unmarshal(rawDataIn, &storage)
		if err != nil {
			log.Fatal("Неправильный формат файла:", err)
		}

		ok := 0
		for _, user := range storage.Users {
			if user.Name == userID {
				response := ""
				for _, friend := range user.Friends {
					response += friend + "\n"
				}
				msg = "Друзья пользователя " + color.Green() + userID + color.Rst() + ":\n" + response
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(msg))
				ok++
			}
		}
		if ok == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Errorf("%s Пользователь %s не найден", s.AppName, userID).Error()))
			return
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// 5. Обработчик, который обновляет возраст пользователя.
// Запрос должен возвращать 200 и сообщение «возраст пользователя успешно обновлён».
func (s *Service) NewAgeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL)

	if r.Method == "PUT" {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		//Парсинг json-файла
		rawDataIn, err := os.ReadFile("users.json")
		if err != nil {
			log.Fatal("Не могу открыть хранилище:", err)
		}

		err = json.Unmarshal(rawDataIn, &storage)
		if err != nil {
			log.Fatal("Неправильный формат файла:", err)
		}

		type AgeUpdateRequest struct {
			NewAge int `json:"new_age"`
		}

		var ur AgeUpdateRequest
		if err := json.Unmarshal(content, &ur); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Println(ur.NewAge)

		requestPath := r.URL.Path
		_, requestFile := path.Split(requestPath)
		userID := requestFile

		ok := 0
		var newStorage userRepo.Storage
		for _, user := range storage.Users {
			if user.Name == userID {
				user.Age = ur.NewAge
				w.WriteHeader(http.StatusOK)
				msg = "Возраст пользователя " + color.Green() + userID + color.Rst() + " обновлен " + color.Green() + "успешно" + color.Rst() + "!\n"
				w.Write([]byte(msg))
				ok++
			}
			newStorage.Users = append(newStorage.Users, user)
		}
		if ok == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Errorf("Пользователь %s не найден", userID).Error()))
			return
		}

		rawDataOut, err := json.MarshalIndent(&newStorage, "", "  ")
		if err != nil {
			log.Fatal("Парсинг в JSON не удался:", err)
		}

		//Запись (добавление) созданного пользователя в файл
		err = os.WriteFile("users.json", rawDataOut, 0)
		if err != nil {
			log.Fatal("Не получилось записать в файл:", err)
		}

		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
