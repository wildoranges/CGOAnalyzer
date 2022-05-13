package main

import (
	"anatool/tools"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)


const funcPath string = "Function"
const includePath string = "Include"

func main() {
	pathPtr := flag.String(`path`,`./`,`dir path of the repos`)
	showAll := flag.Bool(`showall`,false,`show output details`)
	dumpPath := flag.Bool(`dumppath`,false,`dump item path`)
	core := flag.Int(`core`,10,`number of core item`)
	topx := flag.Int(`top`,5,`top x repos having the core item`)
	flag.Parse()
	dirpath := *pathPtr
	tools.SetFlags(*showAll,*dumpPath,*core)
	fileinfo, err := ioutil.ReadDir(dirpath)
	fmt.Printf("now getting valid repos ...\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occured while reading dir:%s:\n", dirpath)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	set := make(tools.StringSet)
	for _, si := range fileinfo {
		if si.IsDir() {
			curRepo := si.Name()
			fullpath := filepath.Join(dirpath, curRepo)
			alldirs, err := tools.GetAllDirs(fullpath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error occured while reading dir:%s:\n", fullpath)
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			for _, subpth := range alldirs {
				isvalid := tools.IsGo2CDir(subpth)
				if isvalid {
					set[fullpath] = tools.VoidMem
					break
				}
			}
		}
	}

	var Validdirs []string
	tmpstr := "["

	fmt.Printf("valid repos:\n")
	for p := range set {
		Validdirs = append(Validdirs, p)
		p = "\"" + p + "\""
		fmt.Println(p)
		tmpstr += "\n\t"
		tmpstr += p
		tmpstr += ","
	}
	if tmpstr[len(tmpstr)-1] == ',' {
		tmpstr = tmpstr[0 : len(tmpstr)-1]
	}

	tmpstr += "\n\t]"
	writestr := "{\n\"Go2CDir\":\n\t" + tmpstr + "\n}"
	f, _ := os.Create("Go2CDir.json")
	f.Write([]byte(writestr))
	f.Close()

	fmt.Printf("now parsing valid repos ...\n")
	funcInfo := tools.CntGo2CFunc(Validdirs)
	err = tools.DumpAll2Csv(funcInfo, "./func_cnt.csv", funcPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	err = tools.DumpDetail2Csv(funcInfo, "func_detail.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} 

	headerInfo := tools.CntGo2CInclude(Validdirs)
	err = tools.DumpAll2Csv(headerInfo, "./include_cnt.csv", includePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	err = tools.DumpDetail2Csv(headerInfo, "include_detail.csv")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	err = tools.DumpTopX(*topx,funcInfo,fmt.Sprintf("func_%d_top%d.csv",*core,*topx))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	err = tools.DumpTopX(*topx,headerInfo,fmt.Sprintf("include_%d_top%d.csv",*core,*topx))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Printf("finish parsing!\n")
}
