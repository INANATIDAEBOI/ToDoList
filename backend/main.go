package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

// โครงสร้าง TodoList
type TodoList struct {
	ID    int    `json:"id"`
	Topic string `json:"topic"`
}

// ตัวแปรสำหรับเชื่อมต่อฐานข้อมูล
var db *sql.DB

// ฟังก์ชันตั้งค่าการเชื่อมต่อฐานข้อมูล
func SetupDatabase() *sql.DB {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "Insafeboy040745"
		dbname   = "mydatabase"
	)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// ตรวจสอบการเชื่อมต่อ
	err = database.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	fmt.Println("Successfully connected to database")
	return database
}
func main() {
	app := fiber.New()
	db = SetupDatabase()
	defer db.Close()

	app.Post("/api/todolist", CreateList)
	app.Get("/api/todolist", GetList)
	log.Fatal(app.Listen(":3000"))

}

func GetList(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, topic FROM list_table") // ✅ ใช้ Query() สำหรับหลายแถว
	if err != nil {
		return err
	}
	defer rows.Close()

	var lists []TodoList
	for rows.Next() {
		var t TodoList
		err := rows.Scan(&t.ID, &t.Topic)
		if err != nil {
			return err
		}
		lists = append(lists, t)
	}

	return c.JSON(lists)
}

func CreateList(c *fiber.Ctx) error {
	t := new(TodoList)
	if err := c.BodyParser(t); err != nil {
		return err
	}

	// ✅ ใช้ QueryRow() สำหรับ RETURNING id
	err := db.QueryRow("INSERT INTO list_table (topic) VALUES ($1) RETURNING id", t.Topic).Scan(&t.ID)

	if err != nil {
		return err
	}

	return c.JSON(t)
}
