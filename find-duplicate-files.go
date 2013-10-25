package main

import "errors"
import "path/filepath"
import "flag"
import "fmt"
import "os"


func findFiles(directories []string) ([]string, error) {
    files := make([]string, 0, 100)
    // Can't use ranged 'for' as the length of directories changes during iteration.
    for x := 0; x < len(directories); x++ {
        directory := directories[x]
        dirFile, err := os.Open(directory)
        if err != nil {
            return nil, err
        }
        contents, err := dirFile.Readdir(0)
        if err != nil {
            return nil, err
        }
        for _, content := range contents {
            path := filepath.Join(directory, content.Name())
            if content.IsDir() {
                directories = append(directories, path)
            } else {
                files = append(files, path)
            }
        }
    }

    return files, nil
}


// Validate the passed-in arguments are directories.
func validateArgs(args []string) error {
    if len(args) < 1 {
        return errors.New("expected 1 or more arguments")
    }
    for _, directory := range args {
        dir, err := os.Open(directory)
        if err != nil {
            return err
        }
        defer dir.Close()
        info, err := dir.Stat()
        if err != nil {
            return err
        } else if !info.IsDir() {
            return errors.New(fmt.Sprintf("%v is not a directory", directory))
        }
    }

    return nil
}


func errorExit(err error) {
    // Would be more correct if flags.out() were publicly available instead of hard coding os.Stderr.
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}

// find-duplicate-files takes 1 or more directories on the command-line,
// recurses into all of them, and prints out what files are duplicates of
// each other.
func main() {
    flag.Parse()
    directories := flag.Args()
    err := validateArgs(directories); if err != nil {
        errorExit(err)
    }
    _, err = findFiles(directories)
    if err != nil {
        errorExit(err)
    }
    // XXX hash files
    // XXX 4096 seems to be a common buffer size for file systems and HDDs
    // XXX collect duplicates
    // XXX print duplicates, sorted
}
