package spotdata

type InstanceType struct {
	EMR   bool    `json:"emr"`
	Cores int64   `json:"cores"`
	RamGB float32 `json:"ram_gb"`
}

type Range struct {
	Index int64  `json:"index"`
	Label string `json:"label"`
	Dots  int64  `json:"dots"`
	Max   int64  `json:"max"`
}

type InstanceTypeStats struct {
	Savings           int64 `json:"s"`
	InterruptionLevel int64 `json:"r"`
}

type OS map[string]InstanceTypeStats

type Region map[string]OS

type SpotDataFile struct {
	InstanceTypes map[string]InstanceType `json:"instance_types"`
	Ranges        []Range                 `json:"ranges"`
	SpotAdvisor   map[string]Region       `json:"spot_advisor"`
	GlobalRate    string                  `json:"global_rate"`
}
