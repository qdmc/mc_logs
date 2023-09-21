package mc_logs

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

type DefaultJsonFormatter struct {
}

func (f DefaultJsonFormatter) ToFields(entry *logrus.Entry) map[string]interface{} {
	timeLayout := "2006-01-02 15:04:05:000"
	return map[string]interface{}{
		logrus.FieldKeyTime:  time.Unix(entry.Time.Unix(), 0).Format(timeLayout),
		logrus.FieldKeyMsg:   entry.Message,
		logrus.FieldKeyLevel: entry.Level.String(),
		"timeStamp":          entry.Time.UnixNano(),
	}
}

func (f DefaultJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fs := f.ToFields(entry)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			fs[k] = v.Error()
		default:
			fs[k] = v
		}
	}
	return json.Marshal(fs)
}
