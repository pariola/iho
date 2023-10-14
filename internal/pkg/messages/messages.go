package messages

type Error struct{ Message string }

func (Error) T() Type { return TypeError }

type Ping struct{}

func (Ping) T() Type { return TypePing }

type OpenTunnel struct{}

func (OpenTunnel) T() Type { return TypeOpenTunnel }

type TunnelOpened struct {
	Port string
}

func (TunnelOpened) T() Type { return TypeTunnelOpened }

type ConnectSession struct {
	Id string
}

func (ConnectSession) T() Type { return TypeConnectSession }

type RequestConnection struct {
	Id string
}

func (RequestConnection) T() Type { return TypeRequestConnection }
