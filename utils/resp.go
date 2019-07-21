package utils

type SuccResp struct {
	Data    interface{} `json: "data"`
	Success bool        `json: "success"`
}

type FailResp struct {
	Error   interface{} `json: "error"`
	Success bool        `json: "success"`
}
