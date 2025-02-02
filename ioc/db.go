package ioc

import (
	"github.com/chenmuyao/go-bootcamp/config"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/pkg/gormx"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	prom "github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
)

func InitDB(l logger.Logger) *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Cfg.DB.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	},
	)
	if err != nil {
		panic("failed to connect database")
	}

	err = db.Use(
		prometheus.New(prometheus.Config{
			DBName:           "wetravel",
			RefreshInterval:  15,
			MetricsCollector: []prometheus.MetricsCollector{
				// &prometheus.MySQL{
				// 	VariableNames: []string{"thread_running"},
				// },
			},
		}))
	if err != nil {
		panic(err)
	}

	db.Use(tracing.NewPlugin(tracing.WithoutMetrics(), tracing.WithDBName("wetravel")))

	err = db.Use(gormx.NewCallbacks(prom.SummaryOpts{
		Namespace: "my_company",
		Subsystem: "wetravel",
		Name:      "gorm_db",
		Help:      "GORM Request data",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
		ConstLabels: prom.Labels{
			"instance_id": "instance",
		},
	}))
	if err != nil {
		panic(err)
	}

	// TODO: Replace by sql migration
	err = dao.InitTable(db)
	if err != nil {
		panic("failed to init tables")
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(s string, i ...interface{}) {
	g("", logger.Field{
		Key:   "args",
		Value: i,
	})
}
