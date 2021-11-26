package reader

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

type Database struct {
	username, password, host, name string
	port                           int
}

/*
func writeToDB(db Database, token_stream <-chan string) {

	// build connection url according to https://github.com/jackc/pgx
	connect_url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		db.username,
		db.password,
		db.host,
		db.port,
		db.name,
	)

	log.Println("Connecting to :" + connect_url)

	k := pgx.ConnConfig{}
	pgx.Connect
	conn, err := pgx.Connect(context.Background(), connect_url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected successfully")

	defer conn.Close(context.Background())

	for token := range token_stream {
		_, err := conn.Exec()
	}
}

*/

func Read(file string, db Database) {

	f, err := os.Open(file)

	if err == nil {
		panic("File does not exist!")
	}

	// after method execution, close the file
	defer f.Close()

	// build connection url according to https://github.com/jackc/pgx
	connect_url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		db.username,
		db.password,
		db.host,
		db.port,
		db.name,
	)

	log.Println("Connecting to :" + connect_url)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connect_url)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected successfully")

	defer conn.Close(ctx)

	// we count the occurrences
	counter := make(map[string]int, 1e7)

	// A scanner will look for lines i.e. by "\n" terminated strings
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	// Send entries over
	// ch_token := make(chan string)
	// ch_done := make(chan bool)

	// Go line by line and count the number of occurrences
	for scanner.Scan() {
		token := scanner.Text()

		// if the token is not in the map, value will be the default int => 0
		value, check := counter[token]
		counter[token] = value + 1

		// send the token over if it is the first occurrence
		if !check {
			// We do not need to prepare/cache statements, the lib will do https://github.com/jackc/pgx/issues/791
			cmdtag, err := conn.Exec(ctx, "INSERT TOKEN(TOKEN) VALUES ($1)", token)
			if err != nil {
				log.Printf("Failed for token %s with result %s", token, cmdtag)
			}
			//ch_token <- token
		}
	}

	// signal that no more tokens will come
	// ch_done <- false

	// https://stackoverflow.com/questions/8593645/is-it-ok-to-leave-a-channel-open
	// defer resource clean up
	// defer close(ch_done)
	// defer close(ch_token)

	// Transfer every token to the data base
	for token, count := range counter {
		if count > 1 {
			fmt.Printf(
				"%s - %d",
				token, count,
			)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
