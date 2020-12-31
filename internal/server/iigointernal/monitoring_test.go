package iigointernal

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

func TestAddToCache(t *testing.T) {
	cases := []struct {
		name        string
		roleID      shared.ClientID
		variables   []rules.VariableFieldName
		values      [][]float64
		expectedVal []shared.Accountability
	}{
		{
			name:      "Basic adding variables with corresponding values",
			roleID:    shared.ClientID(1),
			variables: []rules.VariableFieldName{rules.RuleSelected, rules.VoteCalled},
			values:    [][]float64{{1}, {1}},
			expectedVal: []shared.Accountability{
				{
					ClientID: shared.ClientID(1),
					Pairs: []rules.VariableValuePair{
						rules.MakeVariableValuePair(rules.RuleSelected, []float64{1}),
						rules.MakeVariableValuePair(rules.VoteCalled, []float64{1}),
					},
				},
			},
		},
		{
			name:        "Adding a variable and too many values",
			roleID:      shared.ClientID(1),
			variables:   []rules.VariableFieldName{rules.RuleSelected},
			values:      [][]float64{{1}, {1}},
			expectedVal: []shared.Accountability{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			monitor := &monitor{
				internalIIGOCache: []shared.Accountability{},
			}
			monitor.addToCache(tc.roleID, tc.variables, tc.values)
			res := monitor.internalIIGOCache
			if !reflect.DeepEqual(res, tc.expectedVal) {
				t.Errorf("Expected internalIIGOCache to be %v got %v", tc.expectedVal, res)
			}
		})
	}
}

func TestEvaluateCache(t *testing.T) {
	cases := []struct {
		name        string
		roleID      shared.ClientID
		iigoCache   []shared.Accountability
		expectedVal bool
	}{
		{
			name:   "Basic evaluation of compliant President",
			roleID: shared.ClientID(1),
			iigoCache: []shared.Accountability{
				{
					ClientID: shared.ClientID(1),
					Pairs: []rules.VariableValuePair{
						rules.MakeVariableValuePair(rules.RuleSelected, []float64{1}),
						rules.MakeVariableValuePair(rules.VoteCalled, []float64{1}),
					},
				},
			},
			expectedVal: true,
		},
		{
			name:   "Basic evaluation of non compliant Speaker",
			roleID: shared.ClientID(1),
			iigoCache: []shared.Accountability{
				{
					ClientID: shared.ClientID(1),
					Pairs: []rules.VariableValuePair{
						rules.MakeVariableValuePair(rules.RuleSelected, []float64{0}),
						rules.MakeVariableValuePair(rules.VoteCalled, []float64{1}),
					},
				},
			},
			expectedVal: false,
		},
		{
			name:        "Evaluating with empty cache",
			roleID:      shared.ClientID(1),
			iigoCache:   []shared.Accountability{},
			expectedVal: true,
		},
	}
	ruleStore := registerMonitoringTestRule()
	tempCache := rules.AvailableRules
	rules.AvailableRules = ruleStore
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			monitor := &monitor{
				internalIIGOCache: tc.iigoCache,
			}
			res := monitor.evaluateCache(tc.roleID, ruleStore)
			if !reflect.DeepEqual(res, tc.expectedVal) {
				t.Errorf("Expected evaluation of internalIIGOCache to be %v got %v", tc.expectedVal, res)
			}
		})
	}
	rules.AvailableRules = tempCache
}


func TestFindRoleToMonitor(t *testing.T) {
	cases := []struct{
		name 			string
		roleAccountable shared.ClientID
		expectedVal 	(shared.ClientID, shared.Role, error)
	}{
		{
			name: "Test Speaker to perform monitoring",
			roleAccountable: shared.ClientID(1),
			expectedVal: (shared.ClientID(2), shared.President, nil),
		},
		{
			name: "Test President to perform monitoring",
			roleAccountable: shared.ClientID(2),
			expectedVal: (shared.ClientID(3), shared.Judge, nil),
		},
		{
			name: "Test Judge to perform monitoring",
			roleAccountable: shared.ClientID(3),
			expectedVal: (shared.ClientID(1), shared.Speaker, nil),
		},
		{
			name: "Test non IIGO role trying to perform monitoring",
			roleAccountable: shared.ClientID(4),
			expectedVal: (shared.ClientID(-1), shared.Speaker, errors.Errorf("Monitoring by island that is not an IIGO Role")),
		},
	}
	monitor := &monitor{
		speakerID: 1,
		presidentID: 2,
		judgeID: 3,
	}
	for _, tc := range cases{
		t.Run(tc.name, func(t *testting.T){
			res:=monitor.findRoleToMonitor(tc.roleAccountable)
			if!reflect.DeepEqual(tc.expectedVal, res){
				t.Errorf("Expected role to monitor to be %v got %v", tc.expectedVal, res)
			}
		})
	}
}

func registerMonitoringTestRule() map[string]rules.RuleMatrix {

	rulesStore := map[string]rules.RuleMatrix{}

	name := "vote_called_rule"
	reqVar := []rules.VariableFieldName{
		rules.RuleSelected,
		rules.VoteCalled,
	}

	v := []float64{1, -1, 0}
	CoreMatrix := mat.NewDense(1, 3, v)
	aux := []float64{0}
	AuxiliaryVector := mat.NewVecDense(1, aux)

	rm := rules.RuleMatrix{RuleName: name, RequiredVariables: reqVar, ApplicableMatrix: *CoreMatrix, AuxiliaryVector: *AuxiliaryVector, Mutable: false}
	rulesStore[name] = rm
	return rulesStore
}
