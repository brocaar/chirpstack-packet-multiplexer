package multiplexer

import (
	"net"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func init() {
	log.SetLevel(log.ErrorLevel)
}

type MultiplexerTestSuite struct {
	suite.Suite

	mp     *Multiplexer
	mpAddr *net.UDPAddr

	pfConn *net.UDPConn

	b1Conn *net.UDPConn
	b2Conn *net.UDPConn
	b1Addr *net.UDPAddr
	b2Addr *net.UDPAddr
}

func (ts *MultiplexerTestSuite) SetupTest() {
	assert := require.New(ts.T())
	var err error

	ts.mpAddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:1800")
	assert.NoError(err)

	ts.b1Addr, err = net.ResolveUDPAddr("udp", "127.0.0.1:1801")
	assert.NoError(err)
	ts.b2Addr, err = net.ResolveUDPAddr("udp", "127.0.0.1:1802")
	assert.NoError(err)

	// backends
	ts.b1Conn, err = net.ListenUDP("udp", ts.b1Addr)
	assert.NoError(err)
	ts.b2Conn, err = net.ListenUDP("udp", ts.b2Addr)
	assert.NoError(err)

	// setup packet-forwarder
	ts.pfConn, err = net.DialUDP("udp", nil, ts.mpAddr)
	assert.NoError(err)

	ts.mp, err = New(Config{
		Bind: "127.0.0.1:1800",
		Backends: []BackendConfig{
			{
				Host: "127.0.0.1:1801",
				GatewayIDs: []string{
					"0101010101010101",
					"0202020202020202",
				},
			},
			{
				Host: "127.0.0.1:1802",
				GatewayIDs: []string{
					"0202020202020202",
				},
			},
		},
	})
	assert.NoError(err)
}

func (ts *MultiplexerTestSuite) TearDownTest() {
	assert := require.New(ts.T())
	assert.NoError(ts.mp.Close())
	assert.NoError(ts.pfConn.Close())
	assert.NoError(ts.b1Conn.Close())
	assert.NoError(ts.b2Conn.Close())
}

func (ts *MultiplexerTestSuite) TestUplink() {
	tests := []struct {
		Name                    string
		PacketForwarderSent     []byte
		PacketForwarderReceived []byte
		Backend1Received        []byte
		Backend2Received        []byte
	}{
		{
			Name:                    "PUSH_DATA gw 1",
			PacketForwarderSent:     []byte{0x02, 0x01, 0x02, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x03, 0x04, 0x05},
			PacketForwarderReceived: []byte{0x02, 0x01, 0x02, 0x01},
			Backend1Received:        []byte{0x02, 0x01, 0x02, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x03, 0x04, 0x05},
		},
		{
			Name:                    "PUSH_DATA gw 2",
			PacketForwarderSent:     []byte{0x02, 0x01, 0x02, 0x00, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x04, 0x05},
			PacketForwarderReceived: []byte{0x02, 0x01, 0x02, 0x01},
			Backend1Received:        []byte{0x02, 0x01, 0x02, 0x00, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x04, 0x05},
			Backend2Received:        []byte{0x02, 0x01, 0x02, 0x00, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03, 0x04, 0x05},
		},
		{
			Name:                    "PULL_DATA gw 1",
			PacketForwarderSent:     []byte{0x02, 0x01, 0x02, 0x02, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			PacketForwarderReceived: []byte{0x02, 0x01, 0x02, 0x04},
			Backend1Received:        []byte{0x02, 0x01, 0x02, 0x02, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
		},
		{
			Name:                    "PULL_DATA gw 2",
			PacketForwarderSent:     []byte{0x02, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			PacketForwarderReceived: []byte{0x02, 0x01, 0x02, 0x04},
			Backend1Received:        []byte{0x02, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			Backend2Received:        []byte{0x02, 0x01, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
		},
		{
			Name:                "TX_ACK gw 1",
			PacketForwarderSent: []byte{0x02, 0x01, 0x02, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			Backend1Received:    []byte{0x02, 0x01, 0x02, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
		},
		{
			Name:                "TX_ACK gw 2",
			PacketForwarderSent: []byte{0x02, 0x01, 0x02, 0x05, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			Backend1Received:    []byte{0x02, 0x01, 0x02, 0x05, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			Backend2Received:    []byte{0x02, 0x01, 0x02, 0x05, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
		},
	}

	for _, tst := range tests {
		ts.T().Run(tst.Name, func(t *testing.T) {
			assert := require.New(t)

			_, err := ts.pfConn.Write(tst.PacketForwarderSent)
			assert.NoError(err)

			ts.pfConn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
			b := make([]byte, 65507)
			i, _ := ts.pfConn.Read(b)
			assert.Equal(len(tst.PacketForwarderReceived), i)
			if len(tst.PacketForwarderReceived) > 0 {
				assert.EqualValues(tst.PacketForwarderReceived, b[:i])
			}

			ts.b1Conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
			i, _ = ts.b1Conn.Read(b)
			assert.Equal(len(tst.Backend1Received), i)
			if len(tst.Backend1Received) > 0 {
				assert.EqualValues(tst.Backend1Received, b[:i])
			}

			ts.b2Conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
			i, _ = ts.b2Conn.Read(b)
			assert.Equal(len(tst.Backend2Received), i)
			if len(tst.Backend2Received) > 0 {
				assert.EqualValues(tst.Backend2Received, b[:i])
			}
		})
	}
}

func (ts *MultiplexerTestSuite) TestDownlinkBackend1() {
	tests := []struct {
		Name                    string
		BackendSent             []byte
		BackendReceived         []byte
		PacketForwarderReceived []byte
	}{
		{
			Name:                    "PULL_RESP",
			BackendSent:             []byte{0x02, 0x01, 0x02, 0x03, 0x01, 0x02, 0x03},
			PacketForwarderReceived: []byte{0x02, 0x01, 0x02, 0x03, 0x01, 0x02, 0x03},
		},
	}

	ts.mp.gateways["0101010101010101"] = ts.pfConn.LocalAddr().(*net.UDPAddr)

	for _, tst := range tests {
		ts.T().Run(tst.Name, func(t *testing.T) {
			assert := require.New(t)
			addr := ts.mp.backends["127.0.0.1:1801"]["0101010101010101"].LocalAddr().(*net.UDPAddr)
			_, err := ts.b1Conn.WriteToUDP(tst.BackendSent, addr)
			assert.NoError(err)

			ts.b1Conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
			b := make([]byte, 65507)
			i, _ := ts.b1Conn.Read(b)
			assert.Equal(len(tst.BackendReceived), i)
			if i > 0 {
				assert.EqualValues(tst.BackendReceived, b[:i])
			}

			ts.pfConn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
			i, _ = ts.pfConn.Read(b)
			assert.Equal(len(tst.PacketForwarderReceived), i)
			if i > 0 {
				assert.EqualValues(tst.PacketForwarderReceived, b[:i])
			}
		})
	}
}

func TestMultiplexer(t *testing.T) {
	suite.Run(t, new(MultiplexerTestSuite))
}
