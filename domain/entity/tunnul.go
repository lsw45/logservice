package entity

type UpdateFileReq struct {
	Remote   string
	Server   string
	Preserve bool
}

type TunnelUploadFileRes struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
	} `json:"data"`
}

type ShellTaskReq struct {
	Env           int           `json:"env"`
	Params        []ShellParams `json:"params"`
	Project       int           `json:"project"`
	Asynchronous  bool          `json:"asynchronous"`
	CorporationID string        `json:"corporation_id"`
}

type ShellTaskDeployResp struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data ShellData `json:"data"`
}

type ShellData struct {
	Env           int64         `json:"env"`
	Params        []ShellParams `json:"params"`
	Project       int64         `json:"project"`
	Asynchronous  bool          `json:"asynchronous"`
	CorporationID string        `json:"corporation_id"`
}

type ShellParams struct {
	Shell            string                 `json:"shell"`
	Server           string                 `json:"server"`
	ShellEnvironment map[string]interface{} `json:"environment"`
	Command          []string               `json:"command"`
}

type ShellTaskStateResp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		ID            int           `json:"id"`
		AddTime       string        `json:"add_time"`
		UpdateTime    string        `json:"update_time"`
		StartTime     interface{}   `json:"start_time"`
		FinishTime    interface{}   `json:"finish_time"`
		CorporationID string        `json:"corporation_id"`
		Project       int           `json:"project"`
		Env           string        `json:"env"`
		Params        []ShellParams `json:"params"`
		Asynchronous  bool          `json:"asynchronous"`
		Status        string        `json:"status"`
		Result        ShellResult   `json:"result"`
		UUID          string        `json:"uuid"`
		UsedTime      int           `json:"used_time"`
		Comment       interface{}   `json:"comment"`
		Creator       string        `json:"creator"`
	} `json:"data"`
}

type ShellResult struct {
	Detail []Detail `json:"detail"`
	Status string   `json:"status"`
}

type Data struct {
	ID     int                    `json:"id"`
	Status string                 `json:"status"`
	Result map[string]ShellResult `json:"result"`
	UUID   string                 `json:"uuid"`
}

type Detail struct {
	Exited  int    `json:"exited"`
	Stderr  string `json:"stderr"`
	Stdout  string `json:"stdout"`
	Command string `json:"command"`
}
