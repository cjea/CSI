package main

import (
	"bufio"
	"cj-shell/pkg/dumbstack"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var cdStack = dumbstack.New()

const Prompt = "#> "

var Builtins = []string{"exit", "cd", "pwd"}

var (
	ErrCmdNotFound     = errors.New("command not found")
	ErrBuiltinNotFound = errors.New("builtin not found")
)

type Cmd struct {
	Cmd        string
	Args       []string
	Background bool
}

func (c Cmd) IsBuiltin() bool {
	for i := range Builtins {
		if Builtins[i] == c.Cmd {
			return true
		}
	}
	return false
}

func WritePrompt() {
	fmt.Printf("%s", Prompt)
}

type REPL struct {
	scanner *bufio.Scanner
	done    bool
}

func NewREPL() *REPL {
	scanner := bufio.NewScanner(os.Stdin)
	return &REPL{scanner: scanner}
}

// Read returns the next line of text from stdin.
func (r *REPL) Read() (string, error) {
	if ok := r.scanner.Scan(); !ok {
		err := r.scanner.Err()
		r.done = (err == nil)
		return "", err
	}

	return r.scanner.Text(), nil
}

func pathDirs() []string {
	return strings.Split(os.Getenv("PATH"), ":")
}

// AllPathPrefixes returns a slice of commands, prefixed by different PATH dirs.
// To execute the command, run down the list and keep the first one that works.
//
//      AllPathPrefixes("foo") // => []{"foo", "/bin/foo", "/usr/bin/foo", ...}
func AllPathPrefixes(c Cmd) []*exec.Cmd {
	dirs := pathDirs()
	cmd, args := c.Cmd, c.Args
	ret := make([]*exec.Cmd, len(dirs)+1)
	ret[0] = exec.Command(cmd, args...)

	for idx := range dirs {
		prefixedPath := path.Join(dirs[idx], cmd)
		ret[idx+1] = exec.Command(prefixedPath, args...)
	}
	return ret
}

func ExecuteCd(dest string) {
	if dest == "-" {
		dir, ok := cdStack.Pop()
		if ok {
			dest = dir
		} else {
			dest, _ = os.Getwd()
		}
	}
	fmt.Println("Changing directory to " + dest)
	os.Chdir(dest)
}

func ExecutePwd() {
	wd, _ := os.Getwd()
	fmt.Printf("%s\n", wd)
}

func (r *REPL) ExecuteBuiltin(cmd Cmd) ([]byte, error) {
	c, args := cmd.Cmd, cmd.Args
	switch c {
	case "exit":
		r.done = true
		return nil, nil
	case "cd":
		dest := "~"
		if len(args) > 0 {
			dest = args[0]
		}
		ExecuteCd(dest)
		return nil, nil
	case "pwd":
		ExecutePwd()
		return nil, nil
	}

	return nil, fmt.Errorf("%s is not a builtin: %w", cmd.Cmd, ErrBuiltinNotFound)
}

// Execute forks a process and runs the given string as a command.
func (r *REPL) Execute(cmd Cmd) ([]byte, error) {
	var out []byte
	var err error

	if cmd.IsBuiltin() {
		return r.ExecuteBuiltin(cmd)
	}

	commands := AllPathPrefixes(cmd)
	for _, c := range commands {
		out, err = c.Output()
		if err == nil {
			return out, nil
		}
	}

	return nil, fmt.Errorf("%v: %w", err, ErrCmdNotFound)
}

// ParseLine evaluates a command
func (r *REPL) ParseLine(str string) Cmd {
	parts := strings.Split(str, " ")
	f := parts[0]
	args := parts[1:]
	return Cmd{
		Cmd:  f,
		Args: args,
	}
}

// Write returns the next line of text from stdin.
func (r *REPL) Write(s string) {
	fmt.Println(s)
}

func (r *REPL) Done() bool {
	return r.done
}

func (r *REPL) Loop() {
	for !r.Done() {
		WritePrompt()
		str, err := r.Read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		cmd := r.ParseLine(str)
		out, err := r.Execute(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		r.Write(string(out))
	}
}

func main() {
	fmt.Println("Welcome to CJ shell")
	NewREPL().Loop()
	fmt.Println("** Goodbye **")
	os.Exit(0)
}
