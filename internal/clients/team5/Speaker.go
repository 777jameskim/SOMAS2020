package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type speaker struct {
	*baseclient.BaseSpeaker
	c *client
}

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	c.Logf("Team 5 became speaker")
	return &c.team5Speaker
}

// Pay Judge based on the status of our own wealth
// If we are not doing verywell, pay Judge less so we have more in the CP to take from
func (s *speaker) PayJudge() shared.SpeakerReturnContent {
	JudgeSalaryRule, ok := rules.RulesInPlay["salary_cycle_judge"]
	var JudgeSalary shared.Resources = 0
	if ok {
		JudgeSalary = shared.Resources(JudgeSalaryRule.ApplicableMatrix.At(0, 1))
	}
	if s.c.wealth() == jeffBezos {
		return shared.SpeakerReturnContent{
			ContentType: shared.SpeakerJudgeSalary,
			JudgeSalary: JudgeSalary,
			ActionTaken: true,
		}
	} else if s.c.wealth() == middleClass {
		JudgeSalary = JudgeSalary * 0.8
	} else {
		JudgeSalary = JudgeSalary * 0.5
	}
	return shared.SpeakerReturnContent{
		ContentType: shared.SpeakerJudgeSalary,
		JudgeSalary: JudgeSalary,
		ActionTaken: true,
	}
}

// if the real winner is on our bad side, then we choose our best friend
func (s *speaker) DecideNextJudge(winner shared.ClientID) shared.ClientID {
	if s.c.opinions[winner].getScore() < 0 {
		ballot := s.c.GetVoteForElection(shared.President)
		for _, clientID := range ballot {
			if s.c.isClientAlive(clientID) {
				return clientID
			}
		}
	}
	return winner
}
