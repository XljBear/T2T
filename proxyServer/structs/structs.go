package structs

import "time"

type TrafficData struct {
	DownlinkInSecond uint64 `json:"downlink_in_second"`
	DownlinkTotal    uint64 `json:"downlink_total"`
	UplinkInSecond   uint64 `json:"uplink_in_second"`
	UplinkTotal      uint64 `json:"uplink_total"`
}
type Link struct {
	UUID     string       `json:"uuid"`
	IP       string       `json:"ip"`
	LinkTime time.Time    `json:"link_time"`
	Traffic  *TrafficData `json:"traffic"`
}
type ByLinkTime []Link

func (a ByLinkTime) Len() int {
	return len(a)
}

func (a ByLinkTime) Less(i, j int) bool {
	return a[i].LinkTime.After(a[j].LinkTime)
}

func (a ByLinkTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
