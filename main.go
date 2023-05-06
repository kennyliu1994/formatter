package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	file, err := os.ReadFile("assets/source.sql")
	check(err)
	sql := string(file)
	check(err)
	sql = formatAll(sql)
	lines := strings.Split(sql, "\n")
	words := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "--") {
			words = append(words, line)
			continue
		}
		if line == "" {
			continue
		}
		words = append(words, strings.Split(line, " ")...)
	}

	result := ""
	has1Equals1 := false
	addNewLine := []string{"select", "and", "order"}
	skipEndOnce := false
	if contains(words, "1=1") {
		has1Equals1 = true
	}
	for _, word := range words {
		fmt.Printf("word: %v\n", word)
		if word == "1=1" {
			continue
		}
		if strings.HasPrefix(word, "--") {
			result += "\n" + word
			continue
		}
		if contains(addNewLine, strings.ToLower(word)) {
			if skipEndOnce {
				skipEndOnce = false
			} else {
				result += "\n"
			}
			result += strings.ToLower(word) + " "
			continue
		}
		if strings.Contains(word, ",") {
			cols := strings.Split(word, ",")
			resultCols := []string{}
			for _, col := range cols {
				if col != "" {
					resultCols = append(resultCols, formatCase(col))
				}
			}
			result += strings.Join(resultCols, ",") + " "
			continue
		}
		if strings.Contains(word, ".") && !strings.Contains(word, "=") {
			result += formatCase(word) + " "
			continue
		}
		if word == "a" {
			result += word
			continue
		}
		if strings.ToLower(word) == "where" {
			result += "\n" + strings.ToLower(word) + " 1 = 1" + "\n"
			if !has1Equals1 {
				result += "and "
			}
			continue
		}
		if strings.Contains(word, "=") {
			left := strings.Split(word, "=")[0]
			right := strings.Split(word, "=")[1]
			result += formatCase(left) + " = " + right
			continue
		}
		if strings.ToLower(word) == "between" {
			skipEndOnce = true
			result += strings.ToLower(word) + " "
			continue
		}
		result += strings.ToLower(word) + " "
		continue
	}

	fmt.Printf("%v\n", result)
	os.WriteFile("assets/result_"+time.Now().Format("20060102_150405")+".sql", []byte(result), 0644)
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func formatCase(col string) string {
	left := strings.Split(col, ".")[0]
	right := strings.Split(col, ".")[1]
	return strings.ToLower(left) + "." + strings.ToUpper(right)
}

func formatAll(sql string) string {
	re := regexp.MustCompile(` {2,}`)
	sql = re.ReplaceAllString(sql, " ")
	sql = strings.ReplaceAll(sql, " ,", ",")
	sql = strings.ReplaceAll(sql, ", ", ",")
	sql = strings.ReplaceAll(sql, " =", "=")
	sql = strings.ReplaceAll(sql, "= ", "=")
	sql = strings.ReplaceAll(sql, " ;", ";")
	return sql
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
