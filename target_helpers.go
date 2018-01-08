package bkp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

var targets map[string]*Target

func TargetByName(name string, sourceDirs []string) *Target {
	if targets == nil {
		targets = AllTargets(sourceDirs)
	}

	if t, ok := targets[name]; ok {
		return t
	}
	return nil
}

func AllTargets(sourceDirs []string) map[string]*Target {
	targets := make(map[string]*Target)
	for _, sourceDir := range sourceDirs {
		for _, target := range ScanTargetDir(filepath.Join(sourceDir, "targets")) {
			targets[target.Name] = target
		}
	}

	return targets
}

func ScanTargetDir(path string) map[string]*Target {
	targets := make(map[string]*Target, 0)
	var errParse error

	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		var target *Target
		switch {
		case strings.HasSuffix(path, ".yml"):
			target, errParse = ParseTargetFromYmlFile(path)
		case strings.HasSuffix(path, ".json"):
			target, errParse = ParseTargetFromJSONFile(path)
		default:
			return nil
		}
		if errParse != nil {
			return nil
		}
		targets[target.Name] = target
		return nil
	})
	return targets
}

func ParseTargetFromJSONFile(filename string) (*Target, error) {
	var target *Target
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, target)
	if err != nil {
		return nil, err
	}

	target.Filename = filename

	return target, nil
}

func ParseTargetFromYmlFile(filename string) (*Target, error) {
	var target Target
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &target)
	if err != nil {
		fmt.Println("parsing error", err)
		return nil, err
	}

	target.Filename = filename

	return &target, nil
}
