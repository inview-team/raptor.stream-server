package connector

import (
	"fmt"
	"github.com/inview-team/raptor.stream-server/internal/config"
	"github.com/pion/webrtc/v3"
)

var peerConnectionConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

type Broadcaster struct {
	config *config.Settings
}

func New(conf *config.Settings) *Broadcaster {
	return &Broadcaster{
		config: conf,
	}
}

func (b *Broadcaster) CreateConnection(offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return nil, err
	}

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		return nil, err
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
	return &answer, nil
}
