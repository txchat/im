package acc

type AuthReply struct {
	Result  int                    `json:"result"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type AuthErrorDataReconnectNotAllowed struct {
	Code    int `json:"code"`
	Message struct {
		Datetime   int64  `json:"datetime"`
		Device     int    `json:"device"`
		DeviceName string `json:"deviceName"`
		Uuid       string `json:"uuid"`
	} `json:"message"`
	Service string `json:"service"`
}

type AuthSuccessData struct {
	Address string `json:"address"`
}
