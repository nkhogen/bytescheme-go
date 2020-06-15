package main

import (
	"bytescheme/common/log"
	"bytescheme/common/service"
	"bytescheme/common/util"
	"encoding/json"
	"time"

	gmodels "bytescheme/controller/generated/models"
)

// For generating jsons

func main() {
	pinID := int32(116)
	controller := &gmodels.Controller{
		Pins: []*gmodels.Pin{
			{
				ID:    &pinID,
				Name:  "Fountain",
				Mode:  "Output",
				Value: "Low",
			},
		},
	}
	data := map[string]interface{}{}
	err := util.Convert(controller, &data)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	event := &service.Event{
		ID:        "a9d70934-ea5a-4f9e-92bc-6ad49b184a8c",
		Time:      now,
		RecurMins: 1,
		Data:      data,
	}
	ba, _ := json.MarshalIndent(event, " ", " ")
	log.Infof("JSON:\n%s\n", string(ba))
}
