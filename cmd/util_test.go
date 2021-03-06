package cmd

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_textToMarkdown(t *testing.T) {
	basestring := "This string should have two spaces at the end."
	teststring := basestring + "\n"
	newteststring := textToMarkdown(teststring)
	assert.Equal(t, basestring+"  \n", newteststring)
}

func Test_getCurrentBranchMR(t *testing.T) {
	repo := copyTestRepo(t)

	// make sure the branch does not exist
	cmd := exec.Command("git", "branch", "-D", "mrtest")
	cmd.Dir = repo
	cmd.CombinedOutput()

	cmd = exec.Command(labBinaryPath, "mr", "checkout", "1")
	cmd.Dir = repo
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}

	curDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(repo)
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	mrNum := getCurrentBranchMR("zaquestion/test")
	err = os.Chdir(curDir)
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}

	assert.Equal(t, 1, mrNum)
}

func Test_parseArgsStringAndID(t *testing.T) {
	tests := []struct {
		Name           string
		Args           []string
		ExpectedString string
		ExpectedInt    int64
		ExpectedErr    string
	}{
		{
			Name:           "No Args",
			Args:           nil,
			ExpectedString: "",
			ExpectedInt:    0,
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg remote",
			Args:           []string{"origin"},
			ExpectedString: "origin",
			ExpectedInt:    0,
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg non remote",
			Args:           []string{"foo"},
			ExpectedString: "foo",
			ExpectedInt:    0,
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg page",
			Args:           []string{"100"},
			ExpectedString: "",
			ExpectedInt:    100,
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg invalid page",
			Args:           []string{"asdf100"},
			ExpectedString: "asdf100",
			ExpectedInt:    0,
			ExpectedErr:    "",
		},
		{
			Name:           "2 arg str page",
			Args:           []string{"origin", "100"},
			ExpectedString: "origin",
			ExpectedInt:    100,
			ExpectedErr:    "",
		},
		{
			Name:           "2 arg valid str valid page",
			Args:           []string{"foo", "100"},
			ExpectedString: "foo",
			ExpectedInt:    100,
			ExpectedErr:    "",
		},
		{
			Name:           "2 arg valid str invalid page",
			Args:           []string{"foo", "asdf100"},
			ExpectedString: "foo",
			ExpectedInt:    0,
			ExpectedErr:    "strconv.ParseInt: parsing \"asdf100\": invalid syntax",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test := test
			t.Parallel()
			s, i, err := parseArgsStringAndID(test.Args)
			if err != nil {
				assert.EqualError(t, err, test.ExpectedErr)
			}
			assert.Equal(t, test.ExpectedString, s)
			assert.Equal(t, test.ExpectedInt, i)
		})
	}
}

func Test_parseArgsRemoteAndID(t *testing.T) {
	tests := []struct {
		Name           string
		Args           []string
		ExpectedString string
		ExpectedInt    int64
		ExpectedErr    string
	}{
		{
			Name:           "No Args",
			Args:           nil,
			ExpectedString: "zaquestion/test",
			ExpectedInt:    0,
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg remote",
			Args:           []string{"lab-testing"},
			ExpectedString: "lab-testing/test",
			ExpectedInt:    0,
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg non remote",
			Args:           []string{"foo"},
			ExpectedString: "",
			ExpectedInt:    0,
			ExpectedErr:    "foo is not a valid remote or number",
		},
		{
			Name:           "1 arg page",
			Args:           []string{"100"},
			ExpectedString: "zaquestion/test",
			ExpectedInt:    100,
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg invalid page",
			Args:           []string{"asdf100"},
			ExpectedString: "",
			ExpectedInt:    0,
			ExpectedErr:    "asdf100 is not a valid remote or number",
		},
		{
			Name:           "2 arg remote page",
			Args:           []string{"origin", "100"},
			ExpectedString: "zaquestion/test",
			ExpectedInt:    100,
			ExpectedErr:    "",
		},
		{
			Name:           "2 arg invalid remote valid page",
			Args:           []string{"foo", "100"},
			ExpectedString: "",
			ExpectedInt:    0,
			ExpectedErr:    "foo is not a valid remote",
		},
		{
			Name:           "2 arg invalid remote invalid page",
			Args:           []string{"foo", "asdf100"},
			ExpectedString: "",
			ExpectedInt:    0,
			ExpectedErr:    "strconv.ParseInt: parsing \"asdf100\": invalid syntax",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test := test
			t.Parallel()
			s, i, err := parseArgsRemoteAndID(test.Args)
			if err != nil {
				assert.EqualError(t, err, test.ExpectedErr)
			}
			assert.Equal(t, test.ExpectedString, s)
			assert.Equal(t, test.ExpectedInt, i)
		})
	}
}

func Test_parseArgsRemoteAndProject(t *testing.T) {
	tests := []struct {
		Name           string
		Args           []string
		ExpectedRemote string
		ExpectedString string
		ExpectedErr    string
	}{
		{
			Name:           "No Args",
			Args:           nil,
			ExpectedRemote: "zaquestion/test",
			ExpectedString: "",
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg remote",
			Args:           []string{"lab-testing"},
			ExpectedRemote: "lab-testing/test",
			ExpectedString: "",
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg non remote",
			Args:           []string{"foo123"},
			ExpectedRemote: "zaquestion/test",
			ExpectedString: "foo123",
			ExpectedErr:    "",
		},
		{
			Name:           "1 arg page",
			Args:           []string{"100"},
			ExpectedRemote: "zaquestion/test",
			ExpectedString: "100",
			ExpectedErr:    "",
		},
		{
			Name:           "2 arg remote and string",
			Args:           []string{"origin", "foo123"},
			ExpectedRemote: "zaquestion/test",
			ExpectedString: "foo123",
			ExpectedErr:    "",
		},
		{
			Name:           "2 arg invalid remote and string",
			Args:           []string{"foo", "string123"},
			ExpectedRemote: "",
			ExpectedString: "",
			ExpectedErr:    "foo is not a valid remote",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test := test
			t.Parallel()
			r, s, err := parseArgsRemoteAndProject(test.Args)
			if err != nil {
				assert.EqualError(t, err, test.ExpectedErr)
			}
			assert.Equal(t, test.ExpectedRemote, r)
			assert.Equal(t, test.ExpectedString, s)
		})
	}
}
