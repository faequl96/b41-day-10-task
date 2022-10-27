package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	connection.DatabaseConnect()
	route := mux.NewRouter()

	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/register", registerForm).Methods("GET")
	route.HandleFunc("/login", loginForm).Methods("GET")

	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")

	route.HandleFunc("/form-add-project", addProjectForm).Methods("GET")
	route.HandleFunc("/send-data-add-project", sendDataAddProject).Methods("POST")

	route.HandleFunc("/form-edit-project/{id}", editProjectForm).Methods("GET")
	route.HandleFunc("/send-data-edit-project/{id}", sendDataEditProject).Methods("POST")

	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	fmt.Println("Server running on localhost:8000")
	http.ListenAndServe("localhost:8000", route)
}

type projectDataStruc struct {
	Id              int
	ProjectName     string
	StartDate       time.Time
	EndDate         time.Time
	StartDateFormat string
	EndDateFormat   string
	Duration        string
	Description     string
	Technologies    []string
	Image           string
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/home.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	} else {
		var projectData []projectDataStruc
		allDataFrom_db_project, _ := connection.Conn.Query(context.Background(), "SELECT id, project_name, start_date, end_date, duration, description, technologies, image FROM db_project")
		for allDataFrom_db_project.Next() {
			selectedProjectData := projectDataStruc{}
			err := allDataFrom_db_project.Scan(&selectedProjectData.Id, &selectedProjectData.ProjectName, &selectedProjectData.StartDate, &selectedProjectData.EndDate, &selectedProjectData.Duration, &selectedProjectData.Description, &selectedProjectData.Technologies, &selectedProjectData.Image)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			projectData = append(projectData, selectedProjectData)
		}

		response := map[string]interface{}{
			"ProjectData": projectData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/login.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
	}
}

func registerForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/register.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
	}
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	} else {
		where_id, _ := strconv.Atoi(mux.Vars(r)["id"])

		selectedProjectData := projectDataStruc{}

		err = connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, duration, description, technologies, image FROM db_project WHERE id=$1", where_id).
			Scan(&selectedProjectData.Id, &selectedProjectData.ProjectName, &selectedProjectData.StartDate, &selectedProjectData.EndDate, &selectedProjectData.Duration, &selectedProjectData.Description, &selectedProjectData.Technologies, &selectedProjectData.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		selectedProjectData.StartDateFormat = selectedProjectData.StartDate.Format("2006-01-02")
		selectedProjectData.EndDateFormat = selectedProjectData.EndDate.Format("2006-01-02")

		response := map[string]interface{}{
			"selectedProjectData": selectedProjectData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func addProjectForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
	}
}

func sendDataAddProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		projectName := r.PostForm.Get("project-name")
		startDate := r.PostForm.Get("start-date")
		endDate := r.PostForm.Get("end-date")
		var duration string
		description := r.PostForm.Get("description")
		technologies := []string{r.PostForm.Get("node"), r.PostForm.Get("react"), r.PostForm.Get("vue"), r.PostForm.Get("typescript")}
		image := r.PostForm.Get("project-image")

		layoutDate := "2006-01-02"
		startDateParse, _ := time.Parse(layoutDate, startDate)
		endDateParse, _ := time.Parse(layoutDate, endDate)

		hour := 1
		day := hour * 24
		week := hour * 24 * 7
		month := hour * 24 * 30
		year := hour * 24 * 365

		differHour := endDateParse.Sub(startDateParse).Hours()
		var differHours int = int(differHour)
		// fmt.Println(differHours)
		days := differHours / day
		weeks := differHours / week
		months := differHours / month
		years := differHours / year

		if differHours < week {
			duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < month {
			duration = strconv.Itoa(int(weeks)) + " Weeks"
		} else if differHours < year {
			duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			duration = strconv.Itoa(int(years)) + " Years"
		}

		_, err = connection.Conn.Exec(context.Background(), "INSERT INTO db_project(project_name, start_date, end_date, duration, description, technologies, image) VALUES ($1, $2, $3, $4, $5, $6, $7)", projectName, startDate, endDate, duration, description, technologies, image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func editProjectForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/edit-project.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		where_id, _ := strconv.Atoi(mux.Vars(r)["id"])

		selectedProjectData := projectDataStruc{}

		err = connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, duration, description, technologies, image FROM db_project WHERE id=$1", where_id).
			Scan(&selectedProjectData.Id, &selectedProjectData.ProjectName, &selectedProjectData.StartDate, &selectedProjectData.EndDate, &selectedProjectData.Duration, &selectedProjectData.Description, &selectedProjectData.Technologies, &selectedProjectData.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		selectedProjectData.StartDateFormat = selectedProjectData.StartDate.Format("2006-01-02")
		selectedProjectData.EndDateFormat = selectedProjectData.EndDate.Format("2006-01-02")

		response := map[string]interface{}{
			"selectedProjectData": selectedProjectData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func sendDataEditProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		where_id, _ := strconv.Atoi(mux.Vars(r)["id"])

		projectName := r.PostForm.Get("project-name")
		startDate := r.PostForm.Get("start-date")
		endDate := r.PostForm.Get("end-date")
		var duration string
		description := r.PostForm.Get("description")
		technologies := []string{r.PostForm.Get("node"), r.PostForm.Get("react"), r.PostForm.Get("vue"), r.PostForm.Get("typescript")}
		image := r.PostForm.Get("project-image")

		layoutDate := "2006-01-02"
		startDateParse, _ := time.Parse(layoutDate, startDate)
		endDateParse, _ := time.Parse(layoutDate, endDate)

		hour := 1
		day := hour * 24
		week := hour * 24 * 7
		month := hour * 24 * 30
		year := hour * 24 * 365

		differHour := endDateParse.Sub(startDateParse).Hours()
		var differHours int = int(differHour)
		// fmt.Println(differHours)
		days := differHours / day
		weeks := differHours / week
		months := differHours / month
		years := differHours / year

		if differHours < week {
			duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < month {
			duration = strconv.Itoa(int(weeks)) + " Weeks"
		} else if differHours < year {
			duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			duration = strconv.Itoa(int(years)) + " Years"
		}

		_, err = connection.Conn.Exec(context.Background(), "UPDATE db_project SET project_name=$1, start_date=$2, end_date=$3, duration=$4, description=$5, technologies=$6, image=$7 WHERE id=$8",
			projectName, startDate, endDate, duration, description, technologies, image, where_id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	where_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM db_project WHERE id=$1", where_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
