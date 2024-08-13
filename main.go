package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "8878"
	dbname   = "Question-DB"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	CheckError(err)

	defer db.Close()

	err = db.Ping()
	CheckError(err)

	fmt.Println("Successfully connected!")

	InsertToTable := func() {
		// dynamic
		insertDynStmt := `INSERT INTO C10M1MCQ (question, option1, option2, option3, option4, Answer, Explanation)
	VALUES (
	    'What is the capital of France?',
	    'London',
	    'Berlin',
	    'Paris',
	    'Madrid',
	    3,
	    'Paris is the capital of France.'
	);
	`
		_, err = db.Exec(insertDynStmt)
	}

	ShowAll := func(db *sql.DB) ([]string, error) {
		var data []string
		Query := "Select * from C10M1MCQ;"
		rows, err := db.Query(Query)
		CheckError(err)
		defer rows.Close()

		for rows.Next() {
			var QNo int
			var Question string
			var option1 string
			var option2 string
			var option3 string
			var option4 string
			var Answer int
			var Explanation string

			err = rows.Scan(&QNo, &Question, &option1, &option2, &option3, &option4, &Answer, &Explanation)
			CheckError(err)
			fmt.Println(QNo, Question, option1, option2, option3, option4, Answer, Explanation)
			//data:=[]string{Question, option1, option2, option3, option4, Explanation}
			// Build a single row string
			rowString := fmt.Sprintf("%d, %s, %s, %s, %s, %s, %d, %s", QNo, Question, option1, option2, option3, option4, Answer, Explanation)
			data = append(data, rowString)
		}

		return data, nil
	}
	//fmt.Printf("output: %v\n", output)
	//fmt.Println(output)
	if 7 < 3 {
		InsertToTable()
		ShowAll(db)
	}

	/*
		app := fiber.New()

		app.Static("/", "./index.html",data)
	*/
	app := fiber.New(fiber.Config{
		Views: html.New("./", ".html"),
	})

	app.Get("/", func(c *fiber.Ctx) error {
		data, err := ShowAll(db)
		if err != nil {
			return err // Handle error gracefully
		}
		return c.Render("index", fiber.Map{
			"data": data,
		})
	})
	log.Fatal(app.Listen(":3000"))

}
