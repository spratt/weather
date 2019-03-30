package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	PlotFile   = "weather.png"
	PlotWidth  = 16*vg.Inch
	PlotHeight = 8*vg.Inch
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
	yearlyHighs := make(plotter.XYs, 366)
	yearlyLows := make(plotter.XYs, 366)
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

	// Plot
	{
		p, err := plot.New()
		check(err)

		p.Title.Text = "Waterloo Temperature"
		p.X.Label.Text = "Day"
		p.Y.Label.Text = "Temperature (Celsius)"
		p.X.Tick.Marker = MonthTicks{}
		p.Y.Tick.Marker = TemperatureTicks{}

		check(plotutil.AddLinePoints(p,
			"High", yearlyHighs,
			"Low", yearlyLows,
		))
		// Add a horizontal lines to demarcate 7 & 10 degrees
		for i, v := range []int{7, 10} {
			points := make(plotter.XYs, 366)
			for j := 0; j < 366; j++ {
				points[j].X = float64(j)
				points[j].Y = float64(v)
			}
			l, err := plotter.NewLine(points)
			check(err)
			l.LineStyle.Color = plotutil.SoftColors[i + 2]
			p.Add(l)
			p.Legend.Add(strconv.Itoa(v), l)
		}
		check(p.Save(PlotWidth, PlotHeight, PlotFile))
	}
}
