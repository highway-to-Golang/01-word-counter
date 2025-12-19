package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
)

func main() {
	filename := flag.String("file", "", "")
	top := flag.Int("top", 10, "")
	flag.Parse()

	if filename == nil || top == nil || *filename == "" {
		slog.Info("Usage: word-counter -file <filename> [-top N]")
		os.Exit(1)
	}

	file, err := os.Open(*filename)
	if err != nil {
		slog.Error("Error opening file", "err", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error("Error closing file", "err", err)
		}
	}(file)

	wordCounts := make(map[string]int)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		for _, word := range strings.Fields(scanner.Text()) {
			clean := strings.ToLower(strings.TrimFunc(word, func(r rune) bool {
				return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
			}))
			if clean != "" {
				wordCounts[clean]++
			}
		}
	}

	var words []struct {
		word  string
		count int
	}

	for word, count := range wordCounts {
		words = append(words, struct {
			word  string
			count int
		}{word, count})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].count > words[j].count
	})

	if len(words) > *top {
		words = words[:*top]
	}

	for _, w := range words {
		slog.Info(fmt.Sprintf("%s: %d\n", w.word, w.count))
	}
}
