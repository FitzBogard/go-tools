package map_reduce

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// process big file using map-reduce concept (by ChatGPT ^^)

const (
	chunkSize = 1 << 30 // 1G
)

type element struct {
	value string
	count int
}

func main() {
	splitAndCount("input.txt")
	mergeAndFindTop10()
}

func splitAndCount(inputFile string) {
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	chunkIndex := 0
	for {
		chunkFile := fmt.Sprintf("chunk_%d.txt", chunkIndex)
		w, err := os.Create(chunkFile)
		if err != nil {
			panic(err)
		}
		wr := bufio.NewWriter(w)

		var bytesRead int
		for bytesRead < chunkSize {
			line, err := r.ReadString(',')
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			bytesRead += len(line)
			wr.WriteString(line)
		}

		wr.Flush()
		w.Close()

		go countFrequency(chunkFile, fmt.Sprintf("result_%d.txt", chunkIndex))

		chunkIndex++

		if bytesRead < chunkSize {
			break
		}
	}
}

func countFrequency(inputFile, outputFile string) {
	f, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	frequency := make(map[string]int)

	for {
		line, err := r.ReadString(',')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		value := strings.TrimSpace(line)
		frequency[value]++
	}

	f, err = os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for value, count := range frequency {
		w.WriteString(fmt.Sprintf("%s,%d\n", value, count))
	}

	w.Flush()
}

func mergeAndFindTop10() {
	files, err := filepath.Glob("result_*.txt")
	if err != nil {
		panic(err)
	}

	merged := make(map[string]int)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		r := bufio.NewReader(f)

		for {
			line, err := r.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}

			parts := strings.Split(strings.TrimSpace(line), ",")
			if len(parts) != 2 {
				continue
			}

			value := parts[0]
			count := atoi(parts[1])

			merged[value] += count
		}

		f.Close()
	}

	top10 := findTopN(merged, 10)

	for _, elem := range top10 {
		fmt.Printf("%s: %d\n", elem.value, elem.count)
	}
}

func findTopN(data map[string]int, n int) []element {
	if len(data) < n {
		n = len(data)
	}

	topN := make([]element, n)

	i := 0
	for value, count := range data {
		if i < n {
			topN[i] = element{value, count}
			i++
			continue
		}

		minIndex := 0
		for j := 1; j < n; j++ {
			if topN[j].count < topN[minIndex].count {
				minIndex = j
			}
		}

		if count > topN[minIndex].count {
			topN[minIndex] = element{value, count}
		}
	}

	sort.Slice(topN, func(i, j int) bool {
		return topN[i].count > topN[j].count
	})

	return topN
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}
