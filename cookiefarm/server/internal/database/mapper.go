package database

import (
	"protocols"
	"time"
)

func MapFromFlagToDBParams(params Flag) AddFlagParams {
	return AddFlagParams{
		FlagCode:     params.FlagCode,
		ServiceName:  params.ServiceName,
		TeamID:       params.TeamID,
		Username:     params.Username,
		ExploitName:  params.ExploitName,
		Status:       params.Status,
		Msg:          params.Msg,
		ResponseTime: params.ResponseTime,
		SubmitTime:   params.SubmitTime,
		PortService:  params.PortService,
	}
}

func MapFromResponseProtocolToParamsToUpdate(params protocols.ResponseProtocol) UpdateFlagStatusByCodeParams {
	return UpdateFlagStatusByCodeParams{
		FlagCode:     params.Flag,
		Status:       params.Status,
		Msg:          params.Msg,
		ResponseTime: uint64(time.Now().Unix()),
	}
}

func MapFromGetFilteredFlagsRowToFlag(row []GetFilteredFlagsRow) []Flag {
	var flags []Flag
	for _, r := range row {
		flags = append(flags, Flag{
			ID:           r.ID,
			FlagCode:     r.FlagCode,
			ServiceName:  r.ServiceName,
			PortService:  uint16(r.PortService),
			TeamID:       r.TeamID,
			Username:     r.Username,
			ExploitName:  r.ExploitName,
			Status:       r.Status,
			Msg:          r.Msg,
			SubmitTime:   uint64(r.SubmitTime.Int64),
			ResponseTime: uint64(r.ResponseTime.Int64),
		})
	}

	if len(flags) == 0 {
		return []Flag{}
	}

	return flags
}
