package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Slot struct {
	Title        string   `json:"Title"`
	Owner        string   `json:"Owner"`
	Date         string   `json:"Date"`
	Time         int      `json:"Time"`
	Participants []string `json:"Participants"`
}

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := title + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, 1)
}

func loginRedirecter(w http.ResponseWriter, r *http.Request) {

	credential := r.FormValue("Name")
	if credential == "F10" {
		http.Redirect(w, r, "/dashboard2", http.StatusFound)
	} else {
		http.Redirect(w, r, "/dashboard1", http.StatusFound)
	}
}
func taskRedirecter(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("task")
	if task == "add" {
		http.Redirect(w, r, "/reserveSlot", http.StatusFound)
	} else if task == "retrieve" {
		http.Redirect(w, r, "/retrieveSlot", http.StatusFound)
	} else if task == "delete" {
		http.Redirect(w, r, "/deleteSlot", http.StatusFound)
	} else {
		http.Redirect(w, r, "/monthSlot", http.StatusFound)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func reserveSave(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("slotDate")
	time := r.FormValue("slotTime")
	fmt.Println(time)
	var rSlot Slot
	rSlot.Date = date
	rSlot.Time, _ = strconv.Atoi(time[:2])
	r.ParseForm()
	rSlot.Participants = r.Form["participants"]
	rSlot.Title = r.FormValue("slotTitle")
	rSlot.Owner = r.FormValue("slotOwner")
	//fmt.Println(rSlot.Participants)
	//file, _ := json.Marshal(rSlot)

	//ioutil.WriteFile("test.json", file, os.ModePerm)

	temp_arr := []Slot{}
	file2, _ := ioutil.ReadFile("test.json")
	//fmt.Printf("hey %s", file2)
	json.Unmarshal(file2, &temp_arr)
	// var sSlot Slot
	// sSlot.Owner = "Sparkzz"
	// sSlot.Title = "Hello"

	var flag int = 0

	for i := 0; i < len(temp_arr); i++ {
		if temp_arr[i].Date == rSlot.Date {
			if temp_arr[i].Time == rSlot.Time {
				for j := 0; j < len(temp_arr[i].Participants); j++ {
					for k := 0; k < len(rSlot.Participants); k++ {
						if temp_arr[i].Participants[j] == rSlot.Participants[k] {
							fmt.Println("Meeting cannot be taken")
							flag = 1
							http.Redirect(w, r, "/failure", http.StatusFound)
							break
						}
					}
				}
			}
		}
	}
	if flag == 0 {
		//fmt.Println("heyaaa1  ", temp_arr)
		temp_arr = append(temp_arr, rSlot)
		//fmt.Println("heyaaa2  ", temp_arr)
		file3, _ := json.Marshal(temp_arr)
		ioutil.WriteFile("test.json", file3, os.ModePerm)
	}
	http.Redirect(w, r, "/login", http.StatusFound)
	//fmt.Println("Helllo")
}

func retrieveSchedule(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("slotDate")
	owner := r.FormValue("slotOwner")
	temp_arr := []Slot{}
	file2, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(file2, &temp_arr)
	var result []Slot
	for i := 0; i < len(temp_arr); i++ {
		if temp_arr[i].Date == date {
			length := len(temp_arr[i].Participants)
			for j := 0; j < length; j++ {
				if temp_arr[i].Participants[j] == owner {
					result = append(result, temp_arr[i])
				}
			}
		}
	}
	fmt.Println(result)
	file, _ := json.MarshalIndent(result, "", " ")
	ioutil.WriteFile("result.json", file, os.ModePerm)
	http.Redirect(w, r, "/resultSchedule", http.StatusFound)

}

func monthInfoHandler(w http.ResponseWriter, r *http.Request) {
	year := r.FormValue("slotYear")
	month := r.FormValue("slotMonth")

	combined := year + "-" + month
	file2, _ := ioutil.ReadFile("test.json")
	temp_arr := []Slot{}
	json.Unmarshal(file2, &temp_arr)
	var result []Slot
	//fmt.Println(combined + "\n")
	for i := 0; i < len(temp_arr); i++ {
		if temp_arr[i].Date[:7] == combined {
			result = append(result, temp_arr[i])
			//fmt.Println("Reached here")
		}
	}
	file, _ := json.MarshalIndent(result, "", " ")
	ioutil.WriteFile("resultMonth.json", file, os.ModePerm)
	http.Redirect(w, r, "/monthResult", http.StatusFound)

}

func deleteSlotHandler(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("slotDate")
	owner := r.FormValue("slotOwner")
	time := r.FormValue("slotTime")

	file2, _ := ioutil.ReadFile("test.json")
	temp_arr := []Slot{}
	json.Unmarshal(file2, &temp_arr)
	timing, _ := strconv.Atoi(time[:2])
	var result []Slot
	flag := 0
	for i := 0; i < len(temp_arr); i++ {
		if temp_arr[i].Date != date {
			result = append(result, temp_arr[i])
		} else if temp_arr[i].Owner != owner {
			result = append(result, temp_arr[i])
		} else if temp_arr[i].Time != timing {
			result = append(result, temp_arr[i])
		} else {
			flag = 1
		}
	}
	if flag == 0 {
		http.Redirect(w, r, "/deleteFailure", http.StatusFound)
		return
	}
	fmt.Println(result)
	file, _ := json.MarshalIndent(result, "", " ")
	ioutil.WriteFile("test.json", file, os.ModePerm)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func reserveHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "reserveSlot")
}
func retrieveHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "retrieveSlot")
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "deleteSlot")
}
func monthHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "monthSlot")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "login")
}
func normalDashboard(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "dashboard1")
}
func adminDashboard(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "dashboard2")
}

func failureHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "failure")
}

func deleteFailureHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "deleteFailure")
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {

	p, _ := loadPage("result")

	t, _ := template.ParseFiles("resultSchedule.html")
	t.Execute(w, p)
}

func monthResultHandler(w http.ResponseWriter, r *http.Request) {

	p, _ := loadPage("resultMonth")

	t, _ := template.ParseFiles("monthResult.html")
	t.Execute(w, p)
}

func main() {

	if !fileExists("test.json") {
		var temp_slots []Slot
		file, _ := json.Marshal(temp_slots)
		ioutil.WriteFile("test.json", file, os.ModePerm)
	}
	//fmt.Println(takenSlots[0].owner)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/save/", loginRedirecter)
	http.HandleFunc("/dashboard1/", normalDashboard)
	http.HandleFunc("/dashboard2/", adminDashboard)
	http.HandleFunc("/taskHandler/", taskRedirecter)
	http.HandleFunc("/reserveSlot/", reserveHandler)
	http.HandleFunc("/retrieveSlot/", retrieveHandler)
	http.HandleFunc("/deleteSlot/", deleteHandler)
	http.HandleFunc("/monthSlot/", monthHandler)
	http.HandleFunc("/reserveSave/", reserveSave)
	http.HandleFunc("/failure/", failureHandler)
	http.HandleFunc("/deleteFailure/", deleteFailureHandler)
	http.HandleFunc("/retrieveSchedule/", retrieveSchedule)
	http.HandleFunc("/resultSchedule/", scheduleHandler)
	http.HandleFunc("/retrieveMonth/", monthInfoHandler)
	http.HandleFunc("/monthResult/", monthResultHandler)
	http.HandleFunc("/deleteSchedule/", deleteSlotHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
