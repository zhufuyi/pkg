package snowflake

// This package provides unique id in distribute system
// the algorithm is inspired by Twitter's famous snowflake
// its link is: https://github.com/twitter/snowflake/releases/tag/snowflake-2010
//

// 0               41	           51	 	 64
// +---------------+----------------+------------+
// |timestamp(ms)  | worker node id | sequence	 |
// +---------------+----------------+------------+

// Copyright (C) 2016 by zheng-ji.info

// 修改于 https://github.com/zheng-ji/goSnowFlake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	epoch = 1464613955000 // 毫秒

	sequenceBits = 12 // Num of Sequence Bits
	workerIDBits = 10 // Num of WorkerId Bits

	workerIDShift  = 12 // workerIDBits  = sequenceBits
	timestampShift = 22 // timestampShift = sequenceBits + workerIDBits
)

// IDWorker Struct
type IDWorker struct {
	lastTimeStamp int64

	workerID     int64
	sequence     int64
	sequenceMask int64
	maxWorkerID  int64

	lock *sync.Mutex
}

// NewIDWorker Func: Generate NewIdWorker with Given workerId
func NewIDWorker(workerID int64) (iw *IDWorker, err error) {
	iw = new(IDWorker)

	iw.maxWorkerID = getMaxWorkerID()

	if workerID > iw.maxWorkerID || workerID < 0 {
		return nil, errors.New("worker not fit")
	}

	iw.workerID = workerID
	iw.lastTimeStamp = -1
	iw.sequence = 0
	iw.sequenceMask = getSequenceMask()
	iw.lock = new(sync.Mutex)
	return iw, nil
}

// 最大支持的WorkerID
// 1023
func getMaxWorkerID() int64 {
	return -1 ^ -1<<workerIDBits
}

// 最大支持的Sequence
// 4095
func getSequenceMask() int64 {
	return -1 ^ -1<<sequenceBits
}

// 毫秒数
func (iw *IDWorker) genTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func (iw *IDWorker) reGenTimestamp(last int64) int64 {
	ts := iw.genTimestamp()
	for {
		if ts < last {
			ts = iw.genTimestamp()
		} else {
			break
		}
	}
	return ts
}

// NextID return the id
func (iw *IDWorker) NextID() (ts int64, err error) {
	iw.lock.Lock()
	defer iw.lock.Unlock()

	ts = iw.genTimestamp()

	if ts == iw.lastTimeStamp {
		iw.sequence = (iw.sequence + 1) & iw.sequenceMask

		if iw.sequence == 0 {
			ts = iw.reGenTimestamp(ts)
		}
	} else {
		iw.sequence = 0
	}

	if ts < iw.lastTimeStamp {
		err = errors.New("Clock moved backwards, Refuse gen id")
		return 0, err
	}

	iw.lastTimeStamp = ts

	ts = (ts-epoch)<<timestampShift | iw.workerID<<workerIDShift | iw.sequence

	return ts, nil
}

// ---------------------------------------------------------------------------------------

var iw *IDWorker

// Init 初始化SnowFlake
func Init(workID int64) error {
	var err error
	iw, err = NewIDWorker(workID)
	if err != nil {
		return err
	}

	return nil
}

// NewID create new id
func NewID() int64 {
	id, err := iw.NextID()
	if err != nil {
		time.Sleep(time.Microsecond)
		// try again
		id, err = iw.NextID()
		if err != nil {
			time.Sleep(time.Microsecond)
			// retry
			id, err = iw.NextID()
			if err != nil {
				fmt.Printf("create id failed, %s\n", err.Error())
				return -1
			}
		}
	}

	return id
}
