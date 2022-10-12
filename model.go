package api

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type NotebookLsInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Sort   uint64 `json:"sort"`
	Closed bool   `json:"closed"`
}

type OperationInfo struct {
	Action     string `json:"action"`
	Data       string `json:"data"`
	Id         string `json:"id"`
	ParentId   string `json:"parentID"`
	PreviousId string `json:"previousID"`
	RetData    string `json:"retData"`
}

type responseBlockOperations struct {
	Response
	Data []struct {
		DoOperations   []OperationInfo `json:"doOperations"`
		UndoOperations interface{}     `json:"undoOperations"`
	}
}
