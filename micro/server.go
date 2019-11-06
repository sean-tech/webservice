package micro

import (
	"log"
	"net"
	"net/rpc"
	"runtime"
)

type Server struct {
	Plugins IPluginContainer
}

func NewServer() *Server {
	return &Server{
		Plugins: new(pluginContainer),
	}
}

func (s *Server) RegisterName(nameMap map[string]interface{}) {
	var err error
	for name, instance := range nameMap {
		err = rpc.RegisterName(name, instance)
		if err != nil {
			log.Fatal("Rpc RegisterName error:", err)
		}
	}
}

func (s *Server) Serve(network, address string) {
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		conn, accept := s.Plugins.DoConnAccept(conn)
		if !accept {
			s.Plugins.DoConnClose(conn)
			conn.Close()
		} else  {
			go s.ServeConn(conn)
		}
	}
}

func (s *Server) ServeConn(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			//log.Fatalf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
		}
		//s.mu.Lock()
		//delete(s.activeConn, conn)
		//s.mu.Unlock()
		//conn.Close()
		//
		//if s.Plugins == nil {
		//	s.Plugins = &pluginContainer{}
		//}
		//
		//s.Plugins.DoPostConnClose(conn)
	}()

	rpc.ServeConn(conn)
}
