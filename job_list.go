package bkp

import (
	"github.com/rs/zerolog/log"
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
		log.Debug().Str("job name", job.Name).Msg("checking if job is relevant")
		if job.Target == nil {
			log.Fatal().Str("job name", job.Name).Msg("job does not have a target")
			continue
		}
		if !job.IsRelevant() {
			log.Debug().Str("job name", job.Name).Msg("job is NOT relevant")
			continue
		}
		if !job.Target.IsReady() {
			log.Debug().Str("job name", job.Name).Any("job target", job.Target).Msg("job target is NOT ready")
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
