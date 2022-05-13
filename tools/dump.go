package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func Dump2Csv(total map[string]int, path string) error {
	var err error
	filePath := path
	totalList := SortMapByValue(total)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	column := "Func,Cnt\n"
	f.Write([]byte(column))
	for _, p := range totalList {
		msgReport(os.Stdout,showAll,"item:%s,cnt:%d in %s\n", p.key, p.value, filePath)
		writestr := fmt.Sprintf("%s,%d\n", p.key, p.value)
		_, err = f.Write([]byte(writestr))
	}
	return err
}

func DumpAll2Csv(Info ItemInfo, path string, subDir string) error {
	var err error
	err = Dump2Csv(Info.total, path)
	dir := filepath.Dir(path)
	infoPath := filepath.Join(dir, subDir)
	for repo, m := range Info.local {
		repoBase := filepath.Base(repo)
		repoInfoPath := filepath.Join(infoPath, repoBase)
		os.MkdirAll(repoInfoPath, os.ModePerm)
		csvName := repoBase + ".csv"
		csvPath := filepath.Join(repoInfoPath, csvName)
		err = Dump2Csv(m, csvPath)
	}
	return err
}

func DumpDetail2Csv(Info ItemInfo, path string) error {
	var err error
	var f *os.File
	filePath := path
	totalList := SortMapByValue(Info.total)
	f, err = os.Create(filePath)
	if err != nil {
		msgReport(os.Stderr,true,err.Error()+"\n")
	}
	defer f.Close()
	column := "item,cnt,repo_num,max,max_repo,min,min_repo,all_avg,avg,kind\n"
	_, err = f.Write([]byte(column))
	for _, p := range totalList {
		item := p.key
		value := p.value
		repoNum := CntRepo(item, Info.local)
		maxRepo, max := GetMaxRepo(item, Info.local)
		minRepo, min := GetMinRepo(item, Info.local)
		allAvg := GetAllAvg(item, Info.local)
		avg := GetAvg(item, Info.local)
		kind := Info.i.KindText(item)
		line := fmt.Sprintf("%s,%d,%d,%d,%s,%d,%s,%.3f,%.3f,%s\n", item, value, repoNum, max, maxRepo, min, minRepo, allAvg, avg, kind)
		_, err = f.Write([]byte(line))
	}

	if dumpAll {
		fn := filepath.Join(filepath.Dir(path),Info.i.Label()+`_path.csv`) 
		f2, err := os.Create(fn)
		if err != nil {
			msgReport(os.Stderr,true,err.Error()+"\n")
		}
		defer f2.Close()
		column := "item,path\n"
		f2.Write([]byte(column))
		for name, set := range Info.pathInfo{
			for path := range set{
				line := fmt.Sprintf("%s,%s\n",name,path)
				f2.Write([]byte(line))
			}
		}
	}
	return err
}

func DumpTopX(x int,Info ItemInfo, path string)(err error) {
	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	column := "item,cnt"
	for i:=1;i <= x;i++{
		column += fmt.Sprintf(",top%d,cnt",i)
	}
	column += "\n"
	f.Write([]byte(column))
	totalList := SortMapByValue(Info.total)
	length := len(totalList)
	for i:=0;i < coreItem && i < length;i++{
		p := totalList[i]
		Topx := GetTopX(x,p.key,Info.local)
		base := fmt.Sprintf("%s,%d",p.key,p.value)
		for _, repo := range Topx{
			cnt := Info.local[repo][p.key]
			base += fmt.Sprintf(",%s,%d",repo,cnt)
		}
		base += "\n"
		f.Write([]byte(base))
	}
	return err
}