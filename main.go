package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/progressbar/v3"
	"github.com/schollz/seamlessloop/src/seamless"
)

var flagDebug bool
var flagInput, flagOutput string

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "debug mode")
	flag.StringVar(&flagInput, "in", "", "debug mode")
	flag.StringVar(&flagOutput, "out", "", "debug mode")
}

func main() {
	flag.Parse()
	if flagDebug {
		log.SetLevel("debug")
	} else {
		log.SetLevel("info")
	}

	if flagInput == "" {
		fmt.Println("need to specify input folder with --in")
		return
	}
	if flagOutput == "" {
		fmt.Println("need to specify output folder with --out")
		return
	}
	err := run()
	if err != nil {
		log.Error(err)
	}

	// err := runSpecial()
	// if err != nil {
	// 	log.Error(err)
	// }
}

func run() (err error) {
	var files = []string{}

	err = filepath.Walk(flagInput,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			path, _ = filepath.Abs(path)
			if !info.IsDir() && filepath.Ext(path) == ".wav" {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		return
	}

	log.Debug(files)
	for _, fname := range files {
		if err = loopit(fname); err != nil {
			log.Errorf("could not loop %s", fname)
			return
		}
	}
	return
}

func loopit(fname string) (err error) {
	fname2, bpm, err := seamless.Do(fname)
	if err != nil {
		return
	}

	_, filename2 := path.Split(filepath.ToSlash(fname))
	outFolder := path.Join(flagOutput, fmt.Sprint(bpm))
	outFolder = filepath.ToSlash(outFolder)
	outFile := path.Join(outFolder, filename2)
	err = os.MkdirAll(outFolder, os.ModePerm)
	if err != nil {
		return
	}
	_, err = copy(fname2, outFile)
	if err != nil {
		return
	}
	err = os.Remove(fname2)
	if err != nil {
		return
	}
	fmt.Printf("wrote '%s'\n", outFile)
	return
}

func runSpecial() (err error) {
	lineFile := "/media/zns/backup4tb/splice2/all_files.txt"
	numLines, err := NumLines(lineFile)
	bar := progressbar.Default(int64(numLines))
	numJobs := numLines
	if err != nil {
		return
	}

	f, err := os.Open(lineFile)
	if err != nil {
		return
	}
	defer f.Close()

	type job struct {
		fname string
	}
	type result struct {
		fname string
		err   error
	}

	jobs := make(chan job, numJobs)
	results := make(chan result, numJobs)
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(jobs <-chan job, results chan<- result) {
			for j := range jobs {
				// step 3: specify the work for the worker
				err = makeLoop(j.fname)
				results <- result{j.fname, err}
			}
		}(jobs, results)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fname := scanner.Text()
		jobs <- job{fname}

	}
	close(jobs)

	if err = scanner.Err(); err != nil {
		return
	}

	// step 5: do something with results
	for i := 0; i < numJobs; i++ {
		bar.Add(1)
		r := <-results
		if r.err != nil {
			// do something with error
			log.Errorf("'%s': %s", r.fname, r.err.Error())
		}
	}

	return
}

func makeLoop(filename string) (err error) {
	fname2, bpm, err := seamless.Do(filename)
	if err != nil {
		return
	}

	_, fname2name := path.Split(fname2)
	f1path, _ := path.Split(filename)
	f2path := strings.Replace(f1path, "splice2", "spliceloop", 1)
	f2path = path.Join(f2path, fmt.Sprint(bpm))
	finalName := path.Join(f2path, fname2name)
	err = os.MkdirAll(f2path, os.ModePerm)
	if err != nil {
		return
	}
	_, err = copy(fname2, finalName)
	if err != nil {
		return
	}
	err = os.Remove(fname2)
	if err != nil {
		return
	}
	return
}

func NumLines(fname string) (num int, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return
	}
	defer f.Close()
	num, err = lineCounter(f)
	return
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
