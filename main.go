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

	var db reader.Database = reader.Database{
		Username: "postgres",
		Password: "4eIyCpDzAPumf7WUwixo",
		Host:     "localhost",
		Name:     "interview",
		Port:     5432,
	}

	generator.GenerateRandomStrings(input_file, []byte(alphabet), line_length, total_lines)
	reader.Read(input_file, output_file, db)
}
