package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "loop@54321",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recording",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Album Found By Artist: ", albums)

	alb, err := albumByID(5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album Found By Row: %v\n", alb)

	albID, err := addAlbum(Album{
		Title:  "Avatar",
		Artist: "David Con",
		Price:  70.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of Added Album: %v\n ", albID)
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float64
}

// Query for multiple rows
// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	row, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("AlbumByArtist %q: %v", name, err)
	}
	defer row.Close()

	for row.Next() {
		var alb Album
		if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("AlbumByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("AlbumByArtist %q: %v", name, err)
	}
	return albums, nil
}

//Query for a single row

func albumByID(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE ID = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("AlbumByID %d: no such ID", id)
		}
		return alb, fmt.Errorf("AlbumByID %d: %v", id, err)
	}
	return alb, nil
}

//Add data

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album(title, artist, price) VALUES (?,?,?) ", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("AddAlbum : %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddAlbum : %v", err)
	}
	return id, nil
}
