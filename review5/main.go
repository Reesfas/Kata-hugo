package main

// Написать сервер, отправляющий запрос на API https://demo.apistubs.io/api/v1/users
// и возвращающий общую структуру пользователя с данными профиля,без staticData
// Если сумма на аккаунте превышает 50000 скрыть email, username, firstName, lastName, avatar
import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type Object struct {
	Records     []Record `json:"records"`
	Skip        int      `json:"skip"`
	Limit       int      `json:"limit"`
	TotalAmount int      `json:"totalAmount"`
}

type Record struct {
	ID        int     `json:"id"`
	Email     string  `json:"email"`
	Profile   Profile `json:"profile"`
	Password  string  `json:"password"`
	Username  string  `json:"username"`
	CreatedAt string  `json:"createdAt"`
	CreatedBy string  `json:"createdBy"`
}

type Profile struct {
	Dob       string `json:"dob"`
	Avatar    string `json:"avatar"`
	LastName  string `json:"lastName"`
	Firstname string `json:"firstname"`
	//StaticData string `json:"staticData"`
}

// Если сумма на аккаунте превышает 50000 скрыть email, username, firstName, lastName, avatar
func main() {
	r := chi.NewRouter()
	r.Get("/users", getUsers)
	fmt.Println("Сервер запущен на порту :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var Obj Object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Отправляем GET-запрос с использованием созданного клиента
	resp, err := client.Get("https://demo.apistubs.io/api/v1/users")
	if err != nil {
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		log.Fatal(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&Obj)
	fmt.Println(Obj)
	if err != nil {
		http.Error(w, "Ошибка при сериализации данных", http.StatusInternalServerError)
		return
	}

	if Obj.TotalAmount < 50000 {
		for i, _ := range Obj.Records {
			Obj.Records[i].Email = ""
			Obj.Records[i].Username = ""
			Obj.Records[i].Profile.Firstname = ""
			Obj.Records[i].Profile.LastName = ""
			Obj.Records[i].Profile.Avatar = ""
		}
	}
	err = json.NewEncoder(w).Encode(Obj)
	if err != nil {
		return
	}
}
