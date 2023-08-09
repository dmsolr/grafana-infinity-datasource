package infinity

import (
	"context"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/yesoreyeram/grafana-infinity-datasource/pkg/models"
)

type CustomMeta struct {
	Query                  models.Query  `json:"query"`
	Data                   any           `json:"data"`
	ResponseCodeFromServer int           `json:"responseCodeFromServer"`
	Duration               time.Duration `json:"duration"`
	Error                  string        `json:"error"`
}

func GetDummyFrame(query models.Query) *data.Frame {
	frameName := query.RefID
	if frameName == "" {
		frameName = "response"
	}
	frame := data.NewFrame(frameName)
	frame.Meta = &data.FrameMeta{
		ExecutedQueryString: "This feature is not available for this type of query yet",
		Custom: &CustomMeta{
			Query:                  query,
			Data:                   query.Data,
			ResponseCodeFromServer: 0,
			Error:                  "",
		},
	}
	return frame
}

func WrapMetaForInlineQuery(ctx context.Context, frame *data.Frame, err error, query models.Query) (*data.Frame, error) {
	if frame == nil {
		frame = data.NewFrame(query.RefID)
	}
	customMeta := &CustomMeta{Query: query, Data: query.Data, ResponseCodeFromServer: 0}
	if err != nil {
		customMeta.Error = err.Error()
	}
	frame.Meta = &data.FrameMeta{
		ExecutedQueryString: "This feature is not available for this type of query yet",
		Custom:              customMeta,
	}
	frame = ApplyLogMeta(frame, query)
	frame = ApplyTraceMeta(frame, query)
	return frame, err
}

func WrapMetaForRemoteQuery(ctx context.Context, frame *data.Frame, err error, query models.Query) (*data.Frame, error) {
	if frame == nil {
		frame = data.NewFrame(query.RefID)
	}
	meta := frame.Meta
	if meta == nil {
		customMeta := &CustomMeta{Query: query, Data: query.Data, ResponseCodeFromServer: 0}
		if err != nil {
			customMeta.Error = err.Error()
		}
		frame.Meta = &data.FrameMeta{Custom: customMeta}
	}
	frame = ApplyLogMeta(frame, query)
	frame = ApplyTraceMeta(frame, query)
	return frame, err
}

func ApplyLogMeta(frame *data.Frame, query models.Query) *data.Frame {
	if frame == nil {
		frame = data.NewFrame(query.RefID)
	}
	if frame.Meta == nil {
		frame.Meta = &data.FrameMeta{}
	}
	if query.Format == "logs" {
		doesTimeFieldExist := false
		doesBodyFieldExist := false
		for _, field := range frame.Fields {
			if field.Name == "timestamp" && (field.Type() == data.FieldTypeNullableTime || field.Type() == data.FieldTypeTime) {
				doesTimeFieldExist = true
			}
			if field.Name == "body" && (field.Type() == data.FieldTypeNullableString || field.Type() == data.FieldTypeString) {
				doesBodyFieldExist = true
			}
		}
		if doesTimeFieldExist && doesBodyFieldExist {
			frame.Meta.Type = data.FrameTypeLogLines
			frame.Meta.TypeVersion = data.FrameTypeVersion{0, 0}
		}
		frame.Meta.PreferredVisualization = data.VisTypeLogs
	}
	return frame
}

func ApplyTraceMeta(frame *data.Frame, query models.Query) *data.Frame {
	if frame == nil {
		frame = data.NewFrame(query.RefID)
	}
	if frame.Meta == nil {
		frame.Meta = &data.FrameMeta{}
	}
	if query.Format == "trace" {
		frame.Meta.PreferredVisualization = data.VisTypeTrace
	}
	return frame
}
