package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	TireThreshold = 7
	PlotFile      = "weather.png"
	PlotWidth     = 16*vg.Inch
	PlotHeight    = 8*vg.Inch
	DaysInAYear   = 366
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

type TemperatureTicks struct{}
func (TemperatureTicks) Ticks(min, max float64) []plot.Tick {
	return []plot.Tick{
		plot.Tick{
			Value: min,
			Label: fmt.Sprintf("%.1f", min),
		},
		plot.Tick{
			Value: -10.0,
			Label: "-10",
		},
		plot.Tick{
			Value: 0.0,
			Label: "0",
		},
		plot.Tick{
			Value: 10.0,
			Label: "10",
		},
		plot.Tick{
			Value: 20.0,
			Label: "20",
		},
		plot.Tick{
			Value: 30.0,
			Label: "30",
		},
		plot.Tick{
			Value: max,
			Label: fmt.Sprintf("%.1f", max),
		},
	}
}

type MonthTicks struct{}
func (MonthTicks) Ticks(_, _ float64) []plot.Tick {
	return []plot.Tick{
		plot.Tick{
			Value: 0.0,
			Label: "January",
		},
		plot.Tick{
			Value: 32.0,
			Label: "February",
		},
		plot.Tick{
			Value: 60.0,
			Label: "March",
		},
		plot.Tick{
			Value: 91.0,
			Label: "April",
		},
		plot.Tick{
			Value: 121.0,
			Label: "May",
		},
		plot.Tick{
			Value: 152.0,
			Label: "June",
		},
		plot.Tick{
			Value: 182.0,
			Label: "July",
		},
		plot.Tick{
			Value: 213.0,
			Label: "August",
		},
		plot.Tick{
			Value: 244.0,
			Label: "September",
		},
		plot.Tick{
			Value: 274.0,
			Label: "October",
		},
		plot.Tick{
			Value: 305.0,
			Label: "November",
		},
		plot.Tick{
			Value: 335.0,
			Label: "December",
		},
	}
}

func main() {
	dailyTemps := [][]DailyTemp{}
	filenames := os.Args[1:]
	for _, filename := range filenames {
		dailyTemps = append(dailyTemps, []DailyTemp{})
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

	// Compute the averages
	yearlyHighs := make(plotter.XYs, DaysInAYear)
	yearlyLows := make(plotter.XYs, DaysInAYear)
	for _, yearlyTemps := range dailyTemps {
		for j, dailyTemp := range yearlyTemps {
			if yearlyHighs[j].X == 0 {
				yearlyHighs[j].X = float64(j)
				yearlyLows[j].X = float64(j)
				yearlyHighs[j].Y = dailyTemp.High
				yearlyLows[j].Y = dailyTemp.Low
			} else {
				// Note that we're using a recent-biased average here,
				// which is lazy but may also be more interesting than
				// an evenly-weighted average because we care more
				// about how the temperatures have been recently.
				yearlyHighs[j].Y += dailyTemp.High
				yearlyHighs[j].Y /= 2
				yearlyLows[j].Y += dailyTemp.Low
				yearlyLows[j].Y /= 2
			}
		}
	}

	// Find the longest string of days during which the high/low is at least 7
	type DayCount struct {
		Day int
		Count int
	}
	LowAboveCounts := []DayCount{}
	LowCurrentCount := DayCount{
		Day: -1,
	}
	HighAboveCounts := []DayCount{}
	HighCurrentCount := DayCount{
		Day: -1,
	}
	maxTemp := -999.0
	minTemp := 999.0
	for i := 0; i < DaysInAYear; i++ {
		if yearlyHighs[i].Y > maxTemp {
			maxTemp = yearlyHighs[i].Y
		}
		if yearlyLows[i].Y < minTemp {
			minTemp = yearlyLows[i].Y
		}
		if yearlyHighs[i].Y >= float64(TireThreshold) {
			if HighCurrentCount.Day == -1 {
				HighCurrentCount = DayCount{
					Day: i,
					Count: 0,
				}
			}
			HighCurrentCount.Count++
		} else if HighCurrentCount.Day > -1 {
			HighAboveCounts = append(HighAboveCounts, HighCurrentCount)
			HighCurrentCount.Day = -1
		}
		if yearlyLows[i].Y >= float64(TireThreshold) {
			if LowCurrentCount.Day == -1 {
				LowCurrentCount = DayCount{
					Day: i,
					Count: 0,
				}
			}
			LowCurrentCount.Count++
		} else if LowCurrentCount.Day > -1 {
			LowAboveCounts = append(LowAboveCounts, LowCurrentCount)
			LowCurrentCount.Day = -1
		}
	}
	sort.Slice(HighAboveCounts, func(i, j int) bool {
		return HighAboveCounts[i].Count > HighAboveCounts[j].Count
	})
	sort.Slice(LowAboveCounts, func(i, j int) bool {
		return LowAboveCounts[i].Count > LowAboveCounts[j].Count
	})
	firstHighDay := ParseYearDay("2019", HighAboveCounts[0].Day)
	firstLowDay := ParseYearDay("2019", LowAboveCounts[0].Day)

	// Plot
	{
		p, err := plot.New()
		check(err)

		p.Title.Text = "Waterloo Temperature"
		p.X.Label.Text = " "
		p.Y.Label.Text = "Temperature (Celsius)"
		p.X.Tick.Marker = MonthTicks{}
		p.Y.Tick.Marker = TemperatureTicks{}

		check(plotutil.AddLinePoints(p,
			"High", yearlyHighs,
			"Low", yearlyLows,
			fmt.Sprintf("%d", TireThreshold), plotter.XYs{
				plotter.XY{
					X: 0.0,
					Y: float64(TireThreshold),
				},
				plotter.XY{
					X: 366.0,
					Y: float64(TireThreshold),
				},
			},
			fmt.Sprintf("High consistently above %d after %s %d",
				TireThreshold, firstHighDay.Month().String(), firstHighDay.Day()),
			plotter.XYs{
				plotter.XY{
					X: float64(HighAboveCounts[0].Day),
					Y: minTemp,
				},
				plotter.XY{
					X: float64(HighAboveCounts[0].Day),
					Y: maxTemp,
				},
			},
			fmt.Sprintf("Low consistently above %d after %s %d",
				TireThreshold, firstLowDay.Month().String(), firstLowDay.Day()),
			plotter.XYs{
				plotter.XY{
					X: float64(LowAboveCounts[0].Day),
					Y: minTemp,
				},
				plotter.XY{
					X: float64(LowAboveCounts[0].Day),
					Y: maxTemp,
				},
			},
		))
		check(p.Save(PlotWidth, PlotHeight, PlotFile))
	}
}
