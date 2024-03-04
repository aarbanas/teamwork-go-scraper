package teamwork

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aarbanas/teamwork-go-scraper/config"
)

type LogTimeMetaData struct {
	taskId      string
	hours       int16
	minutes     int16
	description string
	tag         string
}

type TimeLog struct {
	userId          int
	date            string
	time            string
	isBillable      bool
	logTimeMetaData *LogTimeMetaData
}

func prepareData(projectMode bool) (*LogTimeMetaData, error) {
	var taskId string
	var hours int16
	var minutes int16
	hoursReference := "TaskId"
	if projectMode {
		hoursReference = "ProjectId"
	}

	fmt.Printf("Enter %s: \n", hoursReference)
	_, err := fmt.Scan(&taskId)
	if err != nil {
		return nil, err
	}

	fmt.Println("Enter hours and minutes (separate with space): ")
	_, hoursMinutesError := fmt.Scan(&hours, &minutes)
	if hoursMinutesError != nil {
		return nil, hoursMinutesError
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter description: ")
	description, _ := reader.ReadString('\n')

	// To remove the new line character at the end
	description = description[:len(description)-1]

	readerTag := bufio.NewReader(os.Stdin)
	fmt.Println("Enter tags separated with comma (press enter if you want to skip): ")
	tag, _ := readerTag.ReadString('\n')

	// To remove the new line character at the end
	tag = tag[:len(tag)-1]

	return &LogTimeMetaData{taskId: taskId, hours: hours, minutes: minutes, description: description, tag: tag}, nil
}

func prepareDataForRequest(workDays []time.Time, logMetadata *LogTimeMetaData, configuration *config.Config, startTime string, nonBillable bool) *[]TimeLog {
	var timeLogs []TimeLog
	userId, _ := strconv.Atoi(configuration.UserId)
	for _, workDay := range workDays {
		timeLogs = append(timeLogs, TimeLog{userId: userId, date: convertDateToString(workDay), time: startTime, isBillable: !nonBillable, logTimeMetaData: logMetadata})
	}

	return &timeLogs
}

func LogHours(startDate, endDate, startTime string, projectMode, includeCroHolidays, nonBillable bool, configuration config.Config) {
	logMetadata, err := prepareData(projectMode)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	workDays, err := prepareWorkDays(startDate, endDate, includeCroHolidays)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	timeLogs := prepareDataForRequest(*workDays, logMetadata, &configuration, startTime, nonBillable)
	if len(*timeLogs) < 1 {
		fmt.Println("There are no time logs")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	for _, timeLog := range *timeLogs {
		wg.Add(1)
		go func(timeLog TimeLog) {
			date, err := time.Parse("20060102", timeLog.date)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				return
			}
			formattedDate := date.Format(time.DateOnly)

			_, errResponse := postTimeLogs(timeLog, projectMode, configuration)
			if errResponse != nil {
				fmt.Printf("Error sending request for date: %s\n", formattedDate)
				fmt.Printf("Error %v", errResponse)
			} else {
				fmt.Printf("Successfully logged time for date: %s\n", formattedDate)
			}
			wg.Done()
		}(timeLog)
	}
	wg.Wait()
}
