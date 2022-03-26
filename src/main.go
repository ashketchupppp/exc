package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var default_config_path = filepath.Join("./", "exc.config.json")

type Executor interface {
	run(args ...string) int
}

type Command struct {
	Cmd string
}

func (command *Command) run(cmd string) {
	var cmdString = strings.Split(cmd, " ")
	proc := exec.Command(cmdString[0], cmdString[1:]...)

	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	err := proc.Start()
	if err != nil {
		log.Fatal(err)
	}
}

type CmdFlags struct {
	cmd         string
	config_path string
}

type Config struct {
	Tasks []struct {
		Name string `json:"name"`
		Cmd  string `json:"cmd"`
	} `json:"tasks"`
}

func (config *Config) find(cmdname string) string {
	var i int
	for i = 0; i < len(config.Tasks); i++ {
		if config.Tasks[i].Name == cmdname {
			return config.Tasks[i].Cmd
		}
	}
	log.Fatalf("Command not found")
	return ""
}

func parse_flags(args ...string) CmdFlags {
	var config_path = default_config_path
	var cmd string

	flag.StringVar(&config_path, "c", default_config_path, "Path to config file")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	cmd = flag.Arg(0)

	return CmdFlags{
		cmd,
		config_path,
	}
}

func read_config(path string) Config {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config not found: ", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Unable to parse config: ", err)
	}
	return config
}

func main() {
	var cmdFlags = parse_flags(os.Args...)
	var config = read_config(cmdFlags.config_path)

	cmd := Command{config.find(cmdFlags.cmd)}
	cmd.run(cmdFlags.cmd)
}
