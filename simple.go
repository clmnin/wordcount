package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
	"strings"
)

func main() {
	// Start Profiler
	f, err := os.Create("cpuprofile")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
		os.Exit(1)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Fprintf(os.Stderr, "could not start CPU profile: %v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	// init an object of type Map<string, int>
	counts := make(map[string]int)
	// read to the next token. The token is set as "space" scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		// get the text, lower case it and increase the count at HashMap
		word := strings.ToLower(scanner.Text())
		counts[word]++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// ordered is a <List> of type <Count> (a struct)
	var ordered []Count
	// for word, count in range counts
	for word, count := range counts {
		// append to ordered
		ordered = append(ordered, Count{word, count})
	}
	// sort the list of <struct>Count with Count.Count
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Count > ordered[j].Count
	})

	for _, count := range ordered {
		fmt.Println(string(count.Word), count.Count)
	}
	// End Profiler
	defer pprof.StopCPUProfile()
}

type Count struct {
	Word  string
	Count int
}
