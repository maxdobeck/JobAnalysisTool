/* Job Analysis Tool
** Author: Max Dobeck
** Date: 3-25-2016
*/

package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"regexp"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type job struct {
	jobID string
	AUID string
	lineNumber int
	status string
}

func makeJob(lineNumber int, logEntry string) job{
	thisJobID := regexp.MustCompile("(Job \\d)")
	thisJobAUID := regexp.MustCompile("((AUID:) \\d{4,9})")
	thisJobStatus := regexp.MustCompile("completed|started")

	rawJobID := strings.Fields(thisJobID.FindString(logEntry))
	rawJobAUID := strings.Fields(thisJobAUID.FindString(logEntry))

	newJob := job{lineNumber: lineNumber, jobID: rawJobID[1], AUID: rawJobAUID[1], status: thisJobStatus.FindString(logEntry)}
	return newJob
}

func (thisJob job) printJobToFile(fileName string) {
	outputFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	check(err)
	defer outputFile.Close()

	outputFile.WriteString(fmt.Sprintf("Job ID: %s\n", thisJob.jobID))
	outputFile.WriteString(fmt.Sprintf("AUID: %s\n", thisJob.AUID))
	outputFile.WriteString(fmt.Sprintf("Line#: %d\n",thisJob.lineNumber))
	outputFile.WriteString(fmt.Sprintf("Status: %s\n\n", thisJob.status))
}

func isJob(logEntry string) bool {
	matchBool,err := regexp.MatchString("(Job)", logEntry)
	check(err)
	return matchBool
}

func getAllJobs(fileName string) []job {
	allJobs := make([]job,0)
	lineNumber := 0
	logfile, err := os.Open(fileName)
	check(err)
	defer logfile.Close()

	scanner := bufio.NewScanner(logfile)
	
	for scanner.Scan() {
		logEntry := scanner.Text()
		if isJob(logEntry) {
			thisJob := makeJob(lineNumber, logEntry)
			allJobs = append(allJobs,thisJob)
		}
		lineNumber++
	}
	return allJobs
}

func printAllJobs(allJobs []job, fileName string) {
	for job := range allJobs {
		allJobs[job].printJobToFile(fileName)
	}
	fmt.Printf("Wrote all jobs to %s.\n",fileName)
}

func printUnfinishedJobs(allJobs []job, fileName string) {
	unfinishedJobs := make(map[string]job)
	for i := range allJobs {
		if allJobs[i].status == "started"{
			unfinishedJobs[allJobs[i].jobID] = allJobs[i]
		} else if allJobs[i].status == "completed" {
			delete(unfinishedJobs, allJobs[i].jobID)
		}
	}
	for _, v := range unfinishedJobs {
		v.printJobToFile(fileName)
	}
	fmt.Printf("Wrote all unfinishedJobs to %s.\n", fileName)
}

func main() {
	jobs := getAllJobs("test job log file.log")
	printAllJobs(jobs, "allJobs.txt")
	printUnfinishedJobs(jobs, "unfinishedJobs.txt")
}