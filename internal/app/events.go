package app

const (
	eventInit                  = "latency:init"
	eventRecordStart           = "latency:record_start"
	eventRecordStartError      = "latency:record_start_error"
	eventRecordFilish          = "latency:record_filish"
	eventRecordStopError       = "latency:record_stop_error"
	eventTransformStart        = "latency:transform_start"
	eventTransformStartError   = "latency:transform_start_error"
	eventTransformFilish       = "latency:transform_filish"
	eventAnalyseStart          = "latency:analyse_start"
	eventAnalyseFilish         = "latency:analyse_filish"
	eventSetPointerLocationOn  = "latency:set_pointer_location_on"
	eventSetPointerLocationOff = "latency:set_pointer_location_off"
	eventUpdateAvaiable        = "latency:update_available"
	eventUpdateSuccess         = "latency:update_success"
	eventUpdateError           = "latency:update_error"
)
