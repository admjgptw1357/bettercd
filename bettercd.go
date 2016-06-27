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
	MaxHistOpt = flag.Int("i", 200, "Max number of history")
	LogOpt     = flag.String("l", "./directory.log", "Log directory")
	WriteOnly  = flag.Bool("w", false, "For system")
)

func main() {

	current_dir, _ := os.Getwd()
	current_dir = filepath.ToSlash(current_dir)
	os.Chdir(filepath.Dir(os.Args[0]))
	log := make(chan []string)
	go read_log(log)

	flag.Parse()

	if len(flag.Args()) != 1 {
		return

	}else if *WriteOnly{
		target_path := flag.Args()[0]
		Log := <-log
		idx, _ := comp_match_path(Log, target_path);
		rearrange_log(Log,target_path,idx)

	} else if flag.Args()[0] == "-" {
		paths := <-log
		if idx,isOk := comp_match_path(paths, current_dir); isOk{
			paths = del(paths,idx)
		}
		paths = check_exist(paths)
		str := strings.Join(paths, "\n")
		fmt.Println(str)

	} else {
		target_path := strings.Trim(flag.Args()[0],"\"")
		if !filepath.IsAbs(target_path) {
			target_path = filepath.Join(current_dir, target_path)
		}
		target_path = filepath.ToSlash(target_path)
		// todo
		// ignore case
		// file tono kubetu


		if _, err := os.Stat(target_path); err == nil {
			Log := <-log
			if idx, isOk := comp_match_path(Log, target_path); !isOk {
				make_log(Log, target_path)
			} else {
				rearrange_log(Log, target_path, idx)
			}

		} else {
			Log := <-log
			pathlist := make(chan []string)
			go find_match(pathlist, Log, flag.Args()[0])
			paths := check_exist(<-pathlist)

			if idx, isOk := comp_match_path(paths, current_dir); isOk{
				paths = del(paths,idx)
			}

			if len(paths) == 1 {
				idx, _ := comp_match_path(Log, paths[0])
				rearrange_log(Log, paths[0], idx)
			}
			fmt.Println(strings.Join(paths, "\n"))
		}
	}
	return

}

func check_exist(paths []string) []string {
	max_num := len(paths)
	for i := 0; i < max_num; i++ {
		if _, err := os.Stat(paths[i]); err != nil {
			paths = del(paths, i)
			i--
			max_num--
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

func comp_match_path(paths []string, target string) (int, bool) {
	reg := regexp.MustCompile("^" + target + "$")
	for i, path := range paths {
		if reg.MatchString(path) {
			return i, true
		}
	}
	return 0, false
}

func find_match(pathlist chan []string, paths []string, target string) {
	ret := make([]string, 0)
	reg := regexp.MustCompile(target)
	for _, c_path := range paths {
		if reg.MatchString(filepath.Base(c_path)) {
			ret = append(ret, c_path)
		}
	}
	pathlist <- ret
}

func rearrange_log(oldpaths []string, path string, index int) {
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

func make_log(oldpaths []string, path string) {
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

func read_log(r chan []string) {
	fp, _ := os.Open(*LogOpt)
	count := line_counter(os.Open(*LogOpt))
	scanner := bufio.NewScanner(fp)
	var slice []string = make([]string, count)

	defer func() {
		fp.Close()
		r <- slice
	}()

	for i := 0; i < count; i++ {
		scanner.Scan()
		slice[i] = string(scanner.Text())
	}
}

func line_counter(r io.Reader, err error) int {
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
