package micro

import (
	"fmt"
	"github.com/juju/ratelimit"
	"net"
	"sync"
	"testing"
	"time"
)

type HelloService struct {}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}


func TestServer(t *testing.T) {
	nameMap := make(map[string]interface{})
	nameMap["HelloService"] = new(HelloService)

	s := NewServer()
	s.RegisterName(nameMap)
	s.Plugins.Add(new(ServerRateLimitPlugin))
	go func() {
	s.Serve("tcp", ":1234")
	}()

	select {
	}
}


var (
	tokenBucketOnce sync.Once
	tokenBucket		*ratelimit.Bucket
)

func getTokenBucket() *ratelimit.Bucket {
	tokenBucketOnce.Do(func() {
		tokenBucket	= ratelimit.NewBucket(200 * time.Millisecond, 5)
	})
	return tokenBucket
}

type ServerRateLimitPlugin struct{}
func (this *ServerRateLimitPlugin) HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	if getTokenBucket().TakeAvailable(1) > 0 {
		fmt.Println("HandleConnAccept")
		return conn, true
	}
	return conn, false
	//errors.New("服务访问流量已满，请稍后重试\nService access traffic is full, please try again later")
}

func (this *ServerRateLimitPlugin) HandleConnClose(conn net.Conn) bool {
	fmt.Println("HandleConnClose")
	return true
}

