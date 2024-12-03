package brachabroadcast

const (
	initial_type = 0
	echo_type    = 1
	ready_type   = 2
)

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
	Type           int            `json:"type"`
	InitialMessage InitialMessage `json:"initial_message"`
	NodeID         int            `json:"node_id"`
}
