package main

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	"log"
	"math"
	"os/exec"
	"status/conf"
	"status/xtitle"
	"time"
)

// Hard-coded default options
var (
	sym_bat        = "|"
	color_accent   = "#a43052"
	color_bat_low  = "#ef1010"
	color_bat_med  = "#efef10"
	color_bat_high = "#10ef10"
	color_bat_none = "#202020"
)

type BatteryState struct {
	AcStatus int
	Percent  int
	Status   int
}

var options conf.Options

// Set the fg color of text
func setColor(input string, color string) string {
	return "%{F" + color + "}" + input + "%{F-}"
}

// get battery state
func getBattery() BatteryState {
	state := new(BatteryState)

	cmd_out, err := exec.Command("apm", "-alb").Output()
	if err != nil {
		panic(err)
	}
	fmt.Sscanf(string(cmd_out), "%d\n%d\n%d\n", &state.AcStatus, &state.Percent, &state.Status)
	return *state

}

func newBtn(text string, cmd string) string {
	return fmt.Sprintf("%%{A:%s:}%s%%{A}", cmd, text)
}

func batStat() string {
	var output string
	var barColor string = "#ffffff" // TODO ????
	var scaledValue float64

	state := getBattery()

	scaledValue = math.Ceil(float64(state.Percent) / 10)
	//scaledValue = math.Floor(100 / 10)

	for i := 10.0; i > 0; i-- {

		if i >= 6 {
			barColor = color_bat_high
		} else if i >= 3 {
			barColor = color_bat_med
		} else {
			barColor = color_bat_low
		}

		if i < scaledValue {
			output = fmt.Sprintf("%s%s", setColor(sym_bat, barColor), output)
		} else {
			output = fmt.Sprintf("%s%s", setColor(sym_bat, "#1f1f1f"), output)
		}
	}

	return output
}

func clock() string {
	dt := time.Now()
	return dt.Format("1504")
}

func main() {

	options = conf.ParseConfig("./test.conf")
	for key, value := range options {
		fmt.Printf("%s >> %s\n", key, value)
	}

	// initialise connection to X
	X, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	for {
		percent := batStat()
		lockbutton := newBtn("<lock>", "xlock -mode blank")

		title, err := xtitle.Title(X)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%%{l}  %s%%{c}[%s]%%{r}<%s>  <%s>  \n", lockbutton, title, percent, clock())
		time.Sleep(1000 * time.Millisecond)

		/*fmt.Printf("%s %s%s\n", aColor("["), clock(), aColor("]"))
		t.Sleep(500 * t.Millisecond)*/

	}
}
