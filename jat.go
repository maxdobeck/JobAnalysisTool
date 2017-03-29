/* Job Analysis Tool
** Author: Max Dobeck
** Date: 3-25-2016
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

/*type jobsInterface interface {

}*/

type job struct {
	jobID      string
	AUID       string
	lineNumber int
	status     string
}

func makeJob(lineNumber int, logEntry string) job {
	thisJobID := regexp.MustCompile("(Job \\d)")
	thisJobAUID := regexp.MustCompile("((AUID:) \\d{4,9})")
	thisJobStatus := regexp.MustCompile("completed|started")

	rawJobID := strings.Fields(thisJobID.FindString(logEntry))
	rawJobAUID := strings.Fields(thisJobAUID.FindString(logEntry))

	var newJob job
	newJob.lineNumber = lineNumber
	newJob.jobID = rawJobID[1]
	newJob.AUID = rawJobAUID[1]
	newJob.status = thisJobStatus.FindString(logEntry)
	return newJob
}

func (thisJob job) printJobToFile(fileName string) {
	outputFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	check(err)
	defer outputFile.Close()

	outputFile.WriteString(fmt.Sprintf("Job ID: %s\r\n", thisJob.jobID))
	outputFile.WriteString(fmt.Sprintf("AUID: %s\r\n", thisJob.AUID))
	outputFile.WriteString(fmt.Sprintf("Line#: %d\r\n", thisJob.lineNumber))
	outputFile.WriteString(fmt.Sprintf("Status: %s\r\n", thisJob.status))
	outputFile.WriteString(fmt.Sprintf("\r\n\r\n"))
}

func isJob(logEntry string) bool {
	matchBool, err := regexp.MatchString("(Job)", logEntry)
	check(err)
	return matchBool
}

func getAllJobs(fileName string) []job {
	allJobs := make([]job, 0)
	lineNumber := 0
	logfile, err := os.Open(fileName)
	check(err)
	defer logfile.Close()

	scanner := bufio.NewScanner(logfile)

	for scanner.Scan() {
		logEntry := scanner.Text()
		if isJob(logEntry) {
			thisJob := makeJob(lineNumber, logEntry)
			allJobs = append(allJobs, thisJob)
		}
		lineNumber++
	}
	return allJobs
}

func getUnfinishedJobs(allJobs []job) map[string]job {
	unfinishedJobs := make(map[string]job)
	for i := range allJobs {
		if allJobs[i].status == "started" {
			unfinishedJobs[allJobs[i].jobID] = allJobs[i]
		} else if allJobs[i].status == "completed" {
			delete(unfinishedJobs, allJobs[i].jobID)
		}
	}
	return unfinishedJobs
}

func printAllJobs(allJobs []job, fileName string) {
	outputFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	check(err)
	defer outputFile.Close()

	for job := range allJobs {
		outputFile.WriteString(fmt.Sprintf("Job ID: %s\r\n", allJobs[job].jobID))
		outputFile.WriteString(fmt.Sprintf("AUID: %s\r\n", allJobs[job].AUID))
		outputFile.WriteString(fmt.Sprintf("Line#: %d\r\n", allJobs[job].lineNumber))
		outputFile.WriteString(fmt.Sprintf("Status: %s\r\n", allJobs[job].status))
		outputFile.WriteString(fmt.Sprintf("\r\n\r\n"))
	}
	fmt.Printf("Wrote all jobs to %s.\n", fileName)
}
/*
func printUnfinishedJobs(unfinishedJobs map[string]job, fileName string) {
	unfinishedJobs := make(map[string]job)
	for i := range allJobs {
		if allJobs[i].status == "started" {
			unfinishedJobs[allJobs[i].jobID] = allJobs[i]
		} else if allJobs[i].status == "completed" {
			delete(unfinishedJobs, allJobs[i].jobID)
		}
	}
		unfinishedJobOutputFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	check(err)
	defer unfinishedJobOutputFile.Close()

	for _, v := range unfinishedJobs {
		unfinishedJobOutputFile.WriteString(fmt.Sprintf("Job ID: %s\r\n", v.jobID))
		unfinishedJobOutputFile.WriteString(fmt.Sprintf("AUID: %s\r\n", v.AUID))
		unfinishedJobOutputFile.WriteString(fmt.Sprintf("Line#: %d\r\n", v.lineNumber))
		unfinishedJobOutputFile.WriteString(fmt.Sprintf("Status: %s\r\n", v.status))
		unfinishedJobOutputFile.WriteString(fmt.Sprintf("\r\n\r\n"))
	}
	fmt.Printf("Wrote all unfinishedJobs to %s.\n", fileName)
}*/

func main() {
	allJobs := getAllJobs("tests/testSmallBasic.log")
	printAllJobs(allJobs, "output/allJobs.txt")
	unfinishedJobs := getUnfinishedJobs(allJobs)
	printUnfinishedJobs(unfinishedJobs, "output/unfinishedJobs.txt")
}
