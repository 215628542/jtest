package testmodel

type SysConfigInfo struct {
	Variable string `json:"variable,omitempty"`
	Value    string `json:"value,omitempty"`
	SetTime  string `json:"set_time,omitempty"`
	SetBy    string `json:"set_by,omitempty"`
}

func (this SysConfigInfo) TableName() string {
	return "sys_config"
}
