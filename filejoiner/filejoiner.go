package filejoiner

import (
	"bufio"
	"os"
	"sync"
)

func makeWork(textFiles ...string) <-chan string {
	ch := make(chan string)
	go func() {
		for _, textFile := range textFiles {
			ch <- textFile
		}
		close(ch)
	}()
	return ch
}

func readTextFile(textFile string) []string {
	tf, err := os.Open(textFile)
	if err != nil {
		panic(err)
	}
	defer tf.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(tf)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func writeTextFile(textFile string, lines []string) {
	tf, err := os.Create(textFile)
	if err != nil {
		panic(err)
	}
	defer tf.Close()

	for _, line := range lines {
		tf.WriteString(line + "\n")
	}
}

func pipeline[I any, O any](in <-chan I, process func(I) O) <-chan O {
	out := make(chan O)
	go func() {
		for i := range in {
			out <- process(i)
		}
		close(out)
	}()
	return out
}
func JoinFilesAsync() {
	var wg sync.WaitGroup
	textFiles := makeWork("page1.txt", "page2.txt", "page3.txt")
	lines := pipeline(textFiles, readTextFile)
	var joinedLines []string

	for line := range lines {
		wg.Add(1)
		go func(line []string) {
			defer wg.Done()
			joinedLines = append(joinedLines, line...)
		}(line)

	}
	wg.Wait()
	writeTextFile("joinedasync.txt", joinedLines)

}

func JoinFilesSync() {
	pageFiles := []string{"page1.txt", "page2.txt", "page3.txt"}
	var toWriteLines []string

	for _, pageFile := range pageFiles {
		lines := readTextFile(pageFile)
		toWriteLines = append(toWriteLines, lines...)
	}
	writeTextFile("joinedsync.txt", toWriteLines)
}
