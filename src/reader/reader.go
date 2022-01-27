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

const query = "INSERT INTO tokens(token) VALUES ($1)"

/*
	Writes unique tokens to the database in one go
*/
func WriteToDB(conn *pgx.Conn, token_stream <-chan string, query string) {

	log.Println("Start writing to DB")
	var data []string

	// Adding every INSERT to a batch such that we only have one big transfer to the DB
	for token := range token_stream {
		data = append(data, token)
	}

	// Using COPY to write to the DB (https://pkg.go.dev/github.com/jackc/pgx/v4#hdr-Copy_Protocol)
	_, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"tokens"},
		[]string{"token"},
		pgx.CopyFromSlice(len(data), func(i int) ([]interface{}, error) {
			return []interface{}{data[i]}, nil
		}),
	)

	log.Println("Finish writing to DB")

	if err != nil {
		log.Print(err)
	}

	defer conn.Close(context.Background())
}

/*
	Creates a single database connection
*/

func ConnectToDB(db Database) *pgx.Conn {
	// build connection url according to https://github.com/jackc/pgx
	connect_url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)

	log.Println("Connecting to " + connect_url)
	conn, err := pgx.Connect(context.Background(), connect_url)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected successfully")

	return conn
}

/*
	Read an input file, write unique tokens to a data base and output duplicates with their frequency into a file
*/

func Read(input_file, result_file string, conn *pgx.Conn) {

	f_in, err := os.Open(input_file)

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
	ch_token := make(chan string)

	// We do not want to finish this routine if not all tokens are written into the db
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		WriteToDB(conn, ch_token, query)
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

	// header for csv
	writer := bufio.NewWriter(f_out)
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

	// ensure at least one flush such that the buffer is empty
	defer writer.Flush()

	// some statistics
	log.Printf("Observed %d collision with a collision rate of %.7f%%\n",
		collision,
		1-float32(collision)/float32(len(counter)),
	)

	// to ensure that writing to the db finishes
	wg.Wait()
}
