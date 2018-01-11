package bkp

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
	return j.Jobs
}
