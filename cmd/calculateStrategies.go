package main

import "strconv"

// Strategy defines the interface for different strategies.
type Strategy interface {
	CalculateHours(result *Response, value string) float64
}

/* ------------------- Tag concrete implementation------------------- */
type CalculateByTag struct{}

func (a *CalculateByTag) CalculateHours(result *Response, value string) float64 {
	totalOvertimeHours := 0.0
	for _, rec := range result.TimeEntries {
		for _, tag := range rec.Tags {
			if tag.Name == value {
				totalOvertimeHours += rec.HoursDecimal
			}
		}
	}

	return totalOvertimeHours
}

/* ------------------- ProjectId concrete implementation------------------- */
type CalculateByProjectId struct{}

func (s *CalculateByProjectId) CalculateHours(result *Response, value string) float64 {
	totalOvertimeHours := 0.0
	integerValue, _ := strconv.Atoi(value)

	for _, rec := range result.TimeEntries {
		if rec.ProjectID == integerValue {
			totalOvertimeHours += rec.HoursDecimal
		}
	}

	return totalOvertimeHours
}

// Context represents the context in which a strategy is executed.
type Context struct {
	strategy Strategy
}

func NewContext(strategy Strategy) *Context {
	return &Context{strategy: strategy}
}

func (c *Context) ExecuteStrategy(result *Response, value string) float64 {
	return c.strategy.CalculateHours(result, value)
}
