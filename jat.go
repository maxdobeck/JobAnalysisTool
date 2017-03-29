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
	"time"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

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
	outputFile.WriteString(fmt.Sprintf("\r\n"))
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
	for job := range allJobs {
		curJob := allJobs[job]
		if curJob.status == "started" {
			unfinishedJobs[curJob.jobID] = curJob
		} else if curJob.status == "completed" {
			delete(unfinishedJobs, curJob.jobID)
		}
	}
	return unfinishedJobs
}

func printSliceOfJobs(allJobs []job, fileName string) {
	fileName = fmt.Sprintf("%s %s.txt",fileName,generateTimeStamp())
	outputFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	check(err)
	defer outputFile.Close()

	for job := range allJobs {
		outputFile.WriteString(fmt.Sprintf("Job ID: %s\r\n", allJobs[job].jobID))
		outputFile.WriteString(fmt.Sprintf("AUID: %s\r\n", allJobs[job].AUID))
		outputFile.WriteString(fmt.Sprintf("Line#: %d\r\n", allJobs[job].lineNumber))
		outputFile.WriteString(fmt.Sprintf("Status: %s\r\n", allJobs[job].status))
		outputFile.WriteString(fmt.Sprintf("\r\n"))
	}
	fmt.Printf("Wrote all all jobs to %s.\n", fileName)
}

func printMapOfJobs(unfinishedJobs map[string]job, fileName string) {
	fileName = fmt.Sprintf("%s %s.txt",fileName,generateTimeStamp())
	outputFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	check(err)
	defer outputFile.Close()

	for _, v := range unfinishedJobs {
		outputFile.WriteString(fmt.Sprintf("Job ID: %s\r\n", v.jobID))
		outputFile.WriteString(fmt.Sprintf("AUID: %s\r\n", v.AUID))
		outputFile.WriteString(fmt.Sprintf("Line#: %d\r\n", v.lineNumber))
		outputFile.WriteString(fmt.Sprintf("Status: %s\r\n", v.status))
		outputFile.WriteString(fmt.Sprintf("\r\n"))
	}
	fmt.Printf("Wrote all unfinishedJobs to %s.\n", fileName)
}

func generateTimeStamp () string {
	now := time.Now()
	localTimeStamp := fmt.Sprintf("%d-%d-%d %d", now.Month(), now.Day(), now.Year(), now.Nanosecond())
	return localTimeStamp
}

func main() {
	allJobs := getAllJobs("tests/testSmallBasic.log")
	printSliceOfJobs(allJobs, "output/allJobs")
	unfinishedJobs := getUnfinishedJobs(allJobs)
	printMapOfJobs(unfinishedJobs, "output/unfinishedJobs")
}