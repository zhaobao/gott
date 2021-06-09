package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "-h" {
		fmt.Println("🤜\nusage: ls -la | gott\n🤛")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	lines := make([]string, 0)
	var buf bytes.Buffer
	for j := 0; j < len(output); j++ {
		if output[j] == 10 { // break line
			lines = append(lines, buf.String())
			buf.Reset()
			continue
		}
		buf.WriteRune(output[j])
	}

	regSplit := regexp.MustCompile(`\s+`)
	tableData := make([][]string, len(lines)-1)
	for idx, line := range lines {
		if line[0] != '-' && line[0] != 'd' && line[0] != 'l' && line[0] != 'b' && line[0] != 'c' {
			continue
		}
		columns := regSplit.Split(line, 8)
		if len(columns) != 8 {
			log.Fatalln("column.not.right", columns)
		}
		tableData[idx-1] = columns
	}

	mergedData := make([][]string, len(tableData))
	for idx, line := range tableData {
		size, _ := strconv.ParseInt(line[3], 10, 64)
		mergedLine := make([]string, len(line)-2)
		mergedLine[0], mergedLine[1], mergedLine[2], mergedLine[3] = line[0], line[1], line[2], readableBytes(size)
		mergedLine[4] = strings.Join([]string{line[4], line[5], line[6]}, " ")
		mergedLine[5] = line[7]
		mergedData[idx] = mergedLine
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"permission", "links", "owner", "size", "date", "name"})
	table.AppendBulk(mergedData)
	table.Render()
}

func readableBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("B %d", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%c %.1f", "KMGTPE"[exp], float64(b)/float64(div))
}
