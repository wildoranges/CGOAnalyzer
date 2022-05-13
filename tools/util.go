package tools

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var showAll bool = false
var dumpAll bool = false
var coreItem int = 10

type Void struct {

}

type StringSet map[string]Void

var VoidMem Void

func (set StringSet) Insert(item string){
	set[item] = VoidMem
}

type ItemInfo struct{
	total map[string]int
	local map[string]map[string]int
	pathInfo map[string]StringSet
	i ItemKind
}

type ItemKind interface{
	KindText(string) string
	Label() string
}

type Pair struct {
	key   string
	value int
}

type Pairlist []Pair

func (p Pairlist) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Pairlist) Len() int {
	return len(p)
}

func (p Pairlist) Less(i int, j int) bool {
	if p[i].value != p[j].value {
		return p[i].value > p[j].value
	}
	return p[i].key < p[j].key

}

func SortMapByValue(m map[string]int) Pairlist {
	p := make(Pairlist, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

var excludeDir = map[string]Void {
	`test`:VoidMem,
	`vendor`:VoidMem,
}

//GetAllDirs get all dirs in the given dir recursively
func GetAllDirs(dirPth string) ([]string, error) {
	var dirs []string
	var lst []string
	var err error
	s, err := os.Stat(dirPth)
	if err != nil {
		return nil, err
	}
	if !s.IsDir() {
		return nil, err
	}
	isSep := strings.HasSuffix(dirPth, "/")
	if isSep && dirPth != "/" {
		dirPth = dirPth[0 : len(dirPth)-1]
	} else if dirPth == "/" {
		dirPth = ""
	}
	baseDir := filepath.Base(dirPth)
	if _, ok := excludeDir[baseDir]; !ok {
		dirs = append(dirs, dirPth)
		lst = append(lst, dirPth)
	}
	for len(lst) > 0 {
		curdir := lst[0]
		lst = lst[1:]
		di, err := ioutil.ReadDir(curdir)
		if err != nil {
			continue
		} else {
			for _, fi := range di {
				if fi.IsDir() {
					newDirPth := curdir + "/" + fi.Name()
					//newDirPth := filepath.Join(curdir + fi.Name())
					newDirBase := filepath.Base(newDirPth)
					if _, ok := excludeDir[newDirBase]; !ok {
						dirs = append(dirs, newDirPth)
						lst = append(lst, newDirPth)
					}
				}
			}
		}
	}
	return dirs, err
}

func CntRepo(item string, local map[string]map[string]int) int {
	sum := 0
	for _, m := range local {
		if _, ok := m[item]; ok {
			sum++
		}
	}
	return sum
}

func GetMaxRepo(item string, local map[string]map[string]int) (string, int) {
	maxRepo := ""
	maxCnt := 0
	for repo, m := range local {
		if value, ok := m[item]; ok {
			if value > maxCnt {
				maxCnt = value
				maxRepo = repo
			}
		}
	}
	maxRepo = filepath.Base(maxRepo)
	return maxRepo, maxCnt
}

func GetMinRepo(item string, local map[string]map[string]int) (string, int) {
	minRepo := ""
	minCnt := 0
	for repo, m := range local {
		if value, ok := m[item]; ok {
			if minCnt <= 0 {
				minRepo = repo
				minCnt = value
			} else {
				if value < minCnt {
					minCnt = value
					minRepo = repo
				}
			}
		}
	}
	minRepo = filepath.Base(minRepo)
	return minRepo, minCnt
}

func GetAllAvg(item string, local map[string]map[string]int) float64 {
	repoNum := float64(len(local))
	sum := 0.0
	for _, m := range local {
		if value, ok := m[item]; ok {
			sum += float64(value)
		}
	}
	avg := sum / repoNum
	return avg
}

func GetAvg(item string, local map[string]map[string]int) float64 {
	repoNum := 0.0
	sum := 0.0
	for _, m := range local {
		if value, ok := m[item]; ok {
			repoNum += 1.0
			sum += float64(value)
		}
	}
	avg := sum / repoNum
	return avg
}

func GetTopX(x int,item string,local map[string]map[string]int) []string {
	topX := make([]string,0)
	allPair := make(Pairlist,0)
	for repo, m := range local {
		if v, ok := m[item]; ok {
			p := Pair{repo,v}
			allPair = append(allPair, p)
		}
	}
	sort.Sort(allPair)
	for i:=0;i < x&&i < len(allPair);i++{
		topX = append(topX, allPair[i].key)
	}
	return topX
}

func msgReport(w io.Writer,flag bool,format string, a ...interface{}) (int,error){
	if flag {
		return fmt.Fprintf(w,format,a ...)
	}
	return 0,nil
}

func pathRecord(flag bool,info ItemInfo,item string,path string){
	if flag {
		info.pathInfo[item].Insert(path)
	}
}

func SetFlags(showall bool,dumpall bool,core int){
	showAll = showall
	dumpAll = dumpall
	coreItem = core
}