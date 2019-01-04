package service

import (
	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
)

// SpanSQLs ...
type SpanSQLs struct {
	sqls map[int32]*SpanSQL
}

// NewSpanSQLs ...
func NewSpanSQLs() *SpanSQLs {
	return &SpanSQLs{
		sqls: make(map[int32]*SpanSQL),
	}
}

// sqlCounter ...
func (spanSQLs *SpanSQLs) sqlCounter(events []*trace.TSpanEvent, chunkEvents []*trace.TSpanEvent) error {
	for _, event := range events {
		ans := event.GetAnnotations()
		if len(ans) >= 0 {
			for _, an := range ans {
				if an.GetKey() == 20 {
					sql, ok := spanSQLs.sqls[an.GetValue().GetIntValue()]
					if !ok {
						sql = NewSpanSQL()
						spanSQLs.sqls[an.GetValue().GetIntValue()] = sql
					}
					sql.count++
					elapsed := int(event.EndElapsed - event.StartElapsed)
					sql.elapsed += elapsed
					if elapsed > sql.maxElapsed {
						sql.maxElapsed = sql.elapsed
					}

					if sql.minElapsed == 0 || sql.minElapsed > elapsed {
						sql.minElapsed = elapsed
					}

					// 是否有异常抛出
					if event.GetExceptionInfo() != nil {
						sql.errCount++
					}

					sql.averageElapsed = sql.elapsed / sql.count
				}
			}
		}

	}
	for _, event := range chunkEvents {
		ans := event.GetAnnotations()
		if len(ans) >= 0 {
			for _, an := range ans {
				if an.GetKey() == 20 {
					sql, ok := spanSQLs.sqls[an.GetValue().GetIntValue()]
					if !ok {
						sql = NewSpanSQL()
						spanSQLs.sqls[an.GetValue().GetIntValue()] = sql
					}
					sql.count++
					elapsed := int(event.EndElapsed - event.StartElapsed)
					sql.elapsed += elapsed
					if elapsed > sql.maxElapsed {
						sql.maxElapsed = sql.elapsed
					}

					if sql.minElapsed == 0 || sql.minElapsed > elapsed {
						sql.minElapsed = elapsed
					}

					// 是否有异常抛出
					if event.GetExceptionInfo() != nil {
						sql.errCount++
					}

					sql.averageElapsed = sql.elapsed / sql.count
				}
			}
		}
	}
	return nil
}

// SpanSQL ...
type SpanSQL struct {
	averageElapsed int
	elapsed        int
	count          int
	errCount       int
	minElapsed     int
	maxElapsed     int
}

// NewSpanSQL ...
func NewSpanSQL() *SpanSQL {
	return &SpanSQL{}
}
