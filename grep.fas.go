package main

/*
 grep.fas
 ------------------------------------------
 Like grep, but for fasta/fastq sequences
 Search fasta file for sequences matching a certain pattern
 Huan Truong <htruong@tnhh.net> #codencoffee 06-22-2013
*/

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

var printDescr = flag.Bool("d", true, "Prints description lines")
var printSeq = flag.Bool("s", true, "Prints sequences lines")
var printFirst = flag.Bool("1", false, "Match the first sequence only")

func main() {

	flag.Parse()

	s := ""

	// We need to check if we have enough params given
	if flag.NArg() == 0 {
		log.Fatal("Fatal error: Search pattern is needed.")
		os.Exit(-2)
	} else {
		s = flag.Args()[0]
	}

	// Now get the stdin pipe
	scanner := bufio.NewScanner(os.Stdin)

	prevMatched := false

	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}

		if scanner.Text()[0] == '>' {
			// This line has the sequence name

			// If we had the first match, and the user just
			// requests to print the first line only, then quit right now
			if prevMatched && *printFirst {
				os.Exit(0)
			}

			matched, err := regexp.MatchString(s, scanner.Text())

			if err != nil {
				log.Fatal(err)
				os.Exit(-1)
			}

			prevMatched = matched

			if prevMatched && (*printDescr) {
				fmt.Println(scanner.Text())
			}

		} else {
			// This must be the sequence
			if prevMatched && (*printSeq) {
				fmt.Println(scanner.Text())
			}
		}
		//log.Printf("%b, %s", prevMatched, scanner.Text())

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
