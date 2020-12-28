package goutprofile

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ahrtr/gocontainer/stack"
	"gopkg.in/yaml.v2"
)

const (
	typeTest        = "test"
	prefixTest      = "Test"
	typeBenchmark   = "benchmark"
	prefixBenchmark = "Benchmark"
)

// utProfile is the unit test profile definition
type utProfile struct {
	Source      string         `yaml:"source"`
	Test        string         `yaml:"test"`
	Description string         `yaml:"description"`
	Cases       []testCase     `yaml:"cases,omitempty"`
	Categories  []testCategory `yaml:"categories,omitempty"`
}

// iterator is used to iterate all the test cases in UTProfile
type iterator struct {
	Index  int // the index of the next test case to read, starting from 0
	Cases  []testCase
	Stacks stack.Interface
}

// TestCase is the test case definition
type testCase struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

// testCategory is the test case category definition
// It may have cases and sub categories
type testCategory struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Cases       []testCase     `yaml:"cases,omitempty"`
	Categories  []testCategory `yaml:"categories,omitempty"`
}

// testFunc is the generated test function
type testFunc struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

// loadUTProfile loads the UT profile from a file
func loadUTProfile(fileName string) (*utProfile, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	profile := utProfile{}
	if err = yaml.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ut profile, filename: %s, error: %v", fileName, err)
	}

	return &profile, nil
}

// tterator returns the Iterator for the UTProfile
func (utp *utProfile) iterator() *iterator {
	it := &iterator{
		Index:  0,
		Cases:  utp.Cases,
		Stacks: stack.New(),
	}

	for i := len(utp.Categories) - 1; i >= 0; i-- {
		it.Stacks.Push(utp.Categories[i])
	}

	return it
}

// next returns the next test case
func (it *iterator) next() (testCase, error) {
	if len(it.Cases) == 0 || it.Index >= len(it.Cases) {
		for !it.Stacks.IsEmpty() {
			category := it.Stacks.Pop().(testCategory)

			for i := len(category.Categories) - 1; i >= 0; i-- {
				it.Stacks.Push(category.Categories[i])
			}

			if len(category.Cases) > 0 {
				it.Cases = category.Cases
				it.Index = 0
				break
			}
		}
	}

	if len(it.Cases) > 0 && it.Index < len(it.Cases) {
		tc := it.Cases[it.Index]
		it.Index++
		return tc, nil
	}

	return testCase{}, errors.New("test cases not found")
}
