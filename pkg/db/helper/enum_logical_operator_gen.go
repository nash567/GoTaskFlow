// Code generated by "enumer -type=LogicalOperator -text -json -transform=snake -output=enum_logical_operator_gen.go"; DO NOT EDIT.

package helper

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _LogicalOperatorName = "logical_operator_andlogical_operator_or"

var _LogicalOperatorIndex = [...]uint8{0, 20, 39}

const _LogicalOperatorLowerName = "logical_operator_andlogical_operator_or"

func (i LogicalOperator) String() string {
	if i < 0 || i >= LogicalOperator(len(_LogicalOperatorIndex)-1) {
		return fmt.Sprintf("LogicalOperator(%d)", i)
	}
	return _LogicalOperatorName[_LogicalOperatorIndex[i]:_LogicalOperatorIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _LogicalOperatorNoOp() {
	var x [1]struct{}
	_ = x[LogicalOperatorAnd-(0)]
	_ = x[LogicalOperatorOr-(1)]
}

var _LogicalOperatorValues = []LogicalOperator{LogicalOperatorAnd, LogicalOperatorOr}

var _LogicalOperatorNameToValueMap = map[string]LogicalOperator{
	_LogicalOperatorName[0:20]:       LogicalOperatorAnd,
	_LogicalOperatorLowerName[0:20]:  LogicalOperatorAnd,
	_LogicalOperatorName[20:39]:      LogicalOperatorOr,
	_LogicalOperatorLowerName[20:39]: LogicalOperatorOr,
}

var _LogicalOperatorNames = []string{
	_LogicalOperatorName[0:20],
	_LogicalOperatorName[20:39],
}

// LogicalOperatorString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func LogicalOperatorString(s string) (LogicalOperator, error) {
	if val, ok := _LogicalOperatorNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _LogicalOperatorNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to LogicalOperator values", s)
}

// LogicalOperatorValues returns all values of the enum
func LogicalOperatorValues() []LogicalOperator {
	return _LogicalOperatorValues
}

// LogicalOperatorStrings returns a slice of all String values of the enum
func LogicalOperatorStrings() []string {
	strs := make([]string, len(_LogicalOperatorNames))
	copy(strs, _LogicalOperatorNames)
	return strs
}

// IsALogicalOperator returns "true" if the value is listed in the enum definition. "false" otherwise
func (i LogicalOperator) IsALogicalOperator() bool {
	for _, v := range _LogicalOperatorValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for LogicalOperator
func (i LogicalOperator) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for LogicalOperator
func (i *LogicalOperator) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("LogicalOperator should be a string, got %s", data)
	}

	var err error
	*i, err = LogicalOperatorString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for LogicalOperator
func (i LogicalOperator) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for LogicalOperator
func (i *LogicalOperator) UnmarshalText(text []byte) error {
	var err error
	*i, err = LogicalOperatorString(string(text))
	return err
}