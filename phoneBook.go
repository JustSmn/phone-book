package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

var CSVFILE = "./tmp/csv.data"
var data = []Entry{}
var index map[string]int

func readCSVFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	// CSV file read all at once
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		temp := Entry{
			Name:       line[0],
			Surname:    line[1],
			Tel:        line[2],
			LastAccess: line[3],
		}
		// Storing to global variable
		data = append(data, temp)
	}

	return nil
}
func saveCSVFile(filepath string) error {
	csvfile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer csvfile.Close()

	csvwriter := csv.NewWriter(csvfile)
	for _, row := range data {
		temp := []string{row.Name, row.Surname, row.Tel, row.LastAccess}
		_ = csvwriter.Write(temp)
	}
	csvwriter.Flush()
	return nil
}

func createIndex() error {
	index = make(map[string]int)
	for i, k := range data {
		key := k.Tel
		index[key] = i
	}
	return nil
}

// Initialized by the user – returns a pointer
// If it returns nil, there was an error
func initS(N, S, T string) *Entry {
	// Both of them should have a value
	if T == "" || S == "" {
		return nil
	}
	// Give LastAccess a value
	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Entry{Name: N, Surname: S, Tel: T, LastAccess: LastAccess}
}

func insert(pS *Entry) error {
	// If it already exists, do not add it
	_, ok := index[(*pS).Tel]
	if ok {
		return fmt.Errorf("%s already exists", pS.Tel)
	}
	data = append(data, *pS)
	// Update the index
	_ = createIndex()

	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}
	return nil
}

func deleteEntry(key string) error {
	i, ok := index[key]
	if !ok {
		return fmt.Errorf("%s cannot be found!", key)
	}
	data = append(data[:i], data[i+1:]...)
	// Update the index - key does not exist any more
	delete(index, key)

	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}
	return nil
}

func search(key string) *Entry {
	i, ok := index[key]
	if !ok {
		return nil
	}
	data[i].LastAccess = strconv.FormatInt(time.Now().Unix(), 10)
	_ = saveCSVFile(CSVFILE)
	return &data[i]
}
func list() {
	for _, v := range data {
		fmt.Println(v)
	}
}
func matchTel(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`\d+$`)
	return re.Match(t)
}

func main() {
	fmt.Println("Hello World")
}
