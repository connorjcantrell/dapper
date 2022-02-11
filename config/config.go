package config

type ApplicationDetails struct {
	Name              string            `json:"name"`
	AppId             uint              `json:"application_id"`
	Block             uint              `json:"block"`
	Creator           string            `json:"creator"`
	Revision          uint              `json:"revision"`
	Deleted           bool              `json:"deleted"`
	GlobalStateSchema GlobalStateSchema `json:"global_state_schema"`
	LocalStateSchema  LocalStateSchema  `json:"local_state_schema"`
}

type GlobalStateSchema struct {
	NumByteSlice uint `json:"num_byte_slice"`
	NumUint      uint `json:"num_uint"`
}

type LocalStateSchema struct {
	NumByteSlice uint `json:"num_byte_slice"`
	NumUint      uint `json:"num_uint"`
}

func newApplicationDetails(
	name string,
	globalNumByteSlice uint,
	globalNumUint uint,
	localNumByteSlice uint,
	localNumUint uint,
) ApplicationDetails {
	return ApplicationDetails{
		Name: name,
		GlobalStateSchema: GlobalStateSchema{
			NumByteSlice: globalNumByteSlice,
			NumUint:      globalNumUint,
		},
		LocalStateSchema: LocalStateSchema{
			NumByteSlice: localNumByteSlice,
			NumUint:      localNumUint,
		},
	}
}
