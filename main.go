package go_srv

import "github.com/xtulnx/go-srv/netkit"

type jNetKit byte

var NetKit jNetKit = 0

func (jNetKit) PickUnusedPort() (int, error) {
	return netkit.PickUnusedPort()
}
func (jNetKit) LocalIP() (string, error) {
	return netkit.LocalIP()
}
func (jNetKit) LocalMac() (string, error) {
	return netkit.LocalMac()
}
