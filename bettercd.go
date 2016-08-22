package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// MaxHistOpt MAX hisotry size option
	MaxHistOpt = flag.Int("i", 200, "Max number of history")
	// LogOpt Log Directory option
	LogOpt = flag.String("l", "./directory.log", "Log directory")
	// WriteOnly Internal use only
	WriteOnly = flag.Bool("w", false, "For system")
)

func main() {

	currentDir, _ := os.Getwd()
	currentDir = filepath.ToSlash(currentDir)
	os.Chdir(filepath.Dir(os.Args[0]))
	log := make(chan []string)
	go readLog(log)

	flag.Parse()

	if len(flag.Args()) != 1 {
		return

		// writeonly version
	} else if *WriteOnly {
		targetPath := flag.Args()[0]
		Log := <-log
		idx, _ := compMatchPath(Log, targetPath)
		rearrangeLog(Log, targetPath, idx)

		// return list of history
	} else if flag.Args()[0] == "-" {
		paths := <-log
		if idx, isOk := compMatchPath(paths, currentDir); isOk {
			paths = del(paths, idx)
		}
		paths = checkExist(paths)
		str := strings.Join(paths, "\n")
		fmt.Println(str)

	} else {
		targetPath := strings.Trim(flag.Args()[0], "\"")
		if !filepath.IsAbs(targetPath) {
			targetPath = filepath.Join(currentDir, targetPath)
		}
		targetPath = filepath.ToSlash(targetPath)

		//  if a directory exists
		if fInfo, err := os.Stat(targetPath); err == nil {
			if ! fInfo.IsDir(){
				return 
			}
			Log := <-log
			if idx, isOk := compMatchPath(Log, targetPath); !isOk {
				makeLog(Log, targetPath)
			} else {
				rearrangeLog(Log, targetPath, idx)
			}

			// not exists
		} else {
			Log := <-log
			pathlist := make(chan []string)
			go findMatch(pathlist, Log, flag.Args()[0])
			paths := checkExist(<-pathlist)

			if idx, isOk := compMatchPath(paths, currentDir); isOk {
				paths = del(paths, idx)
			}

			if len(paths) == 1 {
				idx, _ := compMatchPath(Log, paths[0])
				rearrangeLog(Log, paths[0], idx)
			}
			fmt.Println(strings.Join(paths, "\n"))
		}
	}
	return

}

func checkExist(paths []string) []string {
	maxNum := len(paths)
	for i := 0; i < maxNum; i++ {
		if _, err := os.Stat(paths[i]); err != nil {
			paths = del(paths, i)
			i--
			maxNum--
		}
	}
	return paths
}

var zero string

func del(a []string, i int) []string {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = zero
	a = a[:len(a)-1]
	return a
}

func compMatchPath(paths []string, target string) (int, bool) {
	reg := regexp.MustCompile("^" + target + "$")
	for i, path := range paths {
		if reg.MatchString(path) {
			return i, true
		}
	}
	return 0, false
}

func findMatch(pathlist chan []string, paths []string, target string) {
	// ret := make([]string, 0)
	var ret []string
	reg := regexp.MustCompile(target)
	for _, path := range paths {
		if reg.MatchString(filepath.Base(path)) {
			ret = append(ret, path)
		}
	}
	pathlist <- ret
}

func rearrangeLog(oldpaths []string, path string, index int) {
	fp, _ := os.OpenFile(*LogOpt, os.O_RDWR, 0664)
	fw := bufio.NewWriter(fp)
	defer fp.Close()

	fmt.Fprint(fw, path+"\n")

	if index != 0 {
		fmt.Fprint(fw, strings.Join(oldpaths[:index], "\n")+"\n")
	}
	if index < len(oldpaths)-1 {
		fmt.Fprint(fw, strings.Join(oldpaths[index+1:], "\n"))
	}
	fw.Flush()
}

func makeLog(oldpaths []string, path string) {
	fp, _ := os.OpenFile(*LogOpt, os.O_RDWR|os.O_CREATE, 0664)
	fp.Truncate(0)
	fw := bufio.NewWriter(fp)
	defer fp.Close()

	length := *MaxHistOpt
	if *MaxHistOpt > len(oldpaths) {
		length = len(oldpaths)
	}

	fmt.Fprint(fw, path+"\n")
	fmt.Fprint(fw, strings.Join(oldpaths[:length], "\n"))
	fw.Flush()
}

func readLog(r chan []string) {
	fp, _ := os.Open(*LogOpt)
	count := lineCounter(os.Open(*LogOpt))
	scanner := bufio.NewScanner(fp)
	// var slice []string = make([]string, count)
	var slice = make([]string, count)

	defer func() {
		fp.Close()
		r <- slice
	}()

	for i := 0; i < count; i++ {
		scanner.Scan()
		slice[i] = string(scanner.Text())
	}
}

func lineCounter(r io.Reader, err error) int {
	buf := make([]byte, 8196)
	count := 1
	lineSep := []byte{'\n'}

	for {
		cbuf, err := r.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		count += bytes.Count(buf[:cbuf], lineSep)

		if err == io.EOF {
			break
		}
	}
	return count
}
