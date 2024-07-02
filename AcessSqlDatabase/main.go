package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float64
}

func getAlbumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist LIKE ?", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func addAlbum(album Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (artist, title, price) VALUES (?, ?, ?)", album.Artist, album.Title, album.Price)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return lastId, nil
}

func main() {
	godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWD")
	dbAddr := fmt.Sprintf("%s:%s", os.Getenv("DB_ADDR"), os.Getenv("DB_PORT"))
	dbName := os.Getenv("DB_NAME") // Must have a valid database, the driver doesn't create one

	cfg := mysql.Config{
		User:   dbUser,
		Passwd: dbPassword,
		Net:    "tcp",
		Addr:   dbAddr,
		DBName: dbName,
	}

	var err error

	album := Album{
		Title:  "New Album",
		Artist: "New Artist",
		Price:  9.99,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("Sucessfully connected to database")

	result, _ := getAlbumsByArtist("a")
	insertedId, err := addAlbum(album)

	fmt.Println("Inserted id: ", insertedId)
	fmt.Println(result)
}
