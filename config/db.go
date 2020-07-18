package config

import (
	"os"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// Server : struct
type Server struct {
	Id       int
	Name     string
	Address  string
	Username string
	Password string
}

func opendbconnection() (database *sql.DB, error error) {

	user_home_dir, _ := os.UserHomeDir()
	db_location :=  user_home_dir + "/.connection-manager"
	db_file := db_location + "/servers.db"

	_, err := os.Stat(db_location)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(db_location, 0755)
		if errDir != nil {
			fmt.Println(err)
		}
	}
	database, _ = sql.Open("sqlite3", db_file)
	return
}

// GetAllServers : lists all the servers in DB
func GetAllServers() (servers []string) {

	database, err := opendbconnection()
	if err != nil {
		fmt.Printf("DB connection error %v\n", err)
		return
	}
	rows, derr := database.Query("SELECT name FROM servers")
	if derr != nil {
		fmt.Println("No server added for ssh connection. Add the server using # connection-manager add")
		return
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
	}
		servers = append(servers, name)
	}
	return
}

// GetDetailOfSpecificServer : Return server's raddress, username, password.
func GetDetailOfSpecificServer(server string) (host Server) {
	database, err := opendbconnection()
	if err != nil {
		fmt.Printf("DB connection error %v\n", err)
		return
	}
	q := fmt.Sprintf("SELECT id, name, address, username, password FROM servers WHERE name = '%s'", server)
	rows, _ := database.Query(q)

	var id int
	var name, address, username, password string
	for rows.Next() {
		err := rows.Scan(&id, &name, &address, &username, &password)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
	}
	}

	host = Server{
		Id:       id,
		Name:     name,
		Address:  address,
		Username: username,
		Password: password,
	}

	return
}

// UpdateHost : update the host details in the DB
func UpdateHost(host Server) {
	if host.Name != "" && host.Address != "" && host.Username != "" {

		database, err := opendbconnection()
		if err != nil {
			fmt.Printf("DB connection error %v\n", err)
			return
		}
		statement, _ := database.Prepare("UPDATE servers SET name = ?, address = ?, username = ?, password = ? WHERE id = ?")
		defer statement.Close()
		_, err = statement.Exec(host.Name, host.Address, host.Username, host.Password, host.Id)
		if err != nil {
			fmt.Printf("Error updaing DB %v\n", err)

		} else {
			fmt.Println("Server updated successfully")
		}
	} else {
		fmt.Println("All the required fields were not filled")
	}
}

// AddServerToDB : add server to DB
func AddServerToDB(host Server) {
	database, err := opendbconnection()
	if err != nil {
		fmt.Printf("DB connection error %v\n", err)
		return
	}
	if host.Name != "" && host.Address != "" && host.Username != "" {

		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS servers (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE, address TEXT, username TEXT, password TEXT )")
		_, err := statement.Exec()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		statement, _ = database.Prepare("INSERT INTO servers (name, address, username, password) VALUES (?, ?, ?, ?)")
		defer statement.Close()
		_, err = statement.Exec(host.Name, host.Address, host.Username, host.Password)
		if err != nil {
			fmt.Printf("fail to add server to DB %v\n", err)

		} else {
			fmt.Println("Server added successfully")
		}

	} else {
		fmt.Println("All the required fields were not filled")
	}
}

// DeleteServerFromDB : Delete the server entry from DB
func DeleteServerFromDB(server string) {
	database, err := opendbconnection()
	if err != nil {
		fmt.Printf("DB connection error %v\n", err)
		return
	}
	query := fmt.Sprintf("DELETE FROM servers WHERE name = '%s'", server)

	statement, _ := database.Prepare(query)
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
