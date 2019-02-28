package service

import (
	"github.com/imdevlab/g"
	"github.com/imdevlab/vgo/analyze/misc"
	"github.com/imdevlab/vgo/proto/pinpoint/thrift/trace"
	"github.com/imdevlab/vgo/util"
	"go.uber.org/zap"
)

// 应用拓扑图计算

// SrvMaps ...
type SrvMaps struct {
	appType int
	srvMaps map[string]*SrvMap
	dbMaps  map[int16]*DBMap // 数据库信息
}

// NewServiceMapsCounter ...
func NewServiceMapsCounter() *SrvMaps {
	return &SrvMaps{
		srvMaps: make(map[string]*SrvMap),
		dbMaps:  make(map[int16]*DBMap),
	}
}

// srvMapsCounter ...
func (s *SrvMaps) srvMapsCounter(appName string, appType int, parentName string, parentType int, elapsed int, isError int, events []*trace.TSpanEvent, chunkEvents []*trace.TSpanEvent) error {
	s.appType = appType
	srv, ok := s.srvMaps[parentName]
	if !ok {
		srv = NewSrvMap()
		srv.parentType = parentType
		s.srvMaps[parentName] = srv
	}

	srv.count++
	if isError != 0 {
		srv.errCount++
	}
	srv.totalelapsed += elapsed

	for _, event := range events {
		isDB := false
		if event.ServiceType == util.MYSQL_EXECUTE_QUERY {
			isDB = true
		} else if event.ServiceType == util.REDIS {
			isDB = true
		} else if event.ServiceType == util.ORACLE_EXECUTE_QUERY {
			isDB = true
		} else if event.ServiceType == util.POSTGRESQL_EXECUTE_QUERY {
			isDB = true
		}
		if isDB {
			db, ok := s.dbMaps[event.ServiceType]
			if !ok {
				db = NewDBMap()
				s.dbMaps[event.ServiceType] = db
			}
			db.Count++
			if event.GetExceptionInfo() != nil {
				db.ErrCount++
			}
			db.Totale += int(event.EndElapsed)
		}
	}

	for _, event := range chunkEvents {
		isDB := false
		if event.ServiceType == util.MYSQL_EXECUTE_QUERY {
			isDB = true
		} else if event.ServiceType == util.REDIS {
			isDB = true
		} else if event.ServiceType == util.ORACLE_EXECUTE_QUERY {
			isDB = true
		} else if event.ServiceType == util.POSTGRESQL_EXECUTE_QUERY {
			isDB = true
		}
		if isDB {
			db, ok := s.dbMaps[event.ServiceType]
			if !ok {
				db = NewDBMap()
				s.dbMaps[event.ServiceType] = db
			}
			db.Count++
			if event.GetExceptionInfo() != nil {
				db.ErrCount++
			}
			db.Totale += int(event.EndElapsed)
		}
	}

	return nil
}

// srvMapsRecord ...
func (s *SrvMaps) srvMapsRecord(app *App, recordTime int64) error {
	for parentName, srv := range s.srvMaps {
		query := gAnalyze.cql.Session.Query(misc.InserServiceMapRecord,
			app.AppName,
			recordTime,
			s.appType,
			parentName,
			srv.parentType,
			srv.count,
			srv.errCount,
			srv.totalelapsed,
		)
		if err := query.Exec(); err != nil {
			g.L.Warn("srvMapsRecord error", zap.String("error", err.Error()), zap.String("sql", query.String()))
		}
	}

	for dbType, db := range s.dbMaps {
		query := gAnalyze.cql.Session.Query(misc.InserDBMapRecord,
			app.AppName,
			recordTime,
			s.appType,
			dbType,
			db.Count,
			db.ErrCount,
			db.Totale,
		)
		if err := query.Exec(); err != nil {
			g.L.Warn("dbMapsRecord error", zap.String("error", err.Error()), zap.String("sql", query.String()))
		}
	}
	return nil
}

// SrvMap ...
type SrvMap struct {
	parentType   int
	totalelapsed int
	count        int
	errCount     int
}

// NewSrvMap ...
func NewSrvMap() *SrvMap {
	return &SrvMap{}
}

// DBMap ...
type DBMap struct {
	Totale   int
	Count    int
	ErrCount int
}

// NewDBMap ...
func NewDBMap() *DBMap {
	return &DBMap{}
}
