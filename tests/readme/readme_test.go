package readme_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bidirekt/go-sdk/internal/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	firstExamples = "\n## Quickstart\n"
	afterExamples = "\n## The `generate` command\n"
)

var filePath = regexp.MustCompile("^`([^`]+)`$")

type example struct {
	heading string
	files   map[string]string
	yaml    string
	stderr  string
}

func readme(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "README.md"))
	require.NoError(t, err)
	return string(content)
}

// Every example between "## Quickstart" and "## The `generate` command" is a whole module: each file is a
// fenced block right below a line that holds only its path in backticks, and its result is the ```yaml
// (contracts/generated.yaml) or the ```text (stderr) block that follows its files.
func readmeToExamples(t *testing.T, text string) []example {
	t.Helper()
	start := strings.Index(text, firstExamples)
	end := strings.Index(text, afterExamples)
	require.True(t, start != -1 && end > start, "the README has its example sections")
	lines := strings.Split(text[start+1:end], "\n")
	var examples []example
	current := example{files: map[string]string{}}
	heading := ""
	path := ""
	for index := 0; index < len(lines); index++ {
		line := lines[index]
		if strings.HasPrefix(line, "## ") {
			heading = strings.TrimPrefix(line, "## ")
		}
		if match := filePath.FindStringSubmatch(line); match != nil {
			path = match[1]
			continue
		}
		language, isFence := strings.CutPrefix(line, "```")
		if !isFence {
			if line != "" {
				path = ""
			}
			continue
		}
		var body strings.Builder
		for index++; index < len(lines) && lines[index] != "```"; index++ {
			body.WriteString(lines[index] + "\n")
		}
		switch {
		case path != "":
			if current.yaml != "" || current.stderr != "" {
				examples = append(examples, current)
				current = example{files: map[string]string{}}
			}
			current.heading = heading
			current.files[path] = body.String()
		case language == "yaml":
			current.yaml = body.String()
		case language == "text":
			current.stderr = body.String()
		}
		path = ""
	}
	if len(current.files) > 0 {
		examples = append(examples, current)
	}
	return examples
}

func newModule(t *testing.T, files map[string]string) string {
	t.Helper()
	checkout, err := filepath.Abs(filepath.Join("..", ".."))
	require.NoError(t, err)
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	files["go.mod"] = "module example.com/petstore\n\ngo 1.26.0\n\nrequire github.com/bidirekt/go-sdk v0.0.0\n\nreplace github.com/bidirekt/go-sdk => " + checkout + "\n"
	for name, content := range files {
		require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
	}
	return dir
}

func TestGeneratesWhatTheReadmeShows(t *testing.T) {
	examples := readmeToExamples(t, readme(t))
	require.NotEmpty(t, examples, "examples")
	for _, each := range examples {
		t.Run(each.heading, func(t *testing.T) {
			t.Parallel()
			require.True(t, each.yaml != "" || each.stderr != "", "the example shows a result")
			dir := newModule(t, each.files)

			result := command.ArgsToResult([]string{"generate"}, dir)

			assert.Equal(t, each.stderr, result.Stderr, "stderr")
			generated, err := os.ReadFile(filepath.Join(dir, "contracts", "generated.yaml"))
			if each.yaml == "" {
				assert.Equal(t, 1, result.ExitCode, "exit code")
				assert.True(t, os.IsNotExist(err), "nothing is written")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, 0, result.ExitCode, "exit code")
			assert.Equal(t, each.yaml, string(generated), "output")
		})
	}
}

func TestShowsTheCheckMessagesGeneratePrints(t *testing.T) {
	text := readme(t)
	dir := newModule(t, map[string]string{"pets.go": `package petstore

import "github.com/bidirekt/go-sdk"

var _ = bidirekt.Provides("GET", "/health", bidirekt.Response(204, nil))
`})

	missing := command.ArgsToResult([]string{"generate", "--check"}, dir)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "contracts"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "contracts", "generated.yaml"), []byte("# stale\n"), 0o644))
	stale := command.ArgsToResult([]string{"generate", "--check"}, dir)

	for _, result := range []command.Result{missing, stale} {
		assert.Equal(t, 1, result.ExitCode, "exit code")
		message := "`" + strings.TrimSuffix(result.Stderr, "\n") + "`"
		assert.True(t, strings.Contains(text, message), "the README shows %s", message)
	}
}
