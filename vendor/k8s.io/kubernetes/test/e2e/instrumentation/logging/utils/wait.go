/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/kubernetes/test/e2e/framework"
)

// LogChecker is an interface for an entity that can check whether logging
// backend contains all wanted log entries.
type LogChecker interface {
	EntriesIngested() (bool, error)
	Timeout() error
}

// IngestionPred is a type of a function that checks whether all required
// log entries were ingested.
type IngestionPred func(string, []LogEntry) (bool, error)

// UntilFirstEntry is a IngestionPred that checks that at least one entry was
// ingested.
var UntilFirstEntry IngestionPred = func(_ string, entries []LogEntry) (bool, error) {
	return len(entries) > 0, nil
}

// UntilFirstEntryFromLog is a IngestionPred that checks that at least one
// entry from the log with a given name was ingested.
func UntilFirstEntryFromLog(log string) IngestionPred {
	return func(_ string, entries []LogEntry) (bool, error) {
		for _, e := range entries {
			if e.LogName == log {
				return true, nil
			}
		}
		return false, nil
	}
}

// TimeoutFun is a function that is called when the waiting times out.
type TimeoutFun func([]string, []bool) error

// JustTimeout returns the error with the list of names for which backend is
// still still missing logs.
var JustTimeout TimeoutFun = func(names []string, ingested []bool) error {
	failedNames := []string{}
	for i, name := range names {
		if !ingested[i] {
			failedNames = append(failedNames, name)
		}
	}
	return fmt.Errorf("timed out waiting for ingestion, still not ingested: %s",
		strings.Join(failedNames, ","))
}

var _ LogChecker = &logChecker{}

type logChecker struct {
	provider      LogProvider
	names         []string
	ingested      []bool
	ingestionPred IngestionPred
	timeoutFun    TimeoutFun
}

// NewLogChecker constructs a LogChecker for a list of names from custom
// IngestionPred and TimeoutFun.
func NewLogChecker(p LogProvider, pred IngestionPred, timeout TimeoutFun, names ...string) LogChecker {
	return &logChecker{
		provider:      p,
		names:         names,
		ingested:      make([]bool, len(names)),
		ingestionPred: pred,
		timeoutFun:    timeout,
	}
}

func (c *logChecker) EntriesIngested() (bool, error) {
	allIngested := true
	for i, name := range c.names {
		if c.ingested[i] {
			continue
		}
		entries := c.provider.ReadEntries(name)
		ingested, err := c.ingestionPred(name, entries)
		if err != nil {
			return false, err
		}
		if ingested {
			c.ingested[i] = true
		}
		allIngested = allIngested && ingested
	}
	return allIngested, nil
}

func (c *logChecker) Timeout() error {
	return c.timeoutFun(c.names, c.ingested)
}

// NumberedIngestionPred is a IngestionPred that takes into account sequential
// numbers of ingested entries.
type NumberedIngestionPred func(string, map[int]bool) (bool, error)

// NumberedTimeoutFun is a TimeoutFun that takes into account sequential
// numbers of ingested entries.
type NumberedTimeoutFun func([]string, map[string]map[int]bool) error

// NewNumberedLogChecker returns a log checker that works with numbered log
// entries generated by load logging pods.
func NewNumberedLogChecker(p LogProvider, pred NumberedIngestionPred,
	timeout NumberedTimeoutFun, names ...string) LogChecker {
	occs := map[string]map[int]bool{}
	return NewLogChecker(p, func(name string, entries []LogEntry) (bool, error) {
		occ, ok := occs[name]
		if !ok {
			occ = map[int]bool{}
			occs[name] = occ
		}
		for _, entry := range entries {
			if no, ok := entry.TryGetEntryNumber(); ok {
				occ[no] = true
			}
		}
		return pred(name, occ)
	}, func(names []string, _ []bool) error {
		return timeout(names, occs)
	}, names...)
}

// NewFullIngestionPodLogChecker returns a log checks that works with numbered
// log entries generated by load logging pods and waits until all entries are
// ingested. If timeout is reached, fraction is lost logs up to slack is
// considered tolerable.
func NewFullIngestionPodLogChecker(p LogProvider, slack float64, pods ...FiniteLoggingPod) LogChecker {
	podsMap := map[string]FiniteLoggingPod{}
	for _, p := range pods {
		podsMap[p.Name()] = p
	}
	return NewNumberedLogChecker(p, getFullIngestionPred(podsMap),
		getFullIngestionTimeout(podsMap, slack), getFiniteLoggingPodNames(pods)...)
}

func getFullIngestionPred(podsMap map[string]FiniteLoggingPod) NumberedIngestionPred {
	return func(name string, occ map[int]bool) (bool, error) {
		p := podsMap[name]
		ok := len(occ) == p.ExpectedLineCount()
		return ok, nil
	}
}

func getFullIngestionTimeout(podsMap map[string]FiniteLoggingPod, slack float64) NumberedTimeoutFun {
	return func(names []string, occs map[string]map[int]bool) error {
		totalGot, totalWant := 0, 0
		lossMsgs := []string{}
		for _, name := range names {
			got := len(occs[name])
			want := podsMap[name].ExpectedLineCount()
			if got != want {
				lossMsg := fmt.Sprintf("%s: %d lines", name, want-got)
				lossMsgs = append(lossMsgs, lossMsg)
			}
			totalGot += got
			totalWant += want
		}
		if len(lossMsgs) > 0 {
			framework.Logf("Still missing logs from:\n%s", strings.Join(lossMsgs, "\n"))
		}
		lostFrac := 1 - float64(totalGot)/float64(totalWant)
		if lostFrac > slack {
			return fmt.Errorf("still missing %.2f%% of logs, only %.2f%% is tolerable",
				lostFrac*100, slack*100)
		}
		framework.Logf("Missing %.2f%% of logs, which is lower than the threshold %.2f%%",
			lostFrac*100, slack*100)
		return nil
	}
}

// WaitForLogs checks that logs are ingested, as reported by the log checker
// until the timeout has passed. Function sleeps for interval between two
// log ingestion checks.
func WaitForLogs(c LogChecker, interval, timeout time.Duration) error {
	err := wait.Poll(interval, timeout, func() (bool, error) {
		return c.EntriesIngested()
	})
	if err == wait.ErrWaitTimeout {
		return c.Timeout()
	}
	return err
}

func getFiniteLoggingPodNames(pods []FiniteLoggingPod) []string {
	names := []string{}
	for _, p := range pods {
		names = append(names, p.Name())
	}
	return names
}
