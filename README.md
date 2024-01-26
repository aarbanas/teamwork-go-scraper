# Teamwork scraper

Simple Go application for calculating hours logged and logging time in Teamwork.

For building binary run:

```bash
$ go build -o bin/run ./cmd
```

Or use directly in terminal:

```bash
$ go run ./cmd
```

For help run:

```bash
$ go run ./cmd -help
```

## Prerequisites

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

## Usage

By default this application is going to calculate all overtime hours logged in teamwork. Here is a list of other options:

### Action

Application is accepting two action values: `tag`(default) and `projectId`. Based on tag script is going to calculate logged hours. (e.g. if projectId is 12345 than it will return all logged hours for selected project)

```bash
go run ./cmd -action=projectId
```

### Value

You can pass any `value` which needs to be in the combination with `action`. Default value is `overtime`.

```bash
go run ./cmd -action=projectId -value=123456789
```

### Dates

For dates two values can be used:

1. `startDate` - default is first day of the current month (must be in format YYYY-MM-DD)
2. `endDate` - default is today (must be in format YYYY-MM-DD)

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

### Example for logging hours January 2024

```bash
go run ./cmd -l -p -startDate=2024-01-01 -endDate=2024-01-31
```
