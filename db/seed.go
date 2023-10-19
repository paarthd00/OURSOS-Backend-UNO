package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"oursos.com/packages/util"
)

func SeedDatabase() {
	db, err := Connection()
	username := "JohnDoe"
	locations := []string{"(40.7128, -74.0060)", "(34.0522, -118.2437)"} // Sample latitude and longitude points
	languagepreference := "English"
	friends := []int{2, 3} // Sample user IDs as integers
	util.CheckError(err)

	createTables := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            locations POINT[],  -- Using the POINT type to store latitude and longitude as an array
            languagepreference VARCHAR(255),
            friends INT[]  -- Storing friend user IDs as an array of integers
        );
    
        CREATE TABLE IF NOT EXISTS alerts (
            id SERIAL PRIMARY KEY,
            message TEXT,
            category VARCHAR(50),
            severity VARCHAR(50),
            time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            latitude VARCHAR(10) NOT NULL,
            longitude VARCHAR(10) NOT NULL,
            radius VARCHAR(10)
        );
    `
	_, err = db.Exec(createTables)
	util.CheckError(err)

	_, err = db.Exec(`
        INSERT INTO users (username, locations, languagepreference, friends)
        VALUES ($1, $2, $3, $4)
    `, username, pq.Array(locations), languagepreference, pq.Array(friends))
	util.CheckError(err)

	insertAlertSQL := `
        INSERT INTO alerts (message, category, severity, latitude, longitude, radius)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err = db.Exec(insertAlertSQL, "Emergency Alert", "Traffic Update", "low", "49.2827", "-123.1207", "5.0")
	util.CheckError(err)
}
