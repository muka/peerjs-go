package server

// IClient client interface
type IClient interface {
	GetID() string
	GetToken() string
	GetSocket() *WebSocket
	SetSocket(socket *WebSocket)
	GetLastPing() int64
	SetLastPing(lastPing int64)
	Send(data []byte) error
}

// Client implementation
type Client struct {
	id       string
	token    string
	socket   *WebSocket
	lastPing int64
}

//NewClient initialize a new client
func NewClient(id string, token string) *Client {
	c := new(Client)
	c.id = id
	c.token = token
	return c
}

//GetID return client id
func (c *Client) GetID() string {
	return c.id
}

//GetToken return client token
func (c *Client) GetToken() string {
	return c.token
}

//GetSocket return the web socket server
func (c *Client) GetSocket() *WebSocket {
	return c.socket
}

//SetSocket set the web socket handler
func (c *Client) SetSocket(socket *WebSocket) {
	c.socket = socket
}

// GetLastPing return the last ping timestamp
func (c *Client) GetLastPing() int64 {
	return c.lastPing
}

//SetLastPing set last ping timestamp
func (c *Client) SetLastPing(lastPing int64) {
	c.lastPing = lastPing
}

//Send send data
func (c *Client) Send(data []byte) error {
	return c.socket.Send(data)
}
