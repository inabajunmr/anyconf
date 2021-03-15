package config

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/inabajunmr/anyconf/statik"

	"github.com/rakyll/statik/fs"
)

type AnyConfConfigs struct {
	TargetConfigPath string
	NodeName         string
	children         map[string]*AnyConfConfigs
}

func ReadConfig() (*AnyConfConfigs, error) {
	r := AnyConfConfigs{children: map[string]*AnyConfConfigs{}}
	confStr := readConfigFile()
	scanner := bufio.NewScanner(strings.NewReader(confStr))

	for scanner.Scan() {
		line := scanner.Text()
		sline := strings.Split(line, " ")
		if len(sline) != 2 {
			return nil, errors.New("static/configs.text is something wrong.")
		}
		key := sline[0]
		configPath := sline[1]

		skey := strings.Split(key, "/")

		tc := r.children
		for i, v := range skey {
			if tc[v] == nil {
				tc[v] = &AnyConfConfigs{children: map[string]*AnyConfConfigs{}, NodeName: v}
			}
			if i == len(skey)-1 {
				tc[v].TargetConfigPath = configPath
			}

			tc = tc[v].children
		}

	}

	return &r, nil
}

func readConfigFile() string {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	r, err := statikFS.Open("/configs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	return string(contents)
}

// Read returns children configs or target config file path
func (c AnyConfConfigs) Read(key string) (*AnyConfConfigs, error) {
	// TODO partial match?
	if c.children[key] == nil {
		return nil, errors.New("No config is matched.")
	}

	return c.children[key], nil
}
