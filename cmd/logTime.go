package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type LogTimeMetaData struct {
	taskId      string
	hours       int16
	minutes     int16
	description string
}

type TimeLog struct {
	userId          int
	date            string
	time            string
	isBillable      bool
	logTimeMetaData *LogTimeMetaData
}

func prepareData() (*LogTimeMetaData, error) {
	var taskId string
	var hours int16
	var minutes int16
	fmt.Println("Enter TaskId: ")
	_, taskError := fmt.Scan(&taskId)
	if taskError != nil {
		return nil, taskError
	}

	fmt.Println("Enter hours and minutes (separate with space): ")
	_, hoursMinutesError := fmt.Scan(&hours, &minutes)
	if taskError != nil {
		return nil, hoursMinutesError
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter description: ")
	description, _ := reader.ReadString('\n')

	// To remove the new line character at the end
	description = description[:len(description)-1]

	return &LogTimeMetaData{taskId: taskId, hours: hours, minutes: minutes, description: description}, nil
}

func prepareDataForRequest(workDays *[]time.Time, logMetadata *LogTimeMetaData) *[]TimeLog {
	var timeLogs []TimeLog
	TIME := "09:00:00"
	userId, _ := strconv.Atoi(os.Getenv("USER_ID"))
	for _, workDay := range *workDays {
		timeLogs = append(timeLogs, TimeLog{userId: userId, date: convertDateToString(workDay), time: TIME, isBillable: true, logTimeMetaData: logMetadata})
	}

	return &timeLogs
}

func logHours(startDate *string, endDate *string) {
	logMetadata, prepareDateErr := prepareData()
	if prepareDateErr != nil {
		fmt.Printf("Error: %s", prepareDateErr)
	}

	workingStartDate, err := convertStringFormatToDate(*startDate)
	workingEndDate, _err := convertStringFormatToDate(*endDate)
	if err != nil || _err != nil {
		os.Exit(1)
	}

	workDays := getWorkingDays(workingStartDate, workingEndDate)
	if len(*workDays) < 1 {
		fmt.Println("There are no workdays in selected period")
		os.Exit(1)
	}

	timeLogs := prepareDataForRequest(workDays, logMetadata)
	if len(*timeLogs) < 1 {
		fmt.Println("There are no time logs")
		os.Exit(1)
	}

	for _, timeLog := range *timeLogs {
		_, errResponse := postTimeLogs(&timeLog)
		if errResponse != nil {
			fmt.Printf("Error sending request for date: %s\n", timeLog.date)
			fmt.Printf("Error %v", errResponse)
		} else {
			fmt.Printf("Successfully logged time for date: %s\n", timeLog.date)
		}
	}
}
