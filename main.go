package main

import (
	"interview/src/generator"
	"interview/src/reader"
)

const (
	input_file  = "data.txt"
	output_file = "data.csv"
	alphabet    = "abcdefghijklmnopqrstuvwxyz"
	line_length = 7
	total_lines = 10_000_000
)

func main() {

	db := reader.Database{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Name:     "interview",
		Port:     5433,
	}

	ctx, conn := reader.ConnectToDB(db)

	//clear the table
	_, err := conn.Exec(ctx, "truncate table tokens")

	if err != nil {
		panic("Could not clear")
	}

	defer conn.Close(ctx)

	generator.GenerateRandomStrings(input_file, []byte(alphabet), line_length, total_lines)
	reader.Read(input_file, output_file, ctx, conn)

}
