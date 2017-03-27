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
	lineNumber int
	jobID string
	AUID string
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
	outputFile, err := os.Create(fileName)
	check(err)
	defer outputFile.Close()

	outputFile.WriteString(fmt.Sprintf("Line#: %d\n",thisJob.lineNumber))
	outputFile.WriteString(fmt.Sprintf("Job ID: %s\n", thisJob.jobID))
	outputFile.WriteString(fmt.Sprintf("AUID: %s\n", thisJob.AUID))
	outputFile.WriteString(fmt.Sprintf("Status: %s", thisJob.status))
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
			//fmt.Printf("#%d %s\n",lineNumber, logEntry)
		}
		lineNumber++
	}
	return allJobs
}

func main() {
	jobs := getAllJobs("test job log file.log")
	jobs[0].printJobToFile("test.txt")
	fmt.Println(jobs)
}