package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	config "github.com/tokopedia/gosample/config"
	redis "github.com/tokopedia/gosample/redis"

	_ "github.com/lib/pq"
)

type UserTable struct {
	ID         int       `db:"user_id"`
	FullName   string    `db:"full_name"`
	UserEmail  string    `db:"user_email"`
	MSISDN     string    `db:"msisdn"`
	BirthDate  time.Time `db:"birth_date"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

type User struct {
	ID         int    `json:"id"`
	FullName   string `json:"full_name"`
	UserEmail  string `json:"email"`
	MSISDN     string `json:"msisdn"`
	BirthDate  string `json:"birth_date"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type RenderUser struct {
	Title     string
	UserData  []User
	ViewCount int
}

const getUserListQuery = `
	SELECT 
		user_id,
		full_name,
		msisdn,
		birth_date,
		create_time,
		update_time
	FROM 
		ws_user 
	ORDER BY 
		full_name ASC
	LIMIT 
		10 
	
`

const filterUserListQuery = `
	SELECT 
		user_id,
		full_name,
		msisdn,
		birth_date,
		create_time,
		update_time
	FROM 
		ws_user 
	WHERE 
		full_name LIKE $1
	ORDER BY 
		full_name ASC
	LIMIT 
		10 
	
`

func GetUsers(w http.ResponseWriter, r *http.Request) {

	//view count
	existKey, err := redis.Exists("tech_curriculum_juna_view_count")

	if err != nil {
		fmt.Println(err.Error())
	}

	var viewCount int

	if existKey == true {
		var err error
		count, err := redis.Get("tech_curriculum_juna_view_count")
		if err != nil {
			fmt.Println(err.Error())
		}

		viewCount, err = strconv.Atoi(string(count))

		fmt.Println("first redis ", viewCount)
		if err != nil {
			fmt.Println(err.Error())
		}
		viewCount++

		redis.Set("tech_curriculum_juna_view_count", []byte(strconv.Itoa(viewCount)))

	} else {
		redis.Set("tech_curriculum_juna_view_count", []byte("1"))
		viewCount = 1
	}

	rows, err := config.Db.Query(getUserListQuery)
	if err != nil {
		fmt.Println(err.Error())
	}
	var Users = []User{}

	for rows.Next() {
		tempUserTable := UserTable{}

		err = rows.Scan(&tempUserTable.ID, &tempUserTable.FullName, &tempUserTable.MSISDN, &tempUserTable.BirthDate, &tempUserTable.CreateTime, &tempUserTable.UpdateTime)

		tempUser := User{
			ID:         tempUserTable.ID,
			FullName:   tempUserTable.FullName,
			UserEmail:  tempUserTable.UserEmail,
			MSISDN:     tempUserTable.MSISDN,
			BirthDate:  time.Time.String(tempUserTable.BirthDate),
			CreateTime: time.Time.String(tempUserTable.CreateTime),
			UpdateTime: time.Time.String(tempUserTable.UpdateTime),
		}

		Users = append(Users, tempUser)
	}

	//fmt.Println(Users)
	defer rows.Close()
	data := RenderUser{
		Title:     "Daftar 10 pengguna teratas",
		UserData:  Users,
		ViewCount: viewCount,
	}

	tmpl, err := template.ParseFiles("view/index.html")
	tmpl.Execute(w, data)

}

func SearchUsers(w http.ResponseWriter, r *http.Request) {

	filter := r.FormValue("key")
	rows, err := config.Db.Query(filterUserListQuery, "%"+filter+"%")
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var Users = []User{}

	for rows.Next() {
		tempUserTable := UserTable{}

		err = rows.Scan(&tempUserTable.ID, &tempUserTable.FullName, &tempUserTable.MSISDN, &tempUserTable.BirthDate, &tempUserTable.CreateTime, &tempUserTable.UpdateTime)

		tempUser := User{
			ID:         tempUserTable.ID,
			FullName:   tempUserTable.FullName,
			UserEmail:  tempUserTable.UserEmail,
			MSISDN:     tempUserTable.MSISDN,
			BirthDate:  time.Time.String(tempUserTable.BirthDate),
			CreateTime: time.Time.String(tempUserTable.CreateTime),
			UpdateTime: time.Time.String(tempUserTable.UpdateTime),
		}

		Users = append(Users, tempUser)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	byt, err := json.Marshal(Users)
	if err != nil {
		w.Write([]byte("Erorr nih" + err.Error()))
		return
	}

	w.Write(byt)

}
