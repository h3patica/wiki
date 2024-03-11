package main

import (
	"os"
	"path/filepath"
	"strings"
	//	"fmt"
)

func prune(path string) (dir, file string) {
	dir, file = filepath.Split(path)
	file, _ = strings.CutSuffix(file, ".html")
	return dir, file
}

func end(name string) string {
	_, name = prune(name)
	nList := strings.Split(name, ".")
	return nList[len(nList)-1]
}

func find(path string, exp func(s1, s2 string) bool) []string {
	var result []string
	dir, file := prune(path)
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		fName, _ := strings.CutSuffix(f.Name(), ".html")
		if exp(fName, file) {
			result = append(result, f.Name())
		}
	}
	return result
}

func getAllChildren(path string) []string {
	children := find(path, func(s1, s2 string) bool {
		return (strings.Contains(s1, s2) && s1 != s2)
	})
	return children
}

func getChildren(path string, span int) []string {
	children := find(path, func(s1, s2 string) bool {
		return (len(strings.Split(s1, ".")) <= len(strings.Split(s2, ".")) + span && strings.Contains(s1, s2) && s1 != s2)
	})
	return children
}

func getAllParents(path string) []string {
	var parents []string
	_, file := prune(path)
	pList := strings.Split(file, ".")
	pList = pList[0:len(pList)-1]
	for i := 0; i < len(pList); i++ {
		p := pList[i]
		for j:= i-1; j >= 0; j-- {
			p = pList[j] + "." + p
		}
		parents = append(parents, p + ".html")
	}
	return parents
}

//func main() {
//	fmt.Println(getChildren("./site/index", 2))
//}
