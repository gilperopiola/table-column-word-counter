package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root"+":"+""+"@tcp("+"localhost"+":"+"3306"+")/"+"suips-insumos-migration"+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	table := "dependencies"
	column := "nombre"

	query := fmt.Sprintf(`SELECT %s FROM %s `, column, table)
	rows, err := db.Query(query)

	defer rows.Close()
	if err != nil {
		log.Fatalf("error getting rows: %v", err)
	}

	totalWords := make(map[string]int)
	for rows.Next() {
		field := ""
		err = rows.Scan(&field)
		if err != nil {
			log.Printf("error scanning field: %v", err)
		}

		fieldWords := strings.Split(field, " ")

		for _, fieldWord := range fieldWords {
			appearances, ok := totalWords[fieldWord]

			if ok {
				totalWords[fieldWord] = appearances + 1
			} else {
				totalWords[fieldWord] = 1
			}
		}
	}

	keys := make([]string, 0, len(totalWords))
	for key := range totalWords {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return totalWords[keys[i]] < totalWords[keys[j]] })

	for _, key := range keys {
		fmt.Printf("%s, %d\n", key, totalWords[key])
	}
}
