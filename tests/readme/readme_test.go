package readme_test

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/bidirekt/go-sdk/internal/command"
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
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	return string(content)
}

// Every example between "## Quickstart" and "## The `generate` command" is a whole module: each file is a
// fenced block right below a line that holds only its path in backticks, and its result is the ```yaml
// (contracts/generated.yaml) or the ```text (stderr) block that follows its files.
func readmeToExamples(t *testing.T, text string) []example {
	t.Helper()
	start := strings.Index(text, firstExamples)
	end := strings.Index(text, afterExamples)
	if start == -1 || end <= start {
		t.Fatalf("the README has its example sections: got false, want true")
	}
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
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	files["go.mod"] = "module example.com/petstore\n\ngo 1.26.0\n\nrequire github.com/bidirekt/go-sdk v0.0.0\n\nreplace github.com/bidirekt/go-sdk => " + checkout + "\n"
	for name, content := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755); err != nil {
			t.Fatalf("got error %v, want none", err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("got error %v, want none", err)
		}
	}
	return dir
}

func TestGeneratesWhatTheReadmeShows(t *testing.T) {
	examples := readmeToExamples(t, readme(t))
	if len(examples) == 0 {
		t.Fatalf("examples: got none, want some")
	}
	for _, each := range examples {
		t.Run(each.heading, func(t *testing.T) {
			t.Parallel()
			if each.yaml == "" && each.stderr == "" {
				t.Fatalf("the example shows a result: got false, want true")
			}
			dir := newModule(t, each.files)

			result := command.ArgsToResult([]string{"generate"}, dir)

			equal(t, "stderr", result.Stderr, each.stderr)
			generated, err := os.ReadFile(filepath.Join(dir, "contracts", "generated.yaml"))
			if each.yaml == "" {
				equal(t, "exit code", result.ExitCode, 1)
				if !os.IsNotExist(err) {
					t.Errorf("nothing is written: got false, want true")
				}
				return
			}
			if err != nil {
				t.Fatalf("got error %v, want none", err)
			}
			equal(t, "exit code", result.ExitCode, 0)
			equal(t, "output", string(generated), each.yaml)
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
	if err := os.MkdirAll(filepath.Join(dir, "contracts"), 0o755); err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "contracts", "generated.yaml"), []byte("# stale\n"), 0o644); err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	stale := command.ArgsToResult([]string{"generate", "--check"}, dir)

	for _, result := range []command.Result{missing, stale} {
		equal(t, "exit code", result.ExitCode, 1)
		message := "`" + strings.TrimSuffix(result.Stderr, "\n") + "`"
		if !strings.Contains(text, message) {
			t.Errorf("the README shows %s: got false, want true", message)
		}
	}
}

func equal(t *testing.T, label string, got, want any) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		return
	}
	if gotText, isText := got.(string); isText {
		t.Errorf("%s:\ngot:\n%s\nwant:\n%s", label, gotText, want)
		return
	}
	t.Errorf("%s:\ngot:  %#v\nwant: %#v", label, got, want)
}
