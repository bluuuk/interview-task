package reader

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v4"
)

type Database struct {
	Username, Password, Host, Name string
	Port                           int
}

func writeToDB(ctx context.Context, conn *pgx.Conn, token_stream <-chan string) {

	log.Println("Start writing to DB")

	// caching happens in the lib pgx
	const query = "INSERT INTO tokens(token) VALUES ($1)"

	batch := &pgx.Batch{}

	for token := range token_stream {
		batch.Queue(query, token)
	}

	batch_request := conn.SendBatch(ctx, batch)
	defer batch_request.Close()
	_, err := batch_request.Exec()

	if err != nil {
		panic(err)
	}

	log.Println("Finish writing to DB")
}

func ConnectToDB(db Database) (context.Context, *pgx.Conn) {
	// build connection url according to https://github.com/jackc/pgx
	connect_url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)

	log.Println("Connecting to " + connect_url)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connect_url)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected successfully")

	return ctx, conn
}

func Read(input_file, result_file string, ctx context.Context, conn *pgx.Conn) {

	f_in, err := os.Open(input_file)
	var wg sync.WaitGroup

	if err != nil {
		panic("File does not exist!")
	}

	// after method execution, close the file
	defer f_in.Close()

	// we count the occurrences
	counter := make(map[string]int, 1e7)

	// A scanner will look for lines i.e. by "\n" terminated strings
	scanner := bufio.NewScanner(f_in)
	scanner.Split(bufio.ScanLines)

	// Prepare go routine for sql insertions
	// We do not want to finish this routine if not all tokens are written into the db
	ch_token := make(chan string)
	wg.Add(1)

	go func() {
		defer wg.Done()
		writeToDB(ctx, conn, ch_token)
	}()

	log.Println("Start scanning file")

	// Go line by line and count the number of occurrences
	for scanner.Scan() {
		token := scanner.Text()

		// if the token is not in the map, value will be the default int => 0
		value, check := counter[token]
		counter[token] = value + 1

		// send the token over if it is the first occurrence
		if !check {
			ch_token <- token
		}
	}

	// close the channel so signal that no new tokens will come
	close(ch_token)

	log.Println("Finish scanning file")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

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

	// some statistics
	log.Printf("Observed %d collision with a collision rate of %.7f%%\n",
		collision,
		1-float32(collision)/float32(len(counter)),
	)

	// to ensure that writing to the db finishes
	wg.Wait()
}
