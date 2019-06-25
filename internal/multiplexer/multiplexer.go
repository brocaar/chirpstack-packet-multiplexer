package multiplexer

import (
	"encoding/base64"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type udpPacket struct {
	addr *net.UDPAddr
	data []byte
}

// Multiplexer forwards packet-forwarder UDP data to multiple backends.
type Multiplexer struct {
	sync.RWMutex
	wg sync.WaitGroup

	conn   *net.UDPConn
	config Config
	closed bool

	backends map[string]map[string]*net.UDPConn // [backendHost][gatewayID]UDPConn
	gateways map[string]*net.UDPAddr            // [gatewayID]UDPAddr
}

// New creates a new multiplexer.
func New(c Config) (*Multiplexer, error) {
	m := Multiplexer{
		backends: make(map[string]map[string]*net.UDPConn),
		gateways: make(map[string]*net.UDPAddr),
	}

	addr, err := net.ResolveUDPAddr("udp", c.Bind)
	if err != nil {
		return nil, errors.Wrap(err, "resolve udp addr error")
	}

	log.WithField("addr", addr).Info("starting listener")
	m.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "listen udp error")
	}

	for _, backend := range c.Backends {
		addr, err := net.ResolveUDPAddr("udp", backend.Host)
		if err != nil {
			return nil, errors.Wrap(err, "resolve udp addr error")
		}

		for _, gatewayID := range backend.GatewayIDs {
			gatewayID = strings.ToLower(gatewayID)

			log.WithFields(log.Fields{
				"gateway_id":  gatewayID,
				"host":        backend.Host,
				"uplink_only": backend.UplinkOnly,
			}).Info("dial udp")
			conn, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				return nil, errors.Wrap(err, "dial udp error")
			}

			if _, ok := m.backends[backend.Host]; !ok {
				m.backends[backend.Host] = make(map[string]*net.UDPConn)
			}

			m.backends[backend.Host][gatewayID] = conn

			go func(backend, gatewayID string, conn *net.UDPConn) {
				m.wg.Add(1)
				err := m.readDownlinkPackets(backend, gatewayID, conn)
				if !m.isClosed() {
					log.WithError(err).Error("read udp packets error")
				}
				m.wg.Done()
			}(backend.Host, gatewayID, conn)
		}
	}

	go func() {
		m.wg.Add(1)
		err := m.readUplinkPackets()
		if !m.isClosed() {
			log.WithError(err).Error("read udp packets error")
		}
		m.wg.Done()
	}()

	return &m, nil
}

// Close closes the multiplexer.
func (m *Multiplexer) Close() error {
	m.Lock()
	m.closed = true

	log.Info("closing listener")
	if err := m.conn.Close(); err != nil {
		return errors.Wrap(err, "close udp listener error")
	}

	log.Info("closing backend connections")
	for _, gws := range m.backends {
		for _, conn := range gws {
			if err := conn.Close(); err != nil {
				return errors.Wrap(err, "close udp connection error")
			}
		}
	}

	m.Unlock()
	m.wg.Wait()
	return nil
}

func (m *Multiplexer) isClosed() bool {
	m.RLock()
	defer m.RUnlock()
	return m.closed
}

func (m *Multiplexer) setGateway(gatewayID string, addr *net.UDPAddr) error {
	m.Lock()
	defer m.Unlock()
	m.gateways[gatewayID] = addr
	return nil
}

func (m *Multiplexer) getGateway(gatewayID string) (*net.UDPAddr, error) {
	m.RLock()
	defer m.RUnlock()

	addr, ok := m.gateways[gatewayID]
	if !ok {
		return nil, errors.New("gateway does not exist")
	}
	return addr, nil
}

func (m *Multiplexer) readUplinkPackets() error {
	buf := make([]byte, 65507) // max udp data size
	for {
		i, addr, err := m.conn.ReadFromUDP(buf)
		if err != nil {
			if m.isClosed() {
				return nil
			}

			log.WithError(err).Error("read from udp error")
			continue
		}

		data := make([]byte, i)
		copy(data, buf[:i])
		up := udpPacket{data: data, addr: addr}

		// handle packet async
		go func(up udpPacket) {
			if err := m.handleUplinkPacket(up); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"data_base64": base64.StdEncoding.EncodeToString(up.data),
					"addr":        up.addr,
				}).Error("could not handle packet")
			}
		}(up)
	}
}

func (m *Multiplexer) readDownlinkPackets(backend, gatewayID string, conn *net.UDPConn) error {
	buf := make([]byte, 65507) // max udp data size
	for {
		i, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if m.isClosed() {
				return nil
			}

			log.WithError(err).Error("read from udp error")
			continue
		}

		data := make([]byte, i)
		copy(data, buf[:i])
		up := udpPacket{data: data, addr: addr}

		// handle packet async
		go func(up udpPacket) {
			if err := m.handleDownlinkPacket(backend, gatewayID, up); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"data_base64": base64.StdEncoding.EncodeToString(up.data),
					"addr":        up.addr,
				}).Error("could not handle packet")
			}
		}(up)
	}
}

