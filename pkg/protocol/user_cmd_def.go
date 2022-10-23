package protocol

import "trcell/pkg/evhub"

const (
	CellCmdClassWebsocket = evhub.HubCmdUserBase + 1
	CellCmdClassTcpsocket = evhub.HubCmdUserBase + 2
)

// Websock
const (
	CmdTypeWebsocketConnect = 1
	CmdTypeWebsocketClosed  = 2
	CmdTypeWebsocketMessage = 3
)

// Tcpsock
const (
	CmdTypeTcpsocketConnect = 1
	CmdTypeTcpsocketClosed  = 2
	CmdTypeTcpsocketMessage = 3
)
