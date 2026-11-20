package fixtures_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bidirekt/go-sdk/internal/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixtures = filepath.Join("..", "..", "testdata", "fixtures")

var sdkReplace = regexp.MustCompile(`(github\.com/bidirekt/go-sdk(?:[ \t]+\S+)?[ \t]*=>[ \t]*)(\.\.?/\S*)`)

func TestMain(m *testing.M) {
	if err := os.Setenv("GOWORK", "off"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestGeneratesTheExpectedResult(t *testing.T) {
	entries, err := os.ReadDir(fixtures)
	require.NoError(t, err)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		t.Run(entry.Name(), func(t *testing.T) {
			t.Parallel()
			generatesTheExpectedResult(t, filepath.Join(fixtures, entry.Name()))
		})
	}
}

func generatesTheExpectedResult(t *testing.T, caseDir string) {
	expectedYaml, hasYaml := readIfExists(t, filepath.Join(caseDir, "expected.yaml"))
	expectedStderr, hasStderr := readIfExists(t, filepath.Join(caseDir, "expected-stderr.txt"))
	require.True(t, hasYaml || hasStderr, "a case has expected.yaml, expected-stderr.txt or both")
	cwd := copyCase(t, caseDir)
	args := append([]string{"generate"}, caseArguments(t, cwd)...)
	output := projectOutput(cwd, args)
	before := filesUnder(t, cwd)

	got := command.ArgsToResult(args, cwd)

	assert.Equal(t, expectedStderr, got.Stderr, "stderr")
	after := filesUnder(t, cwd)
	if !hasYaml {
		assert.Equal(t, 1, got.ExitCode, "exit code")
		assert.Equal(t, before, after, "nothing is written")
		return
	}
	assert.Equal(t, 0, got.ExitCode, "exit code")
	relativeOutput, err := filepath.Rel(cwd, output)
	require.NoError(t, err)
	assert.Equal(t, expectedYaml, after[relativeOutput], "output")
	delete(before, relativeOutput)
	delete(after, relativeOutput)
	assert.Equal(t, before, after, "nothing but the output is written")
}

// The result is read from contracts/generated.yaml of the project folder: the case folder, or the folder
// that --project names in args.txt.
func projectOutput(cwd string, args []string) string {
	project := cwd
	for index, arg := range args {
		if arg == "--project" && index+1 < len(args) {
			project = filepath.Join(cwd, args[index+1])
		}
		if value, found := strings.CutPrefix(arg, "--project="); found {
			project = filepath.Join(cwd, value)
		}
	}
	return filepath.Join(project, "contracts", "generated.yaml")
}

func caseArguments(t *testing.T, cwd string) []string {
	args, _ := readIfExists(t, filepath.Join(cwd, "args.txt"))
	return strings.Fields(args)
}

// The copy points the replace of github.com/bidirekt/go-sdk at the checkout, as the case did from its
// place in the repository.
func copyCase(t *testing.T, caseDir string) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	cwd := filepath.Join(root, filepath.Base(caseDir))
	require.NoError(t, os.CopyFS(cwd, os.DirFS(caseDir)))
	for relative := range filesUnder(t, cwd) {
		if filepath.Base(relative) != "go.mod" && filepath.Base(relative) != "go.work" {
			continue
		}
		original, err := filepath.Abs(filepath.Join(caseDir, filepath.Dir(relative)))
		require.NoError(t, err)
		path := filepath.Join(cwd, relative)
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		rewritten := sdkReplace.ReplaceAllStringFunc(string(content), func(replace string) string {
			parts := sdkReplace.FindStringSubmatch(replace)
			return parts[1] + filepath.Join(original, parts[2])
		})
		require.NoError(t, os.WriteFile(path, []byte(rewritten), 0o644))
	}
	return cwd
}

func readIfExists(t *testing.T, path string) (string, bool) {
	t.Helper()
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", false
	}
	require.NoError(t, err)
	return string(content), true
}

func filesUnder(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(dir, path)
		files[relative] = string(content)
		return err
	})
	require.NoError(t, err)
	return files
}
