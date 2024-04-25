package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist =?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

// editAlbum update the specified album to the database,
// returning the album ID of the updated entry
func editAlbum(id int64, alb Album) (int64, error) {
	result, err := db.Exec("UPDATE album SET title=?, artist=?, price=? WHERE id=?", alb.Title, alb.Artist, alb.Price, id)
	if err != nil {
		return 0, fmt.Errorf("editAlbum: %v", err)
	}
	returnId, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("editAlbum: %v", err)
	}
	return returnId, nil
}

// deleteAlbum delete the specified album  with given id from the database,
// returning the album ID of the deleted entry
func deleteAlbum(id int64) (int64, error) {
	result, err := db.Exec("DELETE FROM album WHERE id=?", id)

	if err != nil {
		return 0, fmt.Errorf("deleteAlbum: %v", err)
	}
	returnId, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("deleteAlbum: %v", err)
	}
	return returnId, nil
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:    os.Getenv("DBUSER"),
		Passwd:  os.Getenv("DBPASS"),
		Net:     "tcp",
		Addr:    "127.0.0.1:3306",
		DBName:  "recordings",
		Timeout: time.Duration(time.Second * 20),
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	album, err := albumByID(4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", album)

	insertedAlbum := Album{
		Title: "Cinta Kasih New",
		Artist: "Firman New",
		Price: 43.33,
	}
	id, err := addAlbum(insertedAlbum)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Inserted album id %v, data: %v\n", id, insertedAlbum)

	newAlbum := Album{
		Title: "New New Album",
		Artist: "New New Firman",
		Price: 99.33,
	}
	updateId, err := editAlbum(2,newAlbum)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Updated album id %v, data: %v\n", updateId, newAlbum)

	rowAffected, err := deleteAlbum(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album deleted, row affected: %v\n", rowAffected )


}
