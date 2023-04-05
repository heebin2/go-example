package helper

import (
	"fmt"
	"log"
	"strconv"
)

// Information for database connection
type DBInfo struct {

	// db host name, ex) 127.0.0.1
	Host string

	// user name, ex) tms
	User string

	// user password
	Password string

	// database name
	Database string

	// database connect port
	Port int
}

// Outputs dbinfo in a format suitable for each dbms.
func (d *DBInfo) Config(dbms string) string {
	switch dbms {
	case "mysql":
		return d.User + ":" + d.Password + "@tcp(" + d.Host + ":" + strconv.Itoa(d.Port) + ")/" + d.Database
	case "postgres":
		return "host=" + d.Host + " port=" + strconv.Itoa(d.Port) + " user=" + d.User + " password=" + d.Password + " dbname=" + d.Database + " sslmode=disable"
	default:
		log.Panic(fmt.Errorf("not supported dbms : %s", dbms))
		return ""
	}
}
