package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type fakeServerHandle struct {
	PresidentID shared.ClientID
	JudgeID     shared.ClientID
	SpeakerID   shared.ClientID
	TermLengths map[shared.Role]uint
}

func (s fakeServerHandle) GetGameState() gamestate.ClientGameState {
	return gamestate.ClientGameState{
		SpeakerID:   s.SpeakerID,
		JudgeID:     s.JudgeID,
		PresidentID: s.PresidentID,
	}
}

func (s fakeServerHandle) GetGameConfig() config.ClientConfig {
	return config.ClientConfig{
		IIGOClientConfig: config.IIGOConfig{
			IIGOTermLengths: map[shared.Role]uint{},
		},
	}
}
