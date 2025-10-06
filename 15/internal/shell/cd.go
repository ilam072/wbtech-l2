package shell

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func HandCd(args []string) error {
	if len(args) > 1 {
		return errors.New("too many arguments")
	}

	path, err := targetPath(args)
	if err != nil {
		return err
	}

	exists, isDir := isDirExists(path)
	if !exists {
		return fmt.Errorf("%s: no such file or directory", path)
	}
	if !isDir {
		return fmt.Errorf("%s: not a directory", path)
	}

	return changeDirectory(path)
}

// определить, куда переходить (HOME, OLDPWD, путь с тильдой/переменными)
func targetPath(args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			return "", errors.New("HOME not set")
		}
		return home, nil
	}

	if args[0] == "-" {
		old := os.Getenv("OLDPWD")
		if old == "" {
			return "", errors.New("OLDPWD not set")
		}
		return old, nil
	}

	return expandPath(args[0]), nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil && home != "" {
			if path == "~" {
				path = home
			} else if strings.HasPrefix(path, "~/") {
				path = filepath.Join(home, path[2:])
			}
		}
	}
	return path
}

func changeDirectory(target string) error {
	oldPwd, _ := os.Getwd()

	if err := os.Chdir(target); err != nil {
		return fmt.Errorf("%s: no such file or directory", target)
	}

	newPwd, _ := os.Getwd()
	_ = os.Setenv("OLDPWD", oldPwd)
	_ = os.Setenv("PWD", newPwd)

	fmt.Println(newPwd)

	return nil
}

func isDirExists(path string) (exists bool, isDir bool) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false
		}
		return true, false
	}
	return true, info.IsDir()
}
