//snowflake algorithm

package util

import (
	"errors"
	"sync"
	"time"
)

const (
	sequenceBits       = 8
	nodeIdBits         = 16
	sequenceMask       = (1 << sequenceBits) - 1
	nodeIdLeftShift    = sequenceBits
	timestampLeftShift = sequenceBits + nodeIdBits
	timeBits           = 63 - timestampLeftShift
	workIdMax          = 1 << nodeIdBits
)

type SFSettings struct {
	NodeId    func() (uint16, error)
	StartTime time.Time
}

type SnowFlake struct {
	sync.Mutex
	sequence   uint16
	latestTime int64
	startTime  int64
	nodeId     uint16
}

func NewSnowFlake(st *SFSettings) (snowFlake *SnowFlake, err error) {

	var startTime time.Time
	if !st.StartTime.IsZero() {
		if st.StartTime.After(time.Now()) {
			err = errors.New("start time cannot be bigger than now")
			return
		}
		startTime = st.StartTime
	} else {
		startTime = time.Date(2015, 9, 7, 0, 0, 0, 0, time.Local)
	}

	var nodeId uint16
	if st.NodeId != nil {
		nodeId, err = st.NodeId()
		if int(nodeId) > workIdMax {
			err = errors.New("nodeId bigger than maxId")
		}
	} else {
		nodeId, err = Low16BitsPrivateIP4()
	}

	if err != nil {
		return
	}

	snowFlake = &SnowFlake{
		nodeId:    nodeId,
		startTime: toMilliSecondTime(startTime),
		sequence:  (1 << sequenceBits) - 1,
	}

	return
}

func (self *SnowFlake) NextId() (uint64, error) {
	self.Lock()
	defer self.Unlock()
	currentMS := CurrentMillis()
	if currentMS < self.latestTime {
		return 0, errors.New("Clock is moving backwards")
	}
	if currentMS == self.latestTime {
		self.sequence = (self.sequence + 1) & sequenceMask
		if 0 == self.sequence {
			currentMS = waitUtilNextTime(currentMS)
			if currentMS-self.startTime >= 1<<timeBits {
				return 0, errors.New("over time limit")
			}
		}
	} else {
		self.sequence = 0
	}

	self.latestTime = currentMS

	return uint64(currentMS-self.startTime)<<timestampLeftShift +
		uint64(self.nodeId)<<nodeIdLeftShift +
		uint64(self.sequence), nil
}

func waitUtilNextTime(latestTime int64) int64 {
	currentMS := CurrentMillis()
	for latestTime < currentMS {
		latestTime = currentMS
	}
	return latestTime
}

func toMilliSecondTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}
