package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/shady831213/jarvism/core/errors"
	"io"
	"time"
)

// junitTestSuites is a collection of junit test suites.
type junitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []junitTestSuite
}

// junitTestSuite is a single junit test suite which may contain many
// testcases.
type junitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Name       string          `xml:"name,attr"`
	Time       string          `xml:"time,attr"`
	Properties []junitProperty `xml:"properties>property,omitempty"`
	TestCases  []junitTestCase
}

// junitTestCase is a single test case with its result.
type junitTestCase struct {
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	Time        string            `xml:"time,attr"`
	Status      string            `xml:"status,attr"`
	SkipMessage *junitSkipMessage `xml:"skipped,omitempty"`
	Failure     *junitFailure     `xml:"failure,omitempty"`
}

// junitSkipMessage contains the reason why a testcase was skipped.
type junitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// junitProperty represents a key/value pair used to define properties.
type junitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// junitFailure contains data related to a failed test.
type junitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

var suites junitTestSuites
var buildSuite *junitTestSuite
var testSuite *junitTestSuite

func initJunitXml(buildTotal, testTotal int) {
	suites = junitTestSuites{}
	buildSuite = &junitTestSuite{
		Tests:      buildTotal,
		Failures:   0,
		Time:       formatTime(0),
		Name:       "Builds",
		Properties: []junitProperty{},
		TestCases:  []junitTestCase{},
	}
	testSuite = &junitTestSuite{
		Tests:      testTotal,
		Failures:   0,
		Time:       formatTime(0),
		Name:       "Tests",
		Properties: []junitProperty{},
		TestCases:  []junitTestCase{},
	}
}

func updateBuild(result *errors.JVSRuntimeResult) {
	buildSuite.TestCases = append(buildSuite.TestCases, updateResult(result))
	if result.Status != errors.JVSRuntimePass {
		buildSuite.Failures++
	}
}

func updateTest(result *errors.JVSRuntimeResult) {
	(*testSuite).TestCases = append(testSuite.TestCases, updateResult(result))
	if result.Status != errors.JVSRuntimePass {
		testSuite.Failures++
	}
}

func updateResult(result *errors.JVSRuntimeResult) junitTestCase {
	test := junitTestCase{
		Classname: result.Name,
		Name:      result.Name,
		Time:      formatTime(0),
		Status:    errors.StatusString(result.Status),
		Failure:   nil,
	}

	if result.Status != errors.JVSRuntimePass {
		test.Failure = &junitFailure{
			Message:  "Failed",
			Type:     errors.StatusString(result.Status),
			Contents: result.GetMsg(),
		}
	}
	return test
}

func writeReport(w io.Writer) error {
	suites.Suites = append(suites.Suites, *buildSuite)
	suites.Suites = append(suites.Suites, *testSuite)
	// to xml
	bytes, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(w)

	if !noXMLHeader {
		writer.WriteString(xml.Header)
	}

	writer.Write(bytes)
	writer.WriteByte('\n')
	writer.Flush()

	return nil
}

func formatTime(d time.Duration) string {
	return fmt.Sprintf("%.3f", d.Seconds())
}
