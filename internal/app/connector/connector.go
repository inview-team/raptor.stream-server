package connector

import (
	"github.com/inview-team/raptor.stream-server/internal/config"
	"github.com/inview-team/raptor.stream-server/internal/logger"
	"github.com/pion/webrtc/v3"
)

type WebRTCStorage struct {
	PendingCandidates []*webrtc.ICECandidate
	PeerConnection    map[string]*webrtc.PeerConnection
}

type Broadcaster struct {
	config     *config.Settings
	WC_storage map[string]*WebRTCStorage
}

func New(conf *config.Settings) *Broadcaster {
	return &Broadcaster{
		config:     conf,
		WC_storage: make(map[string]*WebRTCStorage),
	}
}

func (b *Broadcaster) CreateConnection(peerConnection *webrtc.PeerConnection, offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	if _, err := peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return nil, err
	}

	err := peerConnection.SetRemoteDescription(offer)
	if err != nil {
		return nil, err
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	go func() {
		peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
			logger.Info.Printf("Peer Connection State has changed: %s\n", s.String())

			if s == webrtc.PeerConnectionStateFailed {
				// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
				// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
				// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
				logger.Error.Printf("Peer Connection has gone to failed exiting")
			}
		})
	}()
	return &answer, nil
}
