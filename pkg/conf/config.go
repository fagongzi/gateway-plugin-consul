package conf

// Conf config
type Conf struct {
	Addr string `json:"addr"`

	ConsulAddr string `json:"consulAddr"`

	AuthEnabled  bool   `json:"authEnabled,omitempty"`
	AuthUserName string `json:"authUserName,omitempty"`
	AuthPassword string `json:"authPassword,omitempty"`

	Token                  string `json:"token,omitempty"`
	Timeout                int    `json:"timeout,omitempty"`
	HeartbeatsBeforeRemove int    `json:"heartbeatsBeforeRemove,omitempty"`
}
