package generator

import (
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/wwq1988/leaf/pkg/conf"
)

const (
	seqMagic      = 0x2ff
	workerIdMagic = 0x2ff
)

var seq uint32

type snowflake struct {
	workerId uint64
	step     uint32
}

func NewSnowflake() (Generator, error) {
	workerIdStr := os.Getenv("WORKERID")
	workerIdInt, err := strconv.Atoi(workerIdStr)
	if err != nil {
		return nil, errors.Wrapf(err, "workerid", workerIdStr)
	}
	workerId := (workerIdInt & workerIdMagic) << 10
	return &snowflake{workerId: uint64(workerId), step: conf.GetSnowflakeStep()}, nil
}

func (s *snowflake) Generate(key string) ([]uint64, error) {
	timestamp := uint64(time.Now().Unix()) << 23
	var old uint32
	var new uint32
	for {
		old = seq
		new = old + s.step
		if atomic.CompareAndSwapUint32(&seq, old, new) {
			break
		}
	}
	ids := make([]uint64, 0, s.step)
	for i := old + 1; i < new; i++ {
		ids = append(ids, timestamp|s.workerId|(uint64(i)&seqMagic))
	}

	return ids, nil
}

func (s *snowflake) Stop() {
}
