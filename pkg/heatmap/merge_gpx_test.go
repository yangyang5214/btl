package heatmap

import "testing"

var mg *MergeGpx

func init() {
	mg = NewMergeGpx()
}

func TestGetActivityIds(t *testing.T) {
	ids, err := mg.GetActivityIds()
	if err != nil {
		panic(err)
	}
	t.Log(len(ids))
	t.Log(ids[0:1])
}

func TestDownloadActivity(t *testing.T) {
	bytes, err := mg.DownloadActivity("17258a1d2300538a0f325449c2020cc1")
	if err != nil {
		panic(err)
	}
	t.Log(string(bytes))
}
