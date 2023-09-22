package pkg

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

var specialDelimiter = "~"

type JsonGroup struct {
	inputJson string
	fields    []string
	filename  string
}

type Mind struct {
	key   string
	name  string
	count int
	level int
}

func (m *Mind) String() string {
	s := strings.Builder{}
	for i := 0; i < m.level+1; i++ {
		s.WriteString("#")
	}
	s.WriteString(" ")

	s.WriteString(fmt.Sprintf("%s (%d)", m.name, m.count))

	s.WriteString("\n")
	s.WriteString("\n")
	return s.String()
}

func NewJsonGroup(inputJson string, fields []string) *JsonGroup {
	return &JsonGroup{
		inputJson: inputJson,
		fields:    fields,
	}
}

func (j *JsonGroup) ToMarkdown(titleMap map[string][]string, mindMap map[string]*Mind, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf("# %s", name))
	_, err = f.WriteString("\n")

	for _, subTitles := range titleMap {
		sort.Slice(subTitles, func(i, j int) bool {
			return len(subTitles[i]) < len(subTitles[j])
		})
		for _, subTitle := range subTitles {
			mind := mindMap[subTitle]
			_, err = f.WriteString(mind.String())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (j *JsonGroup) Run() error {
	datas, err := readJsonLineFile(j.inputJson)
	if err != nil {
		return err
	}

	m := make(map[string]int)

	for _, data := range datas {
		for i := 0; i < len(j.fields); i++ {
			var values []string
			for _, field := range j.fields[0 : i+1] {
				v, ok := data[field]
				if !ok {
					return errors.New(fmt.Sprintf("field [%s] not found", field))
				}
				values = append(values, fmt.Sprintf("%v", v))
			}
			m[strings.Join(values, specialDelimiter)]++
		}
	}

	minds := make(map[string]*Mind)
	titleMap := make(map[string][]string)

	for k, v := range m {
		arrs := strings.Split(k, specialDelimiter)
		mind := &Mind{
			key:   k,
			name:  arrs[len(arrs)-1],
			count: v,
			level: len(arrs),
		}
		titleMap[arrs[0]] = append(titleMap[arrs[0]], k)
		minds[k] = mind
	}

	err = j.ToMarkdown(titleMap, minds, j.getMarkdownName())
	if err != nil {
		return err
	}

	return nil
}

func (j *JsonGroup) getMarkdownName() string {
	paths := strings.Split(j.inputJson, "/")
	jsonFileName := paths[len(paths)-1]
	arrs := strings.Split(jsonFileName, ".")
	if len(arrs) == 1 {
		return jsonFileName + "." + ".md"
	}
	arrs[len(arrs)-1] = "md"
	return strings.Join(arrs, ".")
}
