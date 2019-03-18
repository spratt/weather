package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type DailyTemp struct {
	When time.Time
	Low float64
	High float64
}

func ParseYMD(year, month, day string) (time.Time, error) {
	return time.Parse("2006-1-2", fmt.Sprintf("%s-%s-%s", year, month, day))
}

func ParseYearDay(year string, yearDay int) time.Time {
	t, err := ParseYMD(year, "1", "1")
	check(err)
	return t.Add(time.Duration(yearDay * 24) * time.Hour)
}

const (
	YearYeardayFormat = 0
	JunkYearYeardayFormat = 1
	YearMonthDayFormat = 2
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkWithContext(err error, filename string, line int) {
	if err != nil {
		fmt.Printf("Error in %s on line %d\n", filename, line)
		panic(err)
	}
}

func Cleanup(f string) string {
	return strings.Trim(f, " ")
}

func main() {
	dailyTemps := [][]DailyTemp{}
	filenames := os.Args[1:]
	for _, filename := range filenames {
		dailyTemps = append(dailyTemps, []DailyTemp{})
		fmt.Println(filename)
		file, err := os.Open(filename)
		check(err)
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		check(err)
		tempIndex := -1
		for i, word := range records[0] {
			if strings.Contains(strings.ToLower(word), "temperature") {
				tempIndex = i
				break
			}
		}
		var (
			format int
			year string
		)
		if records[0][0] == "Year" {
			format = YearYeardayFormat
			year = records[1][0]
		} else if records[0][1] == "Year" {
			format = JunkYearYeardayFormat
			year = records[1][1]
		} else if records[0][0] == "year" {
			format = YearMonthDayFormat
			year = records[1][0]
		}
		for i, record := range records[1:] {
			var when time.Time
			if format == YearYeardayFormat {
				// Year, YearDay
				yearDay, err := strconv.Atoi(record[1])
				checkWithContext(err, filename, i + 1)
				when = ParseYearDay(year, yearDay)
			} else if format == JunkYearYeardayFormat {
				// junk, Year, YearDay
				yearDay, err := strconv.Atoi(record[2])
				checkWithContext(err, filename, i + 1)
				when = ParseYearDay(year, yearDay)
			} else if format == YearMonthDayFormat {
				// year, month, day
				t, err := ParseYMD(year, Cleanup(record[1]), Cleanup(record[2]))
				checkWithContext(err, filename, i + 1)
				when = t
			} else {
				panic(fmt.Errorf("Unknown format in %s on line %d", filename, (i+1)))
			}
			temp, err := strconv.ParseFloat(Cleanup(record[tempIndex]), 64)
			checkWithContext(err, filename, i + 1)
			lastDL := len(dailyTemps) - 1
			if len(dailyTemps[lastDL]) == 0 || dailyTemps[lastDL][len(dailyTemps[lastDL])-1].When != when {
				dailyTemps[lastDL] = append(dailyTemps[lastDL], DailyTemp{
					When: when,
					Low: temp,
					High: temp,
				})
			} else if temp < dailyTemps[lastDL][len(dailyTemps[lastDL])-1].Low {
				dailyTemps[lastDL][len(dailyTemps[lastDL])-1].Low = temp
			} else if temp > dailyTemps[lastDL][len(dailyTemps[lastDL])-1].High {
				dailyTemps[lastDL][len(dailyTemps[lastDL])-1].High = temp
			}
		}
	}
	for _, yearlyTemps := range dailyTemps {
		year := yearlyTemps[0].When.Year()
		var (
			yearlyLow float64 = math.Inf(+1.0)
			yearlyHigh float64 = math.Inf(-1.0)
			
			start time.Time
			count = 0
			runs = []time.Time{}
		)
		for _, dailyTemp := range yearlyTemps {
			if dailyTemp.Low < yearlyLow {
				yearlyLow = dailyTemp.Low
			}
			if dailyTemp.High > yearlyHigh {
				yearlyHigh = dailyTemp.High
			}
			if dailyTemp.High < 7 {
				if count > 90 {
					runs = append(runs, start)
				}
				count = 0
			} else {
				if count == 0 {
					start = dailyTemp.When
				}
				count++
			}
		}
		fmt.Printf("%d,%f,%f,%+v\n", year, yearlyLow, yearlyHigh, runs)
	}
}
