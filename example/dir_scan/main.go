package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/curlymon/pipes"
	"github.com/curlymon/pipes/async"
)

const (
	Workers  = 10
	ChanSize = 10
)

func main() {
	dir, err := getDir(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
	filePipe := pipeline(true, dir)
	filePipe = async.MapWithErrorSink(Workers, ChanSize, openFile, logError("error opening file"), filePipe)
	filePipe = async.MapWithErrorSink(Workers, ChanSize, multiHash, logError("error multi hashing file"), filePipe)
	filePipe = async.MapWithErrorSink(Workers, ChanSize, closeFile, logError("error closing file"), filePipe)
	// filePipe = pipes.Tap(ChanSize, logFileFound, filePipe)
	resultPipe := pipes.Window(ChanSize, time.Second, compileResult, newResults, filePipe)
	resultPipe = pipes.Tap(ChanSize, logAny[*Results], resultPipe)
	log.Println(pipes.Reduce(compileResults, &Results{}, resultPipe))
}

func logAny[T any](t T) {
	log.Println(t)
}

func newResults() *Results {
	return new(Results)
}

func pipeline(recurse bool, dir string) pipes.ChanPull[FileInfo] {
	out := pipes.New[FileInfo](10)

	go func() {
		defer out.Close()
		if err := filepath.WalkDir(dir, walkFunc(dir, recurse, out)); err != nil {
			log.Printf("error walking directory: dir=%s, err=%s", dir, err)
		}
	}()

	return out.ChanPull()
}

func walkFunc(dir string, recurse bool, out chan<- FileInfo) func(path string, d fs.DirEntry, err error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("path=%s, err=%s\n", path, err)
		}
		switch d.IsDir() {
		case true:
			if !recurse && dir != path {
				return fs.SkipDir
			}
		case false:
			out <- FileInfo{
				Path:  path,
				Entry: d,
			}
		}
		return nil
	}
}

func getDir(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("must pass directory path after binary name")
	}
	dir := args[1]
	if !fs.ValidPath(dir) {
		return "", errors.New("must pass valid path: " + dir)
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("must pass valid path: %s", err)
	}
	dir = strings.Trim(dir, "\"")
	return dir, nil
}

func compileResult(fi FileInfo, results *Results) *Results {
	results.Found++
	results.TotalDuration += time.Since(fi.Start)
	return results
}

func compileResults(result, results *Results) *Results {
	results.Found += result.Found
	results.TotalDuration += result.TotalDuration
	return results
}

type FileInfo struct {
	Path   string
	Entry  fs.DirEntry
	File   *os.File
	MD5    []byte
	SHA1   []byte
	SHA256 []byte
	SHA512 []byte
	Start  time.Time
}

type Results struct {
	Found         int
	TotalDuration time.Duration
}

func (r Results) String() string {
	avg := time.Duration(float64(r.TotalDuration) / float64(r.Found))
	return fmt.Sprintf("Found: %d, Avg: %s, Total: %s", r.Found, avg, r.TotalDuration)
}

func openFile(fi FileInfo) (FileInfo, error) {
	f, err := os.Open(fi.Path)
	if err != nil {
		return FileInfo{}, err
	}
	fi.File = f
	fi.Start = time.Now()
	return fi, nil
}

func closeFile(fi FileInfo) (FileInfo, error) {
	return fi, fi.File.Close()
}

var buffers = sync.Pool{
	New: func() any {
		buf := make([]byte, 32*1024)
		return &buf
	},
}

func multiHash(fi FileInfo) (FileInfo, error) {
	defer fi.File.Seek(0, 0)
	buf := buffers.Get().(*[]byte)
	defer buffers.Put(buf)
	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()
	sha512 := sha512.New()
	w := io.MultiWriter(md5, sha1, sha256, sha512)
	if _, err := io.CopyBuffer(w, fi.File, *buf); err != nil {
		return FileInfo{}, err
	}
	fi.MD5 = md5.Sum(nil)
	fi.SHA1 = sha1.Sum(nil)
	fi.SHA256 = sha256.Sum(nil)
	fi.SHA512 = sha512.Sum(nil)
	return fi, nil
}

func logFileFound(fi FileInfo) {
	info, _ := fi.Entry.Info()
	log.Printf(
		"Found file: name=%s, md5=%X, sha1=%X, sha256=%x, sha512=%x, size=%d, took=%s",
		fi.Entry.Name(), fi.MD5[:8], fi.SHA1[:8], fi.SHA256[:8], fi.SHA512[:8], info.Size(), time.Since(fi.Start),
	)
}

func logError(message string) func(error) {
	return func(err error) {
		log.Println(message, err)
	}
}
