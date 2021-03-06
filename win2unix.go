package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	pwd = flag.String("pwd", "./", "pwd")
)

func main() {
	flag.Parse()

	//f, _ := os.Open(*pwd)
	//defer f.Close()
	n := strings.LastIndex(*pwd, "/")
	tmp := substr(*pwd, n)
	path_split := strings.Split(*pwd, "/")
	file_name := path_split[(len(path_split) - 1)]

	new_path := tmp + "/replace_files"
	if ok, err := replace_word("\r\n", "\n", *pwd, new_path, file_name); ok {
		fmt.Println("find")
		//fmt.Println("find")
	} else {
		fmt.Println("no find")
		fmt.Println(err)
	}

}

func replace_word(search string, replace string, path string, new_path string, file_name string) (bool, error) {

	var err error
	if lines, err := reader_line(path); err == nil {
		for key, line := range lines {
			lines[key] = strings.Replace(string(line), search, replace, -1)
		}
		if err := write_line(new_path, file_name, lines); err == nil {
			return true, err
		}
	}

	return false, err
}

func reader_line(path string) (lines []string, err error) {
	file, _ := os.Open(path)
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, is_prefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		if is_prefix {
			break
		}

		buffer := bytes.NewBuffer(make([]byte, 1024))
		buffer.Write(line)
		lines = append(lines, buffer.String())

	}
	if err == io.EOF {
		err = nil
	}
	return
}

func write_line(path string, file_name string, lines []string) (err error) {
	var file *os.File
	if err = os.Mkdir(path, 0777); err != nil {
		fmt.Println(err)
		//return
	}
	if file, err = os.Create(path + "/" + file_name); err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for _, line := range lines {

		if _, err := file.WriteString(strings.TrimSpace(line) + "\n"); err != nil {
			fmt.Println(err)
			break
		}
	}
	return
}

func substr(str string, length int) string {
	bs := []byte(str)[:length]
	bl := 0
	for i := len(bs) - 1; i >= 0; i-- {
		switch {
		case bs[i] >= 0 && bs[i] <= 127:
			return string(bs[:i+1])
		case bs[i] >= 128 && bs[i] <= 191:
			bl++
		case bs[i] >= 192 && bs[i] <= 253:
			cl := 0
			switch {
			case bs[i]&252 == 252:
				cl = 6
			case bs[i]&248 == 248:
				cl = 5
			case bs[i]&240 == 240:
				cl = 4
			case bs[i]&224 == 224:
				cl = 3
			default:
				cl = 2
			}
			if bl+1 == cl {
				return string(bs[:i+cl])
			}
			return string(bs[:i])
		}
	}
	return ""
}
