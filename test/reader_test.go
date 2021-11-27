package generator_test

import (
	"context"
	"interview/src/reader"
	"testing"
	"time"

	"github.com/google/uuid"
)

var db reader.Database = reader.Database{
	Username: "postgres",
	Password: "postgres",
	Host:     "localhost",
	Name:     "interview",
	Port:     5433,
}

const batch_size int = 10_000_000

func TestDBWrite(test *testing.T) {

	conn := reader.ConnectToDB(db)

	// test creation time

	_, err := conn.Exec(context.Background(), "CREATE TABLE test(token VARCHAR(7) NOT NULL)")

	if err != nil {
		test.Error("Could not create test table")
	}

	tokens := make(chan string, batch_size)

	// writing random uuid string
	for k := 0; k < batch_size; k++ {
		// strings of length 7
		tokens <- uuid.NewString()[:7]
	}
	// signal that we are done
	close(tokens)

	test.Log("Done generating tokens")

	const query = "INSERT INTO test(token) VALUES ($1)"

	start := time.Now()
	reader.WriteToDB(conn, tokens, query)
	duration := time.Since(start)

	test.Logf("Took %s to write %d tokens into the db", duration, batch_size)

	// test number of tokens in the db

	var amount int
	conn.QueryRow(context.Background(), "Select Count(*) from test").Scan(&amount)

	if amount != batch_size {
		test.Errorf("Expected amount %d, actual amount %d", batch_size, amount)
	}

	// cleanup

	_, err = conn.Exec(context.Background(), "DROP TABLE test")

	if err != nil {
		test.Error("Could not clean up table")
	}

	defer conn.Close(context.Background())
}
