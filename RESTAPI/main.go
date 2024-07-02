package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	models "restapi/Models"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func connectDatabase() (*sql.DB, error) {

	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASSWORD")
	dbPort := os.Getenv("DBPORT")
	dbAddr := fmt.Sprintf("%s:%s", os.Getenv("DBADDRESS"), dbPort)
	dbName := os.Getenv("DBNAME")

	cfg := mysql.Config{
		User:   dbUser,
		Passwd: dbPassword,
		Net:    "tcp",
		Addr:   dbAddr,
		DBName: dbName,
	}

	var err error

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal("Erro conectar ao banco de dados", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Erro pingar ao banco de dados", err)
	}

	return db, nil
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	db, err := connectDatabase()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}
	defer db.Close()

	var album models.Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error populating album"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

func getAlbums(c *gin.Context) {
	db, err := connectDatabase()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer rows.Close()

	var albums []models.Album

	for rows.Next() {
		var album models.Album

		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		albums = append(albums, album)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(c *gin.Context) {
	var album models.Album

	if err := c.BindJSON(&album); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db, err := connectDatabase()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"success": lastID})
}

func updateAlbum(c *gin.Context) {
	id := c.Param("id")

	db, err := connectDatabase()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	var album models.Album
	err = c.BindJSON(&album)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := db.Exec("UPDATE album SET title = ?, artist = ?, price = ? WHERE id = ?", album.Title, album.Artist, album.Price, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"success": rowsAffected})
}

func deleteAlbum(c *gin.Context) {
	id := c.Param("id")
	db, err := connectDatabase()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM album WHERE id =?", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{"success": rowsAffected})
}

func main() {
	godotenv.Load()
	serverAddress := fmt.Sprintf("%s:%s", os.Getenv("SERVERADDR"), os.Getenv("SERVERPORT"))

	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.POST("/albums", addAlbum)
	router.PUT("/albums/:id", updateAlbum)
	router.DELETE("/albums/:id", deleteAlbum)

	router.Run(serverAddress)
}
