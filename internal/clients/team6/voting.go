package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// VoteForRule returns the client's vote in favour of or against a rule.
func (c *client) VoteForRule(ruleMatrix rules.RuleMatrix) shared.RuleVoteType {
	// TODO implement decision on voting that considers the rule
	return shared.Abstain
}

// VoteForElection returns the client's Borda vote for the role to be elected.
func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	//Sort candidates according to friendship level as a preference list
	//Every Island votes for itself first
	friendship[id] = 0
	for _, candidateID := range candidateList {
		friendship[id] += friendship[candidateID]
	}
	for i := 0; i < len(candidateList); i++ {
		for j := i; j < len(candidateList); j++ {
			if friendship[candidateList[j]] > friendship[candidateList[i]] {
				candidateList[i], candidateList[j] = candidateList[j], candidateList[i]
			}
		}
	}

	preferenceList := candidateList

	return preferenceList
}
>>>>>>> main
