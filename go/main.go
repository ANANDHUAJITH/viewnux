
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func tasks() {
	fmt.Printf("%7s %5s %9s  %s\n", "PID", "STATE", "RAM(MB)", "COMMAND")

	files, err := os.ReadDir("/proc")
	if err != nil {
		return
	}

	for _, file := range files {
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		path := "/proc/" + file.Name()

		stat, err := os.ReadFile(path + "/stat")
		if err != nil {
			continue
		}

		line := string(stat)
		end := strings.LastIndex(line, ")")
		if end < 0 {
			continue
		}

		name := line[strings.Index(line, "(")+1 : end]
		fields := strings.Fields(line[end+1:])
		if len(fields) < 22 {
			continue
		}

		state := fields[0]

		statm, err := os.ReadFile(path + "/statm")
		if err != nil {
			continue
		}

		mem := strings.Fields(string(statm))
		if len(mem) < 2 {
			continue
		}

		pages, err := strconv.ParseUint(mem[1], 10, 64)
		if err != nil {
			continue
		}

		ram := float64(pages) * float64(os.Getpagesize()) / 1048576

		cmdline, _ := os.ReadFile(path + "/cmdline")
		command := strings.TrimSpace(strings.ReplaceAll(
			string(cmdline), "\x00", " ",
		))

		if command == "" {
			command = "[" + name + "]"
		}

		fmt.Printf("%7d %5s %9.1f  %s\n", pid, state, ram, command)
	}
}

func main() {
	for {
		fmt.Print("\033[H\033[2J")
		tasks()
		fmt.Println("\nPress Ctrl+C to exit.")
		time.Sleep(time.Second)
	}
}
