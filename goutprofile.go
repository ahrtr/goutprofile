package goutprofile

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	lm "github.com/ahrtr/gocontainer/map/linkedmap"
)

// ValidateDir validates all the ut profiles in the specified directory dir.
// If the second parameter recursive is true, then it recursively validates
// all its subdirectories as well.
//
// The first return parameter is a slice including all the profile names which
// passes the validation. It stops validating immediately once any profile's
// validation fails, and return an error in the second return parameter.
func ValidateDir(dir string, recursive bool) ([]string, error) {
	validatedProfile := []string{}
	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to access path %q: %v", path, err)
			}

			if !isValidProfile(info) {
				return nil
			}

			if err := ValidateFile(path); err != nil {
				return err
			}

			validatedProfile = append(validatedProfile, path)
			return nil
		})

		if err != nil {
			return validatedProfile, fmt.Errorf("error walking through the path %q: %v", dir, err)
		}
		return validatedProfile, nil
	}

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return validatedProfile, fmt.Errorf("failed to read dir %q: %v", dir, err)
	}

	for _, info := range fileInfos {
		if !isValidProfile(info) {
			continue
		}
		fileName := filepath.Join(dir, info.Name())
		if err = ValidateFile(fileName); err != nil {
			return validatedProfile, err
		}
		validatedProfile = append(validatedProfile, fileName)
	}

	return validatedProfile, nil
}

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
func ValidateFile(fileName string) error {
	fmt.Printf("Validating %s\n", fileName)
	// load UT profile
	utp, err := loadUTProfile(fileName)
	if err != nil {
		return processError(fmt.Errorf("failed to load profile, error: %v", err))
	}

	// check whether the source file exists or not
	srcFile := getAbsPath(utp.Source, filepath.Dir(fileName))
	if _, err = os.Stat(srcFile); err != nil {
		if os.IsNotExist(err) {
			return processError(fmt.Errorf("the source file %q defined in %q doesn't exit", srcFile, fileName))
		}
	}

	// parse the test cases in the test file
	testFile := getAbsPath(utp.Test, filepath.Dir(fileName))
	casesInFile, err := parseTestFile(testFile)
	if err != nil {
		return processError(err)
	}

	// generate the test cases base on the UT profile definitions
	casesInProfile, err := generateTestFunctions(utp)
	if err != nil {
		return processError(err)
	}

	// diff casesInFile and casesInProfile
	if err = checkTestCases(casesInFile, casesInProfile); err != nil {
		return processError(err)
	}

	fmt.Printf("Validating %s passed\n", fileName)
	return nil
}

func processError(err error) error {
	fmt.Println(err)
	return err
}

// getAbsPath returns the absolute path of the file
func getAbsPath(fileName, path string) string {
	if filepath.IsAbs(fileName) {
		return fileName
	}

	return filepath.Join(path, fileName)
}

// valid utprofile file must have suffix "_utprofile.yml" or "_utprofile.yaml"
func isValidProfile(info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	return strings.HasSuffix(info.Name(), "_utprofile.yml") || strings.HasSuffix(info.Name(), "_utprofile.yaml")
}

// ParseTestFile parses the given test file and returns all the functions starting with "Test" or "Benchmark"
func parseTestFile(testFile string) (lm.Interface, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, testFile, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test file: %s, error: %v", testFile, err)
	}

	v := visitor{
		functions: lm.New(),
	}
	ast.Walk(&v, f)

	return v.functions, nil
}

type visitor struct {
	functions lm.Interface
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.FuncDecl:
		if strings.HasPrefix(d.Name.Name, "Test") || strings.HasPrefix(d.Name.Name, "Benchmark") {
			v.functions.Put(d.Name.Name, struct{}{})
		}
	}
	return v
}

// GenerateTestFunctions generates all the test functions based on the given UTProfile
func generateTestFunctions(utp *utProfile) ([]*testFunc, error) {
	var generatedFuncs []*testFunc

	it := utp.iterator()
	for {
		tc, err := it.next()
		if err != nil {
			break
		}
		typeStr := strings.ToLower(tc.Type)
		switch typeStr {
		case typeTest, "":
			generatedFuncs = append(generatedFuncs, &testFunc{
				Name:        prefixTest + tc.Name,
				Type:        typeTest,
				Description: tc.Description,
			})
		case typeBenchmark:
			generatedFuncs = append(generatedFuncs, &testFunc{
				Name:        prefixBenchmark + tc.Name,
				Type:        typeBenchmark,
				Description: tc.Description,
			})
		default:
			return generatedFuncs, fmt.Errorf("invalid test case type: %s", tc.Type)
		}
	}

	return generatedFuncs, nil
}

// CheckTestCases checks whether the test cases in the test file match the test cases generated from UT profile.
// If they do not match, then an error will be returned.
func checkTestCases(casesInFile lm.Interface, casesInProfile []*testFunc) error {
	// Check whether the test cases number are equal
	if casesInFile.Size() != len(casesInProfile) {
		return fmt.Errorf("there are %d cases in test file, but there are %d cases in ut profile", casesInFile.Size(), len(casesInProfile))
	}

	// go through all the cases in both test file and profile
	it, hasNext := casesInFile.Iterator()
	var key interface{}
	index := 0
	for hasNext {
		key, _, hasNext = it()
		funcInFile := key.(string)
		funcInProfile := casesInProfile[index]
		index++
		//log.Printf("funcInFile: %s, funcInProfile: %s", funcInFile, funcInProfile.Name)

		if !casesInFile.ContainsKey(funcInProfile.Name) {
			return fmt.Errorf("the case %q in ut profile doesn't exist in the test file", funcInProfile.Name)
		}

		if funcInFile != funcInProfile.Name {
			return fmt.Errorf("the case %q in test file doesn't match the test %q in profile at position %d", funcInFile, funcInProfile.Name, index)
		}
	}

	return nil
}
