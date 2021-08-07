package server

import (
	"encoding/json"
	"github.com/inview-team/raptor.stream-server/internal/app/connector"
	"github.com/inview-team/raptor.stream-server/internal/logger"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
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
	r.POST("/offer", s.createNewWebRTCConnection)
	return r
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.http.Close()
}

func (s *Server) createNewWebRTCConnection(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to read request body": err.Error()})
		return
	}

	var offer webrtc.SessionDescription
	if err = json.Unmarshal(bodyBytes, &offer); err != nil {
		logger.Critical.Panic(err.Error())
	}
	logger.Info.Printf(string(offer.SDP))
	answer, err := s.con.CreateConnection(offer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to create task": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *answer)
}
