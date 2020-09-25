package main

type dataSource struct {
	ID            string  `json:"id"`
	Value         float64 `json:"init_value"`
	MaxChangeStep int     `json:"max_change_step"`
}

// to save into storage
type avgData struct {
	ID     string  `json:"id"`
	Value  float64 `json:"value"`
	Lenght int     `json:"-"`
}

type generator struct {
	TimeoutS    int          `json:"timeout_s"`
	SendPeriodS int          `json:"send_period_s"`
	DataSources []dataSource `json:"data_sources"`
}

type agregator struct {
	SubIds          []string   `json:"sub_ids"`
	AgregatePeriodS int        `json:"agregate_period_s"`
	Event           chan event // notify about events
	Start           chan event // notify about start
}

type config struct {
	Generators []generator `json:"generators"`
	Agregators []agregator `json:"agregators"`
	Queue      struct {
		Size int `json:"size"`
	} `json:"queue"`
	StorageType int `json:"storage_type"`
}
