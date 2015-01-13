package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i += 1
	}
	sort.Sort(p)
	return p
}

func parseLogLine(line string) (string, int) {
	strArray := strings.Split(line, ",")
	if len(strArray) < 4 {
		return "", 0
	}

	pageCount, _ := strconv.Atoi(strArray[2])
	copyCount, _ := strconv.Atoi(strArray[3])
	return strArray[1], pageCount * copyCount

}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: papercut_charts [csv] [limit]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var printStats map[string]int
	var limit int

	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Use -help for usage.")
		os.Exit(1)
	}

	if len(args) > 1 {
		l, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("Not a number!")
		}
		limit = l
	}

	file, err := os.Open(args[0])

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	printStats = make(map[string]int)
	for {
		line, isPrefix, err := reader.ReadLine()
		lineStr := string(line)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		if isPrefix {
			log.Fatalln("Line too long!")
		}
		if strings.Contains(lineStr, "PaperCut") {
			continue
		}
		if strings.Contains(lineStr, "Time,User") {
			continue
		}

		user, pageCount := parseLogLine(lineStr)
		printStats[user] += pageCount
	}
	p := sortMapByValue(printStats)

	for i, v := range p {
		if limit > 0 && i == limit {
			break
		}
		fmt.Printf("%2d.) %s: \t %d\n", i+1, v.Key, v.Value)
	}
}
