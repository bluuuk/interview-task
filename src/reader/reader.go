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
	Username, Password, Host, Name string
	Port                           int
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

func ConnectToDB(db Database) (context.Context, *pgx.Conn) {
	// build connection url according to https://github.com/jackc/pgx
	connect_url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)

	log.Println("Connecting to :" + connect_url)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connect_url)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected successfully")

	return ctx, conn
}

func Read(input_file, result_file string, db Database) {

	f_in, err := os.Open(input_file)

	if err != nil {
		panic("File does not exist!")
	}

	// after method execution, close the file
	defer f_in.Close()

	// setup connection
	ctx, conn := ConnectToDB(db)
	defer conn.Close(ctx)

	// we count the occurrences
	counter := make(map[string]int, 1e7)

	// A scanner will look for lines i.e. by "\n" terminated strings
	scanner := bufio.NewScanner(f_in)
	scanner.Split(bufio.ScanLines)

	// Send entries over
	// ch_token := make(chan string)
	// ch_done := make(chan bool)

	log.Println("Start scanning file")

	// Go line by line and count the number of occurrences
	for scanner.Scan() {
		token := scanner.Text()

		// if the token is not in the map, value will be the default int => 0
		value, check := counter[token]
		counter[token] = value + 1

		// send the token over if it is the first occurrence
		if !check {
			// We do not need to prepare/cache statements, the lib will do https://github.com/jackc/pgx/issues/791
			cmdtag, err := conn.Exec(ctx, "INSERT INTO tokens(token) VALUES ($1)", token)
			if err != nil {
				log.Printf("Failed for token %s with result %s", token, cmdtag)
			}
			//ch_token <- token
		}
	}

	log.Println("Finish scanning file")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// signal that no more tokens will come
	// ch_done <- false

	// https://stackoverflow.com/questions/8593645/is-it-ok-to-leave-a-channel-open
	// defer resource clean up
	// defer close(ch_done)
	// defer close(ch_token)

	// Write duplicate tokens into a file with their freq

	f_out, err := os.Create(result_file)

	if err != nil {
		panic("Could not create file")
	}

	// resource clean up
	defer f_out.Close()

	writer := bufio.NewWriter(f_out)
	// header for csv
	writer.WriteString("token,freq\n")

	log.Println("Writing duplicates to file")

	// iterate over map and write to file iff a token occurred more than once
	collision := 0
	for token, count := range counter {
		if count > 1 {
			writer.WriteString(
				fmt.Sprintf(
					"%s,%d\n",
					token, count,
				))
			collision += 1
		}
	}

	defer writer.Flush()

	log.Printf("Observed %d collision with a collision rate of %.7f%%\n",
		collision,
		1-float32(collision)/float32(len(counter)),
	)

}
