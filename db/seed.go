package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"oursos.com/packages/util"
)

func SeedDatabase() {
	db, err := Connection()
	username := "JohnDoe"
	longitude := -123.138570
	latitude := 49.263570
	languagepreference := "en"
	friends := []int{2, 3} // Sample user IDs as integers
	profile := "https://picsum.photos/200/300?grayscale"
	util.CheckError(err)

	createTables := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            latitude FLOAT NOT NULL,
            longitude FLOAT NOT NULL,  
            languagepreference VARCHAR(255),
            friends INT[],
            profile VARCHAR(255)  
        );
    
        CREATE TABLE IF NOT EXISTS alerts (
            id SERIAL PRIMARY KEY,
            message TEXT,
            category VARCHAR(50),
            severity int8,
            time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            latitude FLOAT NOT NULL,
            longitude FLOAT NOT NULL
        );
    `
	_, err = db.Exec(createTables)
	util.CheckError(err)

	_, err = db.Exec(`
        INSERT INTO users (username, longitude, latitude, languagepreference, friends, profile)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, username, longitude, latitude, languagepreference, pq.Array(friends), profile)
	util.CheckError(err)

	insertAlertSQL := `
        INSERT INTO alerts (message, category, severity, latitude, longitude)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = db.Exec(insertAlertSQL, "Emergency Alert", "Traffic Update", 0, 49.2827, -123.1207)
	util.CheckError(err)
}
