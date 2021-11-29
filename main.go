package main

import (
	"context"
	"flag"
	"interview/src/generator"
	"interview/src/reader"
	"log"
	"os"
	"runtime/pprof"
)

const (
	input_file  = "data.txt"
	output_file = "data.csv"
	alphabet    = "abcdefghijklmnopqrstuvwxyz"
	line_length = 7
	total_lines = 10_000_000
)

func main() {

	// taken from https://go.dev/blog/pprof
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	db := reader.Database{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Name:     "interview",
		Port:     5433,
	}

	conn := reader.ConnectToDB(db)

	//clear the table
	_, err := conn.Exec(context.Background(), "truncate table tokens")

	if err != nil {
		panic("Could not clear")
	}

	defer conn.Close(context.Background())

	generator.GenerateRandomStrings(input_file, []byte(alphabet), line_length, total_lines)
	reader.Read(input_file, output_file, conn)

}