func (m *Multiplexer) handleUplinkPacket(up udpPacket) error {
	pt, err := GetPacketType(up.data)
	if err != nil {
		return errors.Wrap(err, "get packet-type error")
	}

	gatewayID, err := GetGatewayID(up.data)
	if err != nil {
		return errors.Wrap(err, "get gateway id error")
	}

	log.WithFields(log.Fields{
		"packet_type": pt,
		"addr":        up.addr,
		"gateway_id":  gatewayID,
	}).Info("packet received from gateway")

	switch pt {
	case PushData:
		return m.handlePushData(gatewayID, up)
	case PullData:
		if err := m.setGateway(gatewayID, up.addr); err != nil {
			return errors.Wrap(err, "set gateway error")
		}
		return m.handlePullData(gatewayID, up)
	case TXACK:
		return m.forwardUplinkPacket(gatewayID, up)
	}

	return nil
}

func (m *Multiplexer) handleDownlinkPacket(backend, gatewayID string, up udpPacket) error {
	pt, err := GetPacketType(up.data)
	if err != nil {
		return errors.Wrap(err, "get packet-type error")
	}

	log.WithFields(log.Fields{
		"packet_type": pt,
		"gateway_id":  gatewayID,
		"host":        backend,
	}).Info("packet received from backend")

	switch pt {
	case PullResp:
		return m.forwardPullResp(backend, gatewayID, up)
	}

	return nil
}

func (m *Multiplexer) handlePushData(gatewayID string, up udpPacket) error {
	if len(up.data) < 12 {
		return errors.New("expected at least 12 bytes of data")
	}

	// respond with PushACK
	log.WithFields(log.Fields{
		"addr":        up.addr,
		"packet_type": PushACK,
		"gateway_id":  gatewayID,
	}).Info("sending packet to gateway")
	b := make([]byte, 4)
	copy(b[:3], up.data[:3])
	b[3] = byte(PushACK)
	if _, err := m.conn.WriteToUDP(b, up.addr); err != nil {
		return errors.Wrap(err, "write to udp error")
	}

	return m.forwardUplinkPacket(gatewayID, up)
}

func (m *Multiplexer) handlePullData(gatewayID string, up udpPacket) error {
	if len(up.data) < 12 {
		return errors.New("expected at least 12 bytes of data")
	}

	// respond with PullACK
	log.WithFields(log.Fields{
		"addr":        up.addr,
		"packet_type": PullACK,
		"gateway_id":  gatewayID,
	}).Info("sending packet to gateway")
	b := make([]byte, 4)
	copy(b[:3], up.data[:3])
	b[3] = byte(PullACK)
	if _, err := m.conn.WriteToUDP(b, up.addr); err != nil {
		return errors.Wrap(err, "write to udp error")
	}

	return m.forwardUplinkPacket(gatewayID, up)
}

func (m *Multiplexer) forwardUplinkPacket(gatewayID string, up udpPacket) error {
	for host, gwIDs := range m.backends {
		for gwID, conn := range gwIDs {
			if gwID == gatewayID {
				pt, err := GetPacketType(up.data)
				if err != nil {
					return errors.Wrap(err, "get packet-type error")
				}
				log.WithFields(log.Fields{
					"from":        up.addr,
					"to":          host,
					"gateway_id":  gatewayID,
					"packet_type": pt,
				}).Info("forwarding packet to backend")
				if _, err := conn.Write(up.data); err != nil {
					log.WithError(err).WithFields(log.Fields{
						"host":       host,
						"gateway_id": gwID,
					}).Error("udp write error")
				}
			}
		}
	}

	return nil
}

func (m *Multiplexer) forwardPullResp(backend, gatewayID string, up udpPacket) error {
	addr, err := m.getGateway(gatewayID)
	if err != nil {
		return errors.Wrap(err, "get gateway error")
	}

	if m.backendIsUplinkOnly(backend) {
		log.WithFields(log.Fields{
			"packet_type": PullResp,
			"gateway_id":  gatewayID,
			"host":        backend,
		}).Info("ignoring downlink packet, backend is uplink only")
		return nil
	}

	log.WithFields(log.Fields{
		"from":        backend,
		"to":          addr,
		"packet_type": PullResp,
		"gateway_id":  gatewayID,
	}).Info("forwarding packet to gateway")
	if _, err := m.conn.WriteToUDP(up.data, addr); err != nil {
		return errors.Wrap(err, "write to udp error")
	}

	return nil
}

func (m *Multiplexer) backendIsUplinkOnly(backend string) bool {
	for _, be := range m.config.Backends {
		if be.Host == backend {
			return be.UplinkOnly
		}
	}

	return true
}
