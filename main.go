package main

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	//"github.com/BurntSushi/xgbutil"
	"math"
	"os/exec"
	"status/conf"
	"status/xtitle"
	t "time"
)

//type apm_state int
//type apm_perfmode int

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

type Options map[string]string

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
	dt := t.Now()
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
		fmt.Printf("Error| %s\n", err)
	}

	// initialise Xutil
	/*Xutil, err := xgbutil.NewConnXgb(X)
	if err != nil {
		fmt.Printf("Error| %s\n", err)
	}*/

	for {
		percent := batStat()
		lockbutton := newBtn("<lock>", "xlock -mode blank")
		title := xtitle.Title(X)
		fmt.Printf("%%{l}  %s%%{c}[%s]%%{r}%s  %s  \n", lockbutton, title, percent, clock())
		t.Sleep(1000 * t.Millisecond)

		/*fmt.Printf("%s %s%s\n", aColor("["), clock(), aColor("]"))
		t.Sleep(500 * t.Millisecond)*/

	}
}
