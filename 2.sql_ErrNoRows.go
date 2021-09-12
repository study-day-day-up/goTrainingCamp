package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Team struct {
	id         uint32
	team_name  string
	team_order int8
	which_year int16
}

func main() {
	db, err := sql.Open("mysql",
		"root:password@tcp(127.0.0.1:3306)/go")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	team := Team{}
	err = db.QueryRow("select team_name from teams_temp where which_year = ?", 2021).Scan(&team.team_name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatal(err) // 输出 sql: no rows in result set
		} else {
			log.Fatal(err)
		}
	}
	fmt.Println(team.team_name)
}
