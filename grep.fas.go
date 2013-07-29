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
	"strconv"
	"strings"
)

var printDescr = flag.Bool("d", true, "Prints description lines")
var printSeq = flag.Bool("s", true, "Prints sequences lines")
var printFirst = flag.Bool("1", false, "Match the first sequence only")
var getNumberedSeq = flag.Bool("n", false, "Match the numbered sequences instead of matching strings, separated by commas (,)")

func main() {

	flag.Parse()

	s := ""
	numberedSeq := make(map[int]bool)
	seqCounter := 0

	// We need to check if we have enough params given
	if flag.NArg() == 0 {
		log.Fatal("Fatal error: Search pattern is needed.")
		os.Exit(-2)
	} else {
		s = flag.Args()[0]
		if *getNumberedSeq {
			numberedSeqStr := strings.Split(s, ",")
			for _, val := range numberedSeqStr {
				v, err := strconv.Atoi(val)
				if err != nil {
					log.Fatal("You have specified an invalid sequence #.")
					os.Exit(-2)
				}
				numberedSeq[v] = true
			}
		}
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

			seqMatched := false

			if !*getNumberedSeq {
				matched, err := regexp.MatchString(s, scanner.Text())
				if err != nil {
					log.Fatal(err)
					os.Exit(-1)
				}
				seqMatched = matched
			} else {
				seqMatched = numberedSeq[seqCounter]
			}

			prevMatched = seqMatched

			if prevMatched && (*printDescr) {
				fmt.Println(scanner.Text())
			}

			seqCounter++

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
