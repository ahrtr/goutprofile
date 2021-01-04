goutprofile
======
goutprofile is the abbreviation for GO Unit Test Profile. It can be used to create a profile for each test file, and ensure that each test file is strictly consistent with the profile.

# A pain point of unit testing
A test file in golang is a file with the suffix "_test.go", and a test function is prefixed with "Test" or "Benchmark". When a test file has lots of test functions or it's a big file containing thousands of lines of code, then it's definitely difficult to understand all the test functions and what's the relationship between them. Accordingly, it is NOT easy to add new test functions and keep them consistent with the original design of test cases. 

# goutprofile (GO UT Profile)
goutprofile is exactly the solution to address the above pain point. 

Firstly, developers can create a profile for each test file. A profile is a YAML file with the suffix "_utprofile.yml" or "utprofile.yaml". It's a detailed design & document for the corresponding test file. It's super useful for anyone to understand the design of the test cases in the test file, such as,
- What test functions are included in the test file
- What's the type of each test function, i.e. test or benchmark
- What's the description for each test function
- How the test functions are classified

Secondly, goutprofile ensures that the naming and order of the test functions must conform exactly to the definition of profile. 

# Examples
The following example comes from **[demo3_utprofile.yml](testdata/subpkg1/demo3_utprofile.yml)**,
```
# go ut profile
source: demo3.go
test: demo3_test.go
description: this is the overview description for the profile
cases:
  - name: Demo3Case1
    type: benchmark
    description: this is top case1
  - name: Demo3Case2
    description: this is top case2
categories:
  - name: category1
    description: This is category 1
    cases:
    - name: Demo3Category1Case1
      type: benchmark 
      description: |
        This is the example description, and 
        it spans more than 
        one line
    - name: Demo3Category1Case2
      description: |
        This is the example description, and 
        it spans more than 
        one line       
```

A profile may contains the following fields,

| Field name | Description |
|------------|-------------|
| source | The source file which the test file tests, it's demo3.go in this example.|
| test   | The test file, it's demo3_test.go in this example |
| description | The description for this profile |
| cases | An array of test cases. Each case contains name, type and description. The type can be "test" or "benchmark". If type is not provided, then its value is "test" by default.|
| categories | An array of categories of test cases. Each category has a name and description, and may contains an array of cases and an arry of subcategories. | 

The test file (as below) conforms exactly to the above profile, 
```go
func BenchmarkDemo3Case1(b *testing.B) {}

func TestDemo3Case2(t *testing.T) {}

func BenchmarkDemo3Category1Case1(b *testing.B) {}

func TestDemo3Category1Case2(t *testing.T) {}

```

Refer to more examples under **[testdata](testdata)**

# How to use goutprofile in your golang project
## Step 1: create profiles
The first step is to create profiles for the test files. Usually you need to create the profiles firstly, and then implement the corresponding test functions in test files according to the profiles. But for some legacy projects, the test files may already exist, then you need to create profiles according to the implementations of test files.

## Step 2: validate test files against the profiles
You have two options. 

### Option 1: use command
Firstly, install the goutprofile command using the following command,
```
$ go install github.com/ahrtr/goutprofile/cmd/goutprofile
```

Note that if you are using an old golang version, i.e. go1.13.12, then the go.mod and go.sum may be modified unexpectedly, please rollback the changes on go.mod and go.sum. Refer to the following links to get more detailed info. Note that it's unrelated to goutprofile.  
- [golang/go#27643](https://github.com/golang/go/issues/27643)
- [golang/go#30515](https://github.com/golang/go/issues/30515)
- [golang/go#40276](https://github.com/golang/go/issues/40276)
- [https://go-review.googlesource.com/c/proposal/+/243077](https://go-review.googlesource.com/c/proposal/+/243077)
- [https://go-review.googlesource.com/c/go/+/254365](https://go-review.googlesource.com/c/go/+/254365)

Secondly, add the validation commands something like below into your Makefile or build script,
```
# This is your Makefile
all: validate xxxx

validate:
    goutprofile -d ./path1
    goutprofile -d ./path2 -r
```

Run "goutprofile -h" to get the usage of the command. 

### Option 2: use package
Firstly, write a validate_profiles.go something like below,
```go
// +build ignore

package main

import (
	"fmt"
	gup "github.com/ahrtr/goutprofile"
	"os"
)

func main() {
	if _, err := gup.ValidateDir(".", false); err != nil {
		fmt.Printf("Error validating goutprofile, error: %v\n", err)
		os.Exit(1)
	}
}
```

Secondly, add the validation commands something like below into your Makefile or build script,
```
# This is your Makefile
all: validate xxxx

validate:
    go run validate_profiles.go
```

# Exported functions
There are two exported functions in package github.com/ahrtr/goutprofile, 

```go
// ValidateDir validates all the ut profiles in the specified directory dir.
// If the second parameter recursive is true, then it recursively validates
// all its subdirectories as well.
//
// The first return parameter is a slice including all the profile names which
// passed the validation. It stops validating immediately once any profile's
// validation fails, and return an error in the second return parameter.
func ValidateDir(dir string, recursive bool) ([]string, error) 

// ValidateFile validates whether the test cases defined in the test file
// matches the UT profile fileName. An error will be returned immediately
// if it sees any mismatch.
//
// The test functions in the *_test.go file must be exactly the same as the case
// definition in the profile, including the order of appearance.
//
// Each test case in the profile is automatically prefixed a specific prefix
// according to its type. If its type is "test" or empty, then the prefix
// "Test". If its type is "benchmark", then the prefix is "Benchmark". All other
// types are invalid.
func ValidateFile(fileName string) error 
```

# Contribute to this repo
Anyone is welcome to contribute to this repo. Please raise an issue firstly, then fork this repo and submit a pull request.

Currently this repo is under heavily development, any helps are appreciated! 

# Support
If you need any support, please raise issues. 

If you have any suggestions or proposals, please also raise issues. Thanks!


