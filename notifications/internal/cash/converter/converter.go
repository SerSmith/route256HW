package converter

import (
	"route256/notifications/internal/cash/schema"
	"route256/notifications/internal/domain"
)

func NotificationMemD2S(domainStruct []domain.NotificationMem) []schema.StatusChangeMessage {

	out := make([]schema.StatusChangeMessage, 0, len(domainStruct))

	for _, ds := range domainStruct {
		scmSchema := schema.StatusChangeMessage{NewStatus: ds.ChangeStatus.NewStatus,
			OldStatus: ds.ChangeStatus.OldStatus,
			OrderID:   ds.ChangeStatus.OrderID,
			DT:        ds.DT,
			UserID:    ds.ChangeStatus.UserID}
		out = append(out, scmSchema)

	}

	return out
}

func NotificationMemS2D(schemaStruct []schema.StatusChangeMessage) []domain.NotificationMem {

	out := make([]domain.NotificationMem, 0, len(schemaStruct))

	for _, ss := range schemaStruct {

		ChangeStatus := domain.StatusChangeMessage{NewStatus: ss.NewStatus,
			OldStatus: ss.OldStatus,
			OrderID:   ss.OrderID,
			UserID:    ss.UserID}

		dnm := domain.NotificationMem{ChangeStatus: &ChangeStatus,
			DT: ss.DT}

		out = append(out, dnm)

	}

	return out
}

func NotificationHistoryRequestD2S(domainStruct domain.NotificationHistoryRequest) schema.NotificationHistoryRequest {
	return schema.NotificationHistoryRequest{
		UserID:   domainStruct.UserID,
		DateFrom: domainStruct.DateFrom,
		DateTo:   domainStruct.DateTo}
}
