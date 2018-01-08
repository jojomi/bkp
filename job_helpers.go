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

func AllJobs(sourceDirs []string) []*Job {
	jobs := make([]*Job, 0)
	for _, sourceDir := range sourceDirs {
		for _, job := range ScanJobDir(filepath.Join(sourceDir, "jobs"), sourceDirs) {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

func ScanJobDir(path string, targetDirs []string) []*Job {
	jobs := make([]*Job, 0)
	var errParse error

	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		var job *Job
		switch {
		case strings.HasSuffix(path, ".yml"):
			job, errParse = ParseJobFromYmlFile(path, targetDirs)
		case strings.HasSuffix(path, ".json"):
			job, errParse = ParseJobFromJSONFile(path, targetDirs)
		default:
			return nil
		}
		if errParse != nil {
			return nil
		}
		jobs = append(jobs, job)
		return nil
	})
	return jobs
}

func ParseJobFromJSONFile(filename string, targetDirs []string) (*Job, error) {
	var job Job
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &job)
	if err != nil {
		return nil, err
	}

	augmentJob(&job, filename, targetDirs)

	return &job, nil
}

func ParseJobFromYmlFile(filename string, targetDirs []string) (*Job, error) {
	var job Job
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &job)
	if err != nil {
		fmt.Println("parsing error", err)
		return nil, err
	}

	augmentJob(&job, filename, targetDirs)

	return &job, nil
}

func augmentJob(job *Job, filename string, targetDirs []string) {
	job.Filename = filename
	job.Target = TargetByName(job.TargetName, targetDirs)
}
