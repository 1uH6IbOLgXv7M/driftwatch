// Package schedule provides a simple interval-based scheduler for running
// periodic drift-detection jobs.
//
// A Scheduler is created with New, passing the desired interval, the Job
// function to execute, and an optional slog.Logger. Calling Run blocks until
// the supplied context is cancelled, invoking the job immediately and then
once per interval thereafter. Job errors are logged but do not stop the loop.
package schedule
