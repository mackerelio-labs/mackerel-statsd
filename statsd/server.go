package statsd

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"

	"github.com/hatena/mackerelstatsd/parser"
)

type Server struct {
	Error func(err error)

	sync.Mutex
	c *net.UDPConn
}

type Handler func(addr string, m *parser.Metric) error

func (s *Server) ListenAndServe(ctx context.Context, addr string, handle Handler) error {
	if err := s.listen(addr); err != nil {
		return err
	}
	c := make(chan struct{})
	go func() {
		s.serve(handle)
		close(c)
	}()
	<-ctx.Done()
	s.closeConn()
	<-c
	return nil
}

func (s *Server) listen(addr string) error {
	udpAddr, err := parseAddr(addr)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	if s.c != nil {
		return errors.New("server is already running")
	}
	c, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.c = c
	return nil
}

func parseAddr(addr string) (*net.UDPAddr, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	return &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	}, nil
}

func (s *Server) serve(handle Handler) {
	handleErr := s.Error
	if handleErr != nil {
		handleErr = func(err error) {}
	}

	buf := make([]byte, 65535)
	var wg sync.WaitGroup
	conn := s.c
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			handleErr(err)
			continue
		}
		a, err := parser.Parse(buf[:n])
		if err != nil {
			handleErr(err)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, m := range a {
				if err := handle(addr.String(), m); err != nil {
					handleErr(err)
				}
			}
		}()
	}
	wg.Wait()
}

func (s *Server) closeConn() error {
	if s.c == nil {
		return nil
	}
	c := s.c
	s.c = nil
	return c.Close()
}
