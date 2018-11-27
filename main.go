package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var baseUser = []User{}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var Token = "ATw402zK"

type User struct {
	Email    string
	Password string
	BadLogin string
	Image    string
	Errors   map[string]string
}

func main() {
	file, err := os.Open("base_users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, &baseUser)

	http.HandleFunc("/signup", CreateSignUpHandler)
	http.HandleFunc("/validate", ValidateSignUpHandler)
	http.HandleFunc("/login", CreateLoginHandler)
	http.HandleFunc("/login/ok", ValidateLoginHandler)
	http.HandleFunc("/login/users", CreateUsersTabelHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":3030", nil))
}

func Render(w http.ResponseWriter, tmpl string, information interface{}) error {

	template, err := template.ParseFiles(tmpl)
	if err != nil {
		return err
	}

	err = template.Execute(w, information)
	if err != nil {
		return err
	}

	return nil
}

func CreateSignUpHandler(w http.ResponseWriter, r *http.Request) {
	err := Render(w, "static/index.html.tmpl", User{})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (user User) Valid(tok string) bool {

	if !emailRegexp.MatchString(user.Email) {
		user.Errors["Email"] = "Your email is invalid"
	}
	if len(user.Password) > 1 && len(user.Password) < 8 {
		user.Errors["Password"] = "Very short password"
	}
	if tok != Token {
		user.Errors["Token"] = "Token is invalid"
	}
	if user.Password == "" {
		user.Errors["PasswordEmpty"] = "Password Empty"
	}
	return len(user.Errors) == 0
}
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func ValidateSignUpHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	url := "https://picsum.photos/250/250/?image="
	url += strconv.Itoa(random(0, 1000))
	response, err := http.Get(url)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	defer response.Body.Close()

	SignUPToken := r.Form.Get("Token")

	user := User{
		Email:    r.Form.Get("Email"),
		Password: r.Form.Get("Password"),
		Image:    url,
		Errors:   map[string]string{},
	}

	if user.Valid(SignUPToken) {
		for _, savedUser := range baseUser {
			if user.Email != savedUser.Email {
				continue
			}
			user.Errors["EmailExist"] = "Email already exits"
			Render(w, "static/index.html.tmpl", user)
			return
		}

		baseUser = append(baseUser, user)
		baseUserJson, err := json.Marshal(baseUser)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		err = ioutil.WriteFile("base_users.json", baseUserJson, 0644)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	Render(w, "static/index.html.tmpl", user)
}

func CreateLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := Render(w, "static/login.html.tmpl", User{})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func ValidateLoginHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	user := User{
		Email:    r.Form.Get("Email"),
		Password: r.Form.Get("Password"),
	}

	for _, savedUser := range baseUser {
		if user.Email == savedUser.Email && user.Password == savedUser.Password {
			http.Redirect(w, r, "/login/users", http.StatusSeeOther)
			return
		}
	}
	user.BadLogin = "Bad Login"
	Render(w, "static/login.html.tmpl", user)

}

func CreateUsersTabelHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, "static/users.html.tmpl", baseUser)
}
