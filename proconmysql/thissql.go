package proconmysql

import (
	"fmt"
	"database/sql"
	"encoding/json"
	
	"github.com/gorilla/websocket"
	_ "github.com/go-sql-driver/mysql"
	"go_systems/pr0config"
	"go_systems/proconutil"
)

var DBCon *sql.DB

func init() {
	var err error
	DBCon, err = sql.Open("mysql", "root:"+pr0config.MysqlPass+"@tcp(localhost:3306)/")
	err = DBCon.Ping();
	if err != nil {
	    fmt.Println(err);
	}else {
		fmt.Println("MySql Connected...")
	}
	DBCon.SetMaxOpenConns(20)		
}



type GetMysqlDbsTask struct {
	ws *websocket.Conn
}

func NewGetMysqlDbsTask(ws *websocket.Conn) *GetMysqlDbsTask {
	return &GetMysqlDbsTask{ws}
}

func (rmdst *GetMysqlDbsTask) Perform() {
	var dbnames []string
	rows, err := DBCon.Query("SHOW DATABASES;")
	if err != nil {
		fmt.Println("Error in Perform of thissql: " err)
	}
	var dbs string
	for rows.Next() {
		rows.Scan(&dbs)
		dbnames = append(dbnames, dbs)
		fmt.Println(dbs)
	}
	jdbnames, _ := json.Marshal(dbnames)
		fmt.Println(string(jdbnames))

		proconutil.SendMsg("vAr", "mysql-dbs-list", string(jdbnames), rmdst.ws)
}

