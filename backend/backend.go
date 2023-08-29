package backend

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	FirstName string
	LastName  string
	email_id  string
	password  string
}

type User struct {
	Username string
	Userid   int
}

type Tasks struct {
	Task_id          int
	Task_name        string
	Task_description string
	Status           string
	StatusOptions    []string
}

type Add_Task struct {
	TaskName        string
	TaskDescription string
}

func (p Person) ToTupleString() string {
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", p.FirstName, p.LastName, p.email_id, p.password)
}

func Login(email string, password string) User {
	dsn := "root:180404@tcp(localhost:3306)/to_do_list"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT userid,FirstName from users where email_id = '%s' and password = '%s'", email, password))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var person User
		err := rows.Scan(&person.Userid, &person.Username)
		if err != nil {
			log.Fatal(err)
		}
		return person
	}

	return User{"Invalid Credentials", -1}
}

func AfterLogin(Userid int) []Tasks {
	StatusnotDone := "not done"
	StatusInProgress := "In Progress"
	StatusDone := "done"

	dsn := "root:180404@tcp(localhost:3306)/to_do_list"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT tasks_info.task_id,task_name,task_description,status from tasks_info left join tasks_status on tasks_info.task_id = tasks_status.task_id where user_id = %d", Userid))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Task_items []Tasks
	for rows.Next() {
		var task Tasks
		err := rows.Scan(&task.Task_id, &task.Task_name, &task.Task_description, &task.Status)
		if err != nil {
			log.Fatal(err)
		}
		task.StatusOptions = append(task.StatusOptions, StatusnotDone, StatusInProgress, StatusDone)
		Task_items = append(Task_items, task)
	}
	return Task_items

}

func GetTaskList() []Add_Task {

	dsn := "root:180404@tcp(localhost:3306)/to_do_list"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT task_name,task_description from tasks_info"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var Task_items []Add_Task
	for rows.Next() {
		var task Add_Task
		err := rows.Scan(&task.TaskName, &task.TaskDescription)
		if err != nil {
			log.Fatal(err)
		}
		Task_items = append(Task_items, task)
	}
	return Task_items

}

func SaveSignup(firstname string, lastname string, username string, password string) {
	dsn := "root:180404@tcp(localhost:3306)/to_do_list"

	person := Person{firstname, lastname, username, password}
	detail := string(person.ToTupleString())
	fmt.Println("(", detail, ")")

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare a SQL statement
	det := "INSERT INTO users(FirstName,LastName,email_id,password)VALUES(" + detail + ")"
	tx, err := db.Begin()

	_, err = tx.Exec(det)
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Values inserted successfully")

}

func AddTask(TaskName string, TaskDescription string, Userid int) {
	dsn := "root:180404@tcp(localhost:3306)/to_do_list"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()

	query1 := fmt.Sprintf("INSERT INTO TASKS_INFO(TASK_NAME,TASK_DESCRIPTION) VALUES('%s','%s')", TaskName, TaskDescription)

	rows, err := db.Query(fmt.Sprintf("SELECT TASK_ID FROM TASKS_INFO WHERE TASK_NAME = '%s'", TaskName))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	TID := -1
	for rows.Next() {
		err := rows.Scan(&TID)
		if err != nil {
			log.Fatal(err)
		}
	}
	if TID == -1 {
		_, err = tx.Exec(query1)
		if err != nil {
			log.Fatal(err)
		}

		rows, err := tx.Query(fmt.Sprintf("SELECT TASK_ID FROM TASKS_INFO WHERE TASK_NAME = '%s'", TaskName))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&TID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	query2 := fmt.Sprintf("INSERT INTO TASKS_STATUS(TASK_ID,USER_ID,STATUS) VALUES(%d,%d,'not done')", TID, Userid)

	_, err = tx.Exec(query2)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()

	fmt.Println("Values inserted successfully")
}

func Senddata(ReceivedData string, Userid int, TaskId string) {
	dsn := "root:180404@tcp(localhost:3306)/to_do_list"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()

	query := fmt.Sprintf("UPDATE TASKS_STATUS SET STATUS = '%s' WHERE USER_ID = %d and task_id = %s", ReceivedData, Userid, TaskId)

	_, err = tx.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
