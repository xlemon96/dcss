package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"time"
)

var (
	_ TaskExecutor = Code{}
)

type Lang uint8

func (l Lang) String() string {
	switch l {
	case shell:
		return "shell"
	case python:
		return "python"
	case python3:
		return "python3"
	case golang:
		return "golang"
	case nodejs:
		return "nodejs"
	case windowsBat:
		return "windowsBat"
	default:
		return "unknown lang"
	}
}

const (
	shell Lang = iota + 1
	python3
	golang
	python
	nodejs
	windowsBat
)

const (
	modContent = `module crocodile

go `
	modName   = "go.mod"
	goNamePre = "crocodile_"
)

type Code struct {
	Lang     Lang   `json:"lang"`
	LangDesc string `json:"lang_desc"`
	Code     string `json:"code"`
}

func (c Code) Type() string {
	return c.Lang.String()
}

func (c Code) Run(ctx context.Context) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		var (
			exitCode = DefaultExitCode
			err      error
			codePath string
			cmd      *exec.Cmd
		)
		defer pw.Close()
		defer func() {
			now := time.Now().Local().Format("2006-01-02 15:04:05: ")
			pw.Write([]byte(fmt.Sprintf("%stask finished, return code:%5d", now, exitCode)))
			if codePath != "" {
				_ = os.Remove(codePath)
			}
		}()
		cmd, codePath, err = getCmd(ctx, c.Lang, c.Code)
		if err != nil {
			pw.Write([]byte(err.Error()))
			return
		}
		cmd.Stdout = pw
		cmd.Stderr = pw
		err = cmd.Start()
		if err != nil {
			pw.Write([]byte(err.Error()))
			return
		}
		err = cmd.Wait()
		if err != nil {
			pw.Write([]byte(err.Error()))
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			}
		} else {
			exitCode = 0
		}

	}()
	return pr
}

func getCmd(ctx context.Context, lang Lang, code string) (*exec.Cmd, string, error) {
	switch lang {
	case shell:
		return shellCmd(ctx, code)
	case python:
		return pythonCmd(ctx, code)
	case python3:
		return python3Cmd(ctx, code)
	case golang:
		return golangCmd(ctx, code)
	case nodejs:
		return nodeJsCmd(ctx, code)
	case windowsBat:
		return windowsBatCmd(ctx, code)
	default:
		return nil, "", fmt.Errorf("can not support lang: %d", lang)
	}
}

func shellCmd(ctx context.Context, code string) (*exec.Cmd, string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	tmpFile, err := ioutil.TempFile("", "*.sh")
	if err != nil {
		return nil, "", err
	}
	shellCodePath := tmpFile.Name()
	_, err = tmpFile.WriteString(code)
	if err != nil {
		return nil, "", err
	}
	tmpFile.Sync()
	tmpFile.Close()
	cmd := exec.CommandContext(ctx, shell, shellCodePath)
	return cmd, shellCodePath, nil
}

func pythonCmd(ctx context.Context, code string) (*exec.Cmd, string, error) {
	tmpFile, err := ioutil.TempFile("", "*.py")
	if err != nil {
		return nil, "", err
	}
	pythonCodePath := tmpFile.Name()
	_, err = tmpFile.WriteString(code)
	if err != nil {
		return nil, "", err
	}
	tmpFile.Sync()
	tmpFile.Close()
	cmd := exec.CommandContext(ctx, "python", pythonCodePath)
	return cmd, pythonCodePath, nil
}

func python3Cmd(ctx context.Context, code string) (*exec.Cmd, string, error) {
	tmpFile, err := ioutil.TempFile("", "*.py")
	if err != nil {
		return nil, "", err
	}
	python3CodePath := tmpFile.Name()
	_, err = tmpFile.WriteString(code)
	if err != nil {
		return nil, "", err
	}
	tmpFile.Sync()
	tmpFile.Close()
	cmd := exec.CommandContext(ctx, "python3", python3CodePath)
	return cmd, python3CodePath, nil
}

func golangCmd(ctx context.Context, code string) (*exec.Cmd, string, error) {
	cmd := exec.CommandContext(context.Background(), "go", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", err
	}
	pattern := `[0-1]\.[0-9]{1,2}`
	re := regexp.MustCompile(pattern)
	goVersion := re.FindString(string(out))
	if goVersion < "1.11" {
		err := errors.New("go version must rather equal go1.11 and enable go module")
		return nil, "", err
	}
	if os.Getenv("GO111MODULE") != "on" {
		os.Setenv("GO111MODULE", "on")
	}
	modContent := modContent + goVersion + "\n"
	tmpdir, err := ioutil.TempDir("", "crocodile_")
	if err != nil {
		return nil, "", err
	}
	err = ioutil.WriteFile(path.Join(tmpdir, modName), []byte(modContent), os.ModePerm)
	if err != nil {
		return nil, "", err
	}
	goNameFile := goNamePre + strconv.FormatInt(time.Now().Unix(), 10) + ".go"
	err = ioutil.WriteFile(path.Join(tmpdir, goNameFile), []byte(code), os.ModePerm)
	if err != nil {
		return nil, "", err
	}
	os.Chdir(tmpdir)
	goCmd := exec.CommandContext(ctx, "go", "run", goNameFile)
	return goCmd, tmpdir, nil
}

func nodeJsCmd(ctx context.Context, code string) (*exec.Cmd, string, error) {
	tmpFile, err := ioutil.TempFile("", "*.js")
	if err != nil {
		return nil, "", err
	}
	nodejsCodePath := tmpFile.Name()
	_, err = tmpFile.WriteString(code)
	if err != nil {
		return nil, "", err
	}
	tmpFile.Sync()
	tmpFile.Close()
	cmd := exec.CommandContext(ctx, "node", nodejsCodePath)
	return cmd, nodejsCodePath, nil
}

func windowsBatCmd(ctx context.Context, code string) (*exec.Cmd, string, error) {
	tmpFile, err := ioutil.TempFile("", "*.bat")
	if err != nil {
		return nil, "", err
	}
	batCodePath := tmpFile.Name()
	_, err = tmpFile.WriteString(code)
	if err != nil {
		return nil, "", err
	}
	tmpFile.Sync()
	tmpFile.Close()
	cmd := exec.CommandContext(ctx, "cmd", "/C", batCodePath)
	return cmd, batCodePath, nil
}
