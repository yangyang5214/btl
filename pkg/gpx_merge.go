package pkg

import (
	"bufio"
	"github.com/pkg/errors"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type GpxMerge struct {
	CurrentDir string
}

type GpxFile struct {
	start   []string
	end     []string
	content []string
}

func (g *GpxMerge) Run() error {
	dirs, err := os.ReadDir(g.CurrentDir)
	if err != nil {
		return errors.WithStack(err)
	}

	var gpxFiles []string
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		if strings.HasSuffix(dir.Name(), ".gpx") {
			gpxFiles = append(gpxFiles, path.Join(g.CurrentDir, dir.Name()))
		}
	}

	if len(gpxFiles) == 0 {
		log.Info("not find gpx files in current directory")
		return nil
	}
	gpxFiles, err = g.sortByDate(gpxFiles)
	if err != nil {
		return errors.WithStack(err)
	}

	resultFile, err := os.Create(path.Join(g.CurrentDir, "result.gpx"))
	if err != nil {
		return errors.WithStack(err)
	}
	for index, gFile := range gpxFiles {
		r := g.parseTrkseg(gFile)

		//start
		if index == 0 {
			_, err = resultFile.WriteString(strings.Join(r.start, "\n"))
			if err != nil {
				return errors.WithStack(err)
			}
			_, _ = resultFile.WriteString("\n")
		}

		//content
		_, err = resultFile.WriteString(strings.Join(r.content, "\n"))
		if err != nil {
			return errors.WithStack(err)
		}
		_, _ = resultFile.WriteString("\n")

		//end
		if index == len(gpxFiles)-1 {
			_, err = resultFile.WriteString(strings.Join(r.end, "\n"))
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

// sortByDate is sort gpx files by date
func (g *GpxMerge) sortByDate(files []string) ([]string, error) {
	dateMap := make(map[int64]string)
	for _, f := range files {
		date, err := g.getDate(f)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		dateMap[date] = f
	}

	//get sorted files
	keys := make([]int64, 0, len(dateMap))
	for key := range dateMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	var result []string
	for _, k := range keys {
		result = append(result, dateMap[k])
	}

	return result, nil
}

func (g *GpxMerge) getDate(f string) (int64, error) {
	file, err := os.Open(f)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "<time>") {
			line = strings.Replace(line, "<time>", "", 1)
			line = strings.Replace(line, "</time>", "", 1)

			t, errTime := time.Parse(time.RFC3339, line)
			if errTime != nil {
				return 0, errTime
			}
			return t.Unix(), nil
		}
	}
	return 0, nil
}

// parseTrkseg parse <trkseg>xxx</trkseg> xxx 内容
func (g *GpxMerge) parseTrkseg(f string) *GpxFile {
	bytes, _ := os.ReadFile(f)
	content := string(bytes)
	lines := strings.Split(content, "\n")
	var start, end int
	for index, line := range lines {
		line = strings.TrimLeft(line, " ")
		if strings.HasSuffix(line, "<trkseg>") {
			start = index + 1
		} else if strings.HasSuffix(line, "</trkseg>") {
			end = index
		}
	}
	return &GpxFile{
		start:   lines[:start],
		content: lines[start:end],
		end:     lines[end:],
	}
}
