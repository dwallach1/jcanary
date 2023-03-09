package interoperator

type VariableBag map[string]Variable

type VariableType string

const (
	Constant VariableType = "constant"
)

type Variable struct {
	Type VariableType
}

type MissingVariableErr struct {
	Var string
}

func (v *MissingVariableErr) Error() string {
	return "reference to undefined variable " + v.Var
}

// BuildString
//
//	template -> ${basurl}/checkins , bag -> { "baseurl": "localhost:8080" }
//	should return "localhost:8080/checkins"
func BuildString(template string, bag *VariableBag) string {

	// ${basurl}/checkins ->
	return ""
}
