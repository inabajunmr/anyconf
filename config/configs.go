package config

import (
	"bufio"
	"errors"
	"io"
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

func ReadConfigs() (*AnyConfConfigs, error) {
	r := &AnyConfConfigs{children: map[string]*AnyConfConfigs{}}
	staticConf := readStaticConfig()
	r, _ = readConfig(staticConf, r)
	localConf, err := readLocalConfig()
	if err != nil && err != io.EOF {
		return r, nil
	}
	r, _ = readConfig(localConf, r)

	return r, nil
}

func readConfig(rawConf string, conf *AnyConfConfigs) (*AnyConfConfigs, error) {
	scanner := bufio.NewScanner(strings.NewReader(rawConf))

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		sline := strings.Split(line, " ")
		if len(sline) != 2 {
			return &AnyConfConfigs{}, errors.New("static/configs.txt is something wrong")
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
	conf, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(conf, ".anyconf", "configs.txt")
	bytes, err := os.ReadFile(path)
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
	contents, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	return string(contents)
}

// Read returns children configs or target config file path
func (c AnyConfConfigs) Read(key string) (*AnyConfConfigs, error) {
	// TODO partial match?
	if c.children[key] == nil {
		return nil, errors.New("no config is matched")
	}

	return c.children[key], nil
}

type DisplayOption struct {
	Display string
	Key     string
}

func (c AnyConfConfigs) NextKeysWithDisplay() ([]string, map[string]string) {
	var displays []string
	displayToKey := make(map[string]string)
	
	for k, v := range c.children {
		var display string
		// Check if the config file exists and add emoji prefix
		if v.TargetConfigPath != "" {
			path := GetPath(v.TargetConfigPath)
			if _, err := os.Stat(path); err == nil {
				// File exists
				display = "✅ " + k + " (" + v.TargetConfigPath + ")"
			} else {
				// File doesn't exist
				display = "❌ " + k + " (" + v.TargetConfigPath + ")"
			}
		} else {
			// Directory/category (no target config path)
			display = "📁 " + k
		}
		
		displays = append(displays, display)
		displayToKey[display] = k
	}
	return displays, displayToKey
}

func (c AnyConfConfigs) NextKeys() []string {
	var keys []string
	for k, _ := range c.children {
		keys = append(keys, k)
	}
	return keys
}

func GetPath(configPath string) string {
	h, _ := os.UserHomeDir()
	return strings.Replace(configPath, "~", h, 1)
}
