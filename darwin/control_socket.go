package darwin

/*
# include "control_socket.h"
*/
import "C"
import (
	"log"
	"os/exec"
)

// ControlSocket is for control kext
type ControlSocket struct {
	sock int
}

// NewControlSocket creates an instance of ControlSocket
func NewControlSocket() *ControlSocket {

	cmd := exec.Command("bash", "darwin/ensureload.sh")
	err := cmd.Run()
	if err != nil {
		log.Println("cmd.Run", err)
		return nil
	}
	sock := int(C.Connect())
	if sock == -1 {
		return nil
	}
	return &ControlSocket{sock: sock}
}

// StartRedirect starts redirect traffic to port
func (cs *ControlSocket) StartRedirect(port uint16) bool {
	return bool(C.StartRedirect(C.int(cs.sock), C.ushort(port)))
}

// Close stops redirect
func (cs *ControlSocket) Close() bool {
	return bool(C.StopClose(C.int(cs.sock)))
}
