package pkg

import (
	"errors"
	"fmt"
	"github.com/yangyang5214/btl/pkg/model"
	"os"
	"sort"
	"strings"
)

var specialDelimiter = "~"

type Json2Mind struct {
	inputJson string
	fields    []string
}

func NewJsonGroup(inputJson string, fields []string) *Json2Mind {
	return &Json2Mind{
		inputJson: inputJson,
		fields:    fields,
	}
}

func (j *Json2Mind) ToMarkdown(titleMap map[string][]string, mindMap map[string]*model.Mind, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = f.WriteString("# " + getFilePrefix(j.inputJson))
	_, err = f.WriteString("\n")
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

func (j *Json2Mind) Run() error {
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

	minds := make(map[string]*model.Mind)
	titleMap := make(map[string][]string)

	for k, v := range m {
		arrs := strings.Split(k, specialDelimiter)
		mind := &model.Mind{
			Key:   k,
			Name:  arrs[len(arrs)-1],
			Count: v,
			Level: len(arrs),
		}
		titleMap[arrs[0]] = append(titleMap[arrs[0]], k)
		minds[k] = mind
	}

	err = j.ToMarkdown(titleMap, minds, getFilePrefix(j.inputJson)+".md")
	if err != nil {
		return err
	}

	return nil
}
