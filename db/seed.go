package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"oursos.com/packages/util"
)

func SeedDatabase() {
	db, err := Connection()
	deviceId := "ingrid"
	username := "Jun"
	lat := 37.938194
	long := 126.587791
	languagepreference := "en"
	friends := []int{1} // Sample user IDs as integers
	profile := "https://picsum.photos/200/300?grayscale"
	util.CheckError(err)

	createTables := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
		deviceId VARCHAR(255) NOT NULL,
            username VARCHAR(255) NOT NULL,
            lat FLOAT NOT NULL,
            long FLOAT NOT NULL,  
            languagepreference VARCHAR(255),
            friends INT[],
            profile VARCHAR(255)  
        );
    
        CREATE TABLE IF NOT EXISTS alerts (
            id SERIAL PRIMARY KEY,
            message TEXT,
            type VARCHAR(50),
            severity int8,
            time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            lat FLOAT NOT NULL,
            long FLOAT NOT NULL
        );
    `
	_, err = db.Exec(createTables)
	util.CheckError(err)

	_, err = db.Exec(`
        INSERT INTO users (deviceId,username, lat, long, languagepreference, friends, profile)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, deviceId, username, lat, long, languagepreference, pq.Array(friends), profile)
	util.CheckError(err)

	insertAlertSQL := `
        INSERT INTO alerts (message, type, severity, lat, long)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = db.Exec(insertAlertSQL, "Emergency Alert", "Traffic Update", 0, 49.2827, -123.1207)
	util.CheckError(err)
}
