package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/inview-team/raptor.stream-server/internal/app/connector"
	"github.com/inview-team/raptor.stream-server/internal/config"
	"github.com/inview-team/raptor.stream-server/internal/logger"
	"github.com/pion/webrtc/v3"
	"io/ioutil"
	"net/http"
)

type Server struct {
	http *http.Server
	con  *connector.Broadcaster
}

func New(addr string, con *connector.Broadcaster) *Server {
	var server = &Server{
		http: &http.Server{
			Addr: addr,
		},
		con: con,
	}
	server.http.Handler = server.setupRouter()
	return server
}

func (s *Server) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/candidate", s.addNewICECandidate)
	r.POST("/offer", s.createNewWebRTCConnection)
	return r
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.http.Close()
}

func (s *Server) addNewICECandidate(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Failed to read body request", "msg_err": err.Error()})
		return
	}
	taskId := c.Request.Header.Get("uuid")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Task ID didn't provide"})
		return
	}
	connection := s.con.WC_storage[taskId]
	peerConnection := connection.PeerConnection["worker"]
	peerConnection, err = webrtc.NewPeerConnection(config.PeerConnectionConfig)
	if err != nil {
		logger.Critical.Panic(err)
	}
	defer func() {
		if err := peerConnection.Close(); err != nil {
			fmt.Printf("cannot close peerConnection: %v\n", err)
		}
	}()

	if candidateErr := peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: string(bodyBytes)}); candidateErr != nil {
		panic(candidateErr)
	}
}

func (s *Server) createNewWebRTCConnection(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)

	taskId := c.Request.Header.Get("uuid")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Task ID didn't provide"})
		return
	}

	connection := s.con.WC_storage[taskId]
	peerConnection := connection.PeerConnection["worker"]

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Failed to read body request", "msg_err": err.Error()})
		return
	}

	var offer webrtc.SessionDescription
	if err = json.Unmarshal(bodyBytes, &offer); err != nil {
		logger.Critical.Panic(err.Error())
	}
	logger.Info.Printf(offer.SDP)

	answer, err := s.con.CreateConnection(peerConnection, offer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Failed to create connection", "msg_err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, answer)
}
