package main

import (
	"backend"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type DashboardView struct {
	User               backend.User
	Task_items         []backend.Tasks
	Existing_Task_List []backend.Add_Task
}

type DataFromClient struct {
	SelectedOption string
	TDValues       []string
}

type LoginView struct {
	ErrMessage string
}

var task_list = []backend.Tasks{}

var loggedInUser backend.User
var ErrMessage string

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the user is logged in, show the dashboard
		if loggedInUser.Username != "" {
			tmpl := template.Must(template.ParseFiles("static/dashboard.html"))
			err := tmpl.Execute(w, DashboardView{loggedInUser, backend.AfterLogin(loggedInUser.Userid), backend.GetTaskList()})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			tmp2 := template.Must(template.ParseFiles("static/login.html"))
			err := tmp2.Execute(w, LoginView{ErrMessage})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		loggedInUser = backend.User{}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/signup.html")
	})

	http.HandleFunc("/saveSignup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fullname := r.FormValue("fullname")
			email := r.FormValue("email")
			password := r.FormValue("password")
			backend.SaveSignup(fullname, fullname, email, password)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/checklogin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")
			userCred := backend.Login(email, password)
			if userCred.Userid != -1 {
				loggedInUser = userCred
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/addTask", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			Task_name := r.FormValue("TaskName")
			Task_des := r.FormValue("TaskDescription")
			backend.AddTask(Task_name, Task_des, loggedInUser.Userid)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/SendData", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}

			var receivedData DataFromClient
			err = json.Unmarshal(body, &receivedData)
			if err != nil {
				http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
				return
			}

			fmt.Println("Received selected option:", receivedData.SelectedOption)
			fmt.Println("Received TD values:", receivedData.TDValues)
			err = json.Unmarshal(body, &receivedData)
			backend.Senddata(receivedData.SelectedOption, loggedInUser.Userid, receivedData.TDValues[0])

			// Process and use the received data as needed

			// Send a response back to the client if needed
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Data received successfully"))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/DelTask", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}

			var receivedData DataFromClient
			err = json.Unmarshal(body, &receivedData)
			if err != nil {
				http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
				return
			}

			fmt.Println("Received TD values:", receivedData.TDValues)
			err = json.Unmarshal(body, &receivedData)
			backend.Deldata(loggedInUser.Userid, receivedData.TDValues[0])

			// Process and use the received data as needed

			// Send a response back to the client if needed
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Data deleted successfully"))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.ListenAndServe(":8080", nil)

}
