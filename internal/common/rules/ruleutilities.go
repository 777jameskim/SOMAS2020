package rules

import "github.com/pkg/errors"

// PickUpRulesByVariable returns a list of rule_id's which are affected by certain variables.
func PickUpRulesByVariable(variableName VariableFieldName, ruleStore map[string]RuleMatrix) ([]string, bool) {
	var Rules []string
	if _, ok := VariableMap[variableName]; ok {
		for k, v := range ruleStore {
			_, found := searchForVariableInArray(variableName, v.RequiredVariables)
			if found {
				Rules = append(Rules, k)
			}
		}
		return Rules, true
	}
	return []string{}, false
}

func searchForVariableInArray(val VariableFieldName, array []VariableFieldName) (int, bool) {
	for i, v := range array {
		if v == val {
			return i, true
		}
	}
	return -1, false
}

// MakeVariableValuePair creates a VariableValuePair from the variable name and its value
func MakeVariableValuePair(variable VariableFieldName, value []float64) VariableValuePair {
	return VariableValuePair{
		VariableName: variable,
		Values:       value,
	}
}

// GetValueFromRule evaluates a rule passed as an argument and returns the tax/allocation value
func GetValueFromRule(r RuleMatrix) (bool, float64, error) {
	variableVect := []float64{}
	variableVect = append(variableVect, 1)
	for range r.RequiredVariables {
		variableVect = append(variableVect, 1)
	}

	//Checking dimensions line up
	_, nCols := r.ApplicableMatrix.Dims()

	if nCols != len(variableVect) {
		return false, 0, errors.Errorf(
			"dimension mismatch in evaluating rule: '%v' rule matrix has '%v' columns, while we sourced '%v' variables",
			r.RuleName,
			nCols,
			len(variableVect),
		)
	}

	c := ruleMul(variableVect, r.ApplicableMatrix)

	resultVect, outputVal, err := genRealResult(r.AuxiliaryVector, c)
	if err != nil {
		return false, 0, err
	}

	return checkForFalse(resultVect), outputVal, nil

}
