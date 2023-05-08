package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func get_all_files() []string {
	files, err := filepath.Glob("*.log")
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func get_tail_file(n int, file string) ([]string, []string) {
	cmd := exec.Command("tail", fmt.Sprintf("-n%s", strconv.Itoa(n)), file)
	out, _ := cmd.Output()
	finish := strings.Split(string(out), "\n")

	description := []string{}
	numbers := []string{}

	for i := 0; i < len(finish); i++ {

		all := strings.Fields(finish[i])
		if len(all) == 0 {
			continue
		}
		if all[1] == "" || all[1] == " " || all[1] == "\n" || all[1] == "\r" || all[1] == "\t" {
			continue
		}

		description = append(description, all[1])
		numbers = append(numbers, all[0])
		//fmt.Println(all[0], "\n")
	}
	for i := range description {
		description[i] = strings.TrimSpace(description[i])
	}

	for i := range numbers {
		numbers[i] = strings.TrimSpace(numbers[i])
	}

	return description, numbers
}

func main() {

	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Error! You need 2 args:\n1)Len of tail in logs\n2)Len of split group (0 if you are not need split)")
		return
	}

	len_tail, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("First arg must be INT! Exit...")
		return
	}

	len_split, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Second arg must be INT! Exit...")
		return
	}

	files := get_all_files()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new sheet.
	index, err := f.NewSheet("Perf_statistic")
	if err != nil {
		fmt.Println(err)
		return
	}

	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#A9D08E"}, Pattern: 1}, //#DDEBF7
	})
	err = f.SetCellStyle("Perf_statistic", "B2", "AI2", style)

	for i := 0; i < len(files); i++ {
		description, numbers := get_tail_file(len_tail, files[i])
		//benchmark_name := strings.Split(files[i],".") []
		//input_name := strings.Split(files[i],".") []
		for j := 0; j < len(description); j++ {
			cell, _ := excelize.ColumnNumberToName(j + 2)
			f.SetCellValue("Perf_statistic", fmt.Sprintf("%s2", cell), description[j])
			f.SetCellValue("Perf_statistic", fmt.Sprintf("A%d", i+3), files[i])
			f.SetCellValue("Perf_statistic", fmt.Sprintf("%s%d", cell, i+3), numbers[j])

		}
	}

	if len_split != 0 {
		for i := 5; i < len(files)*3; i = i + len_split + 1 {
			f.InsertRows("Perf_statistic", i, 1)
		}
	}

	//f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	//"Book1.xlsx"
	if err := f.SaveAs(fmt.Sprintf("perf_%d.xlsx", time.Now().Unix())); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successful done!")
}
