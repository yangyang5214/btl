package pkg

import (
	"encoding/csv"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type ChargingCDZ struct {
	csvFile string
	log     *log.Helper
}

func NewChargingCDZ(csvFile string) *ChargingCDZ {
	return &ChargingCDZ{
		csvFile: csvFile,
		log:     log.NewHelper(log.DefaultLogger),
	}
}

type PoleOrder struct {
	StationId string
	PoleID    string
	StartTime string
	EndTime   string
}

type PoleUsage struct {
	StationId  string
	PoleID     string
	Date       string
	Hours      []string
	UsageHours int
}

func parserOrder(csvFile string) ([]*PoleOrder, error) {
	csvF, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}

	var orders []*PoleOrder
	reader := csv.NewReader(csvF)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		orders = append(orders, &PoleOrder{
			StationId: record[0],
			PoleID:    record[2],
			StartTime: record[3],
			EndTime:   record[4],
		})
	}
	return orders, nil
}

func strconvAtoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func (c *ChargingCDZ) Run() error {
	orders, err := parserOrder(c.csvFile)
	if err != nil {
		return err
	}

	//按天分组
	m := make(map[string][]*PoleOrder)
	for _, order := range orders {
		dateStr := order.StartTime[:10]
		k := order.PoleID + "xxx" + dateStr
		m[k] = append(m[k], order)
	}

	var usages []*PoleUsage

	poleStationMap := make(map[string]string)
	for k, poleOrders := range m {
		hoursMap := make(map[int]bool)
		for _, order := range poleOrders {
			poleStationMap[order.PoleID] = order.StationId
			startHour := order.StartTime[11:13]
			endHour := order.EndTime[11:13]
			if order.StartTime[0:10] != order.EndTime[0:10] {
				hoursMap[strconvAtoi(startHour)] = true
				continue //跨天了
			}
			for i := strconvAtoi(startHour); i < strconvAtoi(endHour)+1; i++ {
				hoursMap[i] = true
			}
		}

		var hours []string
		for h := range hoursMap {
			hours = append(hours, strconv.Itoa(h))
		}

		slices.Sort(hours)
		arrs := strings.Split(k, "xxx")
		usages = append(usages, &PoleUsage{
			StationId:  poleStationMap[arrs[0]],
			PoleID:     arrs[0],
			Date:       arrs[1],
			Hours:      hours,
			UsageHours: len(hours),
		})
	}

	outF, err := os.Create("result.csv")
	if err != nil {
		return err
	}
	csvW := csv.NewWriter(outF)
	_ = csvW.Write([]string{"StationId", "PoleID", "Date", "Hours", "Usage"})
	for _, item := range usages {
		_ = csvW.Write([]string{
			item.StationId, item.PoleID, item.Date,
			strings.Join(item.Hours, ","),
			strconv.Itoa(item.UsageHours)})
	}
	return nil
}
