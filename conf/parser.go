package conf

import (
	"bufio"
	"fmt"
	"os"
	//"strings"
)

/*type Options struct {
	bar_h            int
	bar_w            int
	color_foreground string
	color_background string
	color_accent     string
	color_bat_low    string
	color_bat_med    string
	color_bat_high   string
	color_bat_none   string
}*/

type Options map[string]string

//var colors map[string]string

func ParseConfig(filePath string) Options {

	options := make(Options)

	confFd, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer confFd.Close()

	scanner := bufio.NewScanner(confFd)
	for scanner.Scan() {
		var opt string
		var arg string

		line := scanner.Text()
		if line != "" {
			fmt.Sscanf(line, "%s = %s", &opt, &arg)
			options[opt] = arg
		}
	}

	return options
}
