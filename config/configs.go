package config

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	staticConf := readStaticConfig()
	r, _ = readConfig(staticConf, r)
	localConf, err := readLocalConfig()
	if err != nil {
		return &r, nil
	}
	r, _ = readConfig(localConf, r)

	return &r, nil
}

func readConfig(rawConf string, conf AnyConfConfigs) (AnyConfConfigs, error) {
	scanner := bufio.NewScanner(strings.NewReader(rawConf))

	for scanner.Scan() {
		line := scanner.Text()
		sline := strings.Split(line, " ")
		if len(sline) != 2 {
			return AnyConfConfigs{}, errors.New("static/configs.text is something wrong.")
		}
		key := sline[0]
		configPath := sline[1]

		skey := strings.Split(key, "/")

		tc := conf.children
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

	return conf, nil
}

func readLocalConfig() (string, error) {
	conf, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(conf, ".anyconf", "configs.txt")
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil

}

func readStaticConfig() string {
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

func (c AnyConfConfigs) NextKeys() []string {
	var keys []string
	for k, _ := range c.children {
		keys = append(keys, k)
	}
	return keys
}
