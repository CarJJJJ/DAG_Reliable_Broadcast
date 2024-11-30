package brachabroadcast

type InitialMessage struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	NodeID  int    `json:"node_id"`
}

type EchoMessage struct {
	Type           int            `json:"type"`
	InitialMessage InitialMessage `json:"initial_message"`
	NodeID         int            `json:"node_id"`
}

type ReadyMessage struct {
	Type           int         `json:"type"`
	InitialMessage EchoMessage `json:"echo_message"`
	NodeID         int         `json:"node_id"`
}
