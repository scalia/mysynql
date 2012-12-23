package mysql

import (
	"github.com/scalia/mysynql/log"
	"fmt"
)

func (database *Database) ReadConnection(conn *Connection) {
	dbname := conn.DbName
	sql := "SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME" +
		" FROM INFORMATION_SCHEMA.SCHEMATA" +
		" WHERE SCHEMA_NAME=?"
	stmt, err := conn.Prepare(sql)
	if nil != err { // Unknown error happened.
		panic(err)
	}
	stmt.Bind(dbname)

	rows, res, err := stmt.Exec()
	if nil != err {
		panic(err)
	}

	if 1 != len(rows) {
		panic(fmt.Sprintf("database %s not found", dbname))
	}
	row := rows[0]

	database.Name    = dbname
	database.Charset = row.Str(res.Map("DEFAULT_CHARACTER_SET_NAME"))
	database.Collation = row.Str(res.Map("DEFAULT_COLLATION_NAME"))

	sql = "SELECT TABLE_NAME, TABLE_TYPE, ENGINE, TABLE_COLLATION" +
		" FROM INFORMATION_SCHEMA.TABLES" +
		" WHERE TABLE_SCHEMA=?"
	stmt, err = conn.Prepare(sql)
	if nil != err { // Unknown error happened.
		panic(err)
	}
	stmt.Bind(dbname)

	rows, res, err = stmt.Exec()
	if nil != err {
		panic(err)
	}

	database.Tables = make([]Table, 0)
	for _, row := range rows {
		tableType := row.Str(res.Map("TABLE_TYPE"))
		if "BASE TABLE" != tableType {
			log.Error(fmt.Sprintf("Unsupported table type \"%s\"", tableType))
			continue
		}

		var table Table
		table.Name = row.Str(res.Map("TABLE_NAME"))
		table.Type = tableType
		table.Engine = row.Str(res.Map("ENGINE"))
		table.Collation = row.Str(res.Map("TABLE_COLLATION"))

		database.Tables = append(database.Tables, table)
	}

	channel := make(chan bool)
	for index, _ := range database.Tables {
		go database.Tables[index].ReadConnection(conn, channel)
	}

	result := true
	for i:= len(database.Tables); i>0; i-- {
		result = result && <- channel
	}
	close(channel)

	if result {
		log.Verbose("COMMIT")
	} else {
		log.Verbose("ROLLBACK")
	}
}
