package bkp

import (
	"sort"
	"strings"
)

type JobList struct {
	Jobs []*Job
}

func (j *JobList) Load(sourceDirs []string) {
	j.Jobs = AllJobs(sourceDirs)
}

func (j *JobList) Relevant() []*Job {
	relevantJobs := make([]*Job, 0)
	jobs := j.All()
	for _, job := range jobs {
		if !job.IsRelevant() {
			continue
		}
		if job.Target != nil && !job.Target.IsReady() {
			continue
		}
		relevantJobs = append(relevantJobs, job)
	}
	return relevantJobs
}

func (j *JobList) All() []*Job {
	// sort
	sort.Slice(j.Jobs, func(i, k int) bool {
		var (
			a = j.Jobs[i]
			b = j.Jobs[k]
		)
		if a.Weight == b.Weight {
			return strings.Compare(a.Name, b.Name) == -1
		}
		return a.Weight < b.Weight
	})
	return j.Jobs
}
