package proconmysql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"go_systems/pr0config"
	"go_systems/proconutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)
// DBCon is a reference to the DB
var DBCon *sql.DB

func init() {
	var err error
	DBCon, err = sql.Open("mysql", "root:"+pr0config.MysqlPass+"@tcp(localhost:3306)/")
	err = DBCon.Ping();
	if err != nil {
	    fmt.Println(err);
	}
	fmt.Println("MySql Connected...")

	DBCon.SetMaxOpenConns(20)		
}


// GetMysqlDbsTask is a struct to pass the websocket connection
type GetMysqlDbsTask struct {
	ws *websocket.Conn
}

// NewGetMysqlDbsTask creates teh ws link
func NewGetMysqlDbsTask(ws *websocket.Conn) *GetMysqlDbsTask {
	fmt.Println("In NewGetMysqlDvsTask line37")
	return &GetMysqlDbsTask{ws}
}
// Perform is used in the async task channels
func (rmdst *GetMysqlDbsTask) Perform() {
	var dbnames []string
	rows, err := DBCon.Query("SHOW DATABASES;")
	fmt.Println("In GetMysqlDvsTask line44")
	if err != nil {
		fmt.Println("Error in Perform of thissql: ", err)
	}
	var dbs string
	fmt.Println("In GetMysqlDvsTask line49")
	for rows.Next() {
		rows.Scan(&dbs)
		dbnames = append(dbnames, dbs)
		fmt.Println(dbs)
	}
	jdbnames, _ := json.Marshal(dbnames)
		fmt.Println(string(jdbnames))

		proconutil.SendMsg("vAr", "mysql-dbs-list", string(jdbnames), rmdst.ws)
		fmt.Println("In GetMysqlDvsTask line59", string(jdbnames) )
}

