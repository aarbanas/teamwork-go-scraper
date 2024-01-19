# Teamwork scraper

Simple Go application for calculating hours logged in Teamwork. 

Application must be used through command line and can calculate hours based on: 
1. `action` - "tag" or "projectId" 
2. `value` - for specified action (e.g. 1234)
3. `startDate`
4. `endDate` 


For more info run in the terminal next command.
```bash
$ go run ./cmd -help
```

## Prerequisites
You must have a running go version locally. Run the next command to check if you have a running go version.
```bash
$ go version
```

## Usage
Just follow the next steps.
1. Prepare teamwork [API key](https://apidocs.teamwork.com/docs/teamwork/df5a63302d729-getting-started-with-the-teamwork-com-api). 
2. Export necessary environment variables (check `.envrc-example`)
   1. For the `USER_ID` get it from Teamwork.
   2. I am using `direnv` for env variables
3. Run the script
```bash
$ go run ./cmd
```
Example above is executing the script using default values but feel free to pass your own command-line arguments.