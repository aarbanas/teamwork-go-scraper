# Teamwork scraper

Simple Go application for calculating hours logged and logging time in Teamwork.

This app can be used without any dependencies.
Download latest version from [releases](https://github.com/aarbanas/teamwork-go-scraper/releases) and add `config.json` (check the [config.json-example](https://github.com/aarbanas/teamwork-go-scraper/blob/main/config.json-example)) to the
same directory where the downloaded binary file is. Than just run on UNIX:

```bash
$ ./teamwork-go-scraper
```

For Windows users

```powershell
$ .\teamwork-go-scraper.exe
```

For help run:

```bash
$ ./teamwork-go-scraper -help

# String flags (must specify value -action=tag)
-action        Action on which hours will be calculated (tag or projectId) (default "tag")
-startDate     Must be in format YYYY-MM-DD (default "first day of current month (e.g. 2024-01-01)")
-endDate       Must be in format YYYY-MM-DD (default "today (e.g. 2024-01-15)")
-t             Start time when to log hours from (HH:mm) (default "09:00")
-value         Value for the specified action. (default "overtime")

# Boolean flags (just passed as arguments)
-c             Use for checking if there are some days where hours are not logged
-h             Use for including Croatian national holidays in the calculations
-l             Enter the logging mode (default: reading logged hours)
-p             If selected hours will be logged by project id (default: log by task id )
-n             Non billable hours in log mode (default: isBillable)
```

# Table of Contents

1. [Prerequisites Go users](#for-users-having-go-installed)
2. [Usage](#usage)
3. [Action](#action)
4. [Value](#value)
5. [Dates](#dates)
6. [Logging hours](#log-mode)
7. [Logging example](#example-for-logging-hours-january-2024)
8. [Holidays](#croatian-national-holidays)
9. [Check logged-in hours](#check-if-there-are-days-missing-hours)

## For users having Go installed

You must have a running go version locally [Go.dev](https://go.dev/doc/install).
Run the next command to check if you have a running go version.

```bash
$ go version
```

Just follow the next steps.

1. Prepare teamwork [API key](https://apidocs.teamwork.com/docs/teamwork/df5a63302d729-getting-started-with-the-teamwork-com-api).
2. Create `config.json` in the root
3. Export necessary environment variables (check `.config.json-example`)
   1. Get the `userId` from Teamwork

## Installation

To install the package just run:

```bash
go install github.com/aarbanas/teamwork-go-scraper/cmd@latest
```

## Usage

By default this application is going to calculate all overtime hours logged in teamwork. Here is a list of other options:

### Action

Application is accepting two action values: `tag`(default) and `projectId`. Based on tag script is going to calculate logged hours. (e.g. if projectId is 12345 than it will return all logged hours for selected project)

```bash
go run ./cmd -action=projectId
```

### Value

You can pass any `-value` which needs to be in the combination with `-action`. Default value is `overtime`.

```bash
go run ./cmd -action=projectId -value=123456789
```

### Dates

For dates two values can be used:

1. `-startDate` - default is first day of the current month (must be in format YYYY-MM-DD)
2. `-endDate` - default is today (must be in format YYYY-MM-DD)

```bash
go run ./cmd -startDate=2024-01-01 -endDate=2024-01-15
```

### Log mode

In case you want to log new hours, not check for the already logged ones than flag `-l` needs to be passed.

```bash
go run ./cmd -l
```

Logging hours is selecting `startDate` and `endDate` flags, so if they are not passed the default values are used.
Using log mode is prompt couple of more questions like taskId, hours/minutes and description.

Log mode by default is logging hours by `taskId`, but in case hours need to be logged on the project directly there is
a flag `-p`.

```bash
go run ./cmd -l -p
```

### Start time

To specify the start time use the flag `-t`. Time must be provided in format `HH:mm`

```bash
go run ./cmd -l -t=08:30
```

### Example for logging hours January 2024

```bash
go run ./cmd -l -p -startDate=2024-01-01 -endDate=2024-01-31
```

### Croatian national holidays

In case you want to log hours and skip Croatian national holidays use `-h` flag.

```bash
go run ./cmd -l -h -startDate=2024-01-01 -endDate=2024-01-31
```

### Check if there are days missing hours

If you want to validate if in the specified period you logged all your hours there is a flag `-c` made for that purpose.
Here you can add Croatian national holidays in the equation.

```bash
go run ./cmd -c -h -startDate=2024-01-01 -endDate=2024-01-31
```
