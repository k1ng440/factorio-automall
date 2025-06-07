package blueprint

// Factorio represents the top-level structure containing a blueprint
type Factorio struct {
	Blueprint *Blueprint `json:"blueprint"`
}

type Blueprint struct {
	Description string    `json:"description"`
	Icons       []*Icon   `json:"icons"`
	Entities    []*Entity `json:"entities"`
	Wires       [][]int   `json:"wires,omitempty"`
	Item        string    `json:"item"`
	Version     string    `json:"version"`
}

type Entity struct {
	EntityNumber      int              `json:"entity_number"`
	Name              string           `json:"name"`
	Position          Position         `json:"position"`
	Direction         *int64           `json:"direction,omitempty"` // Enum: [0, 1, 2, 3, 4, 5, 6, 7], omitted if 0
	ControlBehavior   *ControlBehavior `json:"control_behavior"`
	PlayerDescription string           `json:"player_description"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Icon struct {
	Signal *Signal `json:"signal"`
	Index  int     `json:"index"`
}

type ControlBehavior struct {
	DeciderConditions *DeciderConditions  `json:"decider_conditions,omitempty"` // For decider-combinator
	Sections          *ConstantCombinator `json:"sections,omitempty"`           // For constant-combinator
}

// DeciderConditions represents conditions for a decider combinator.
type DeciderConditions struct {
	Conditions []*DeciderConditionsCondition `json:"conditions"`
	Outputs    []*DeciderConditionOutputs    `json:"outputs"`
}

type DeciderConditionsCondition struct {
	FirstSignal          *Signal         `json:"first_signal,omitempty"`
	Comparator           string          `json:"comparator,omitempty"`
	FirstSignalNetworks  *SignalNetworks `json:"first_signal_networks"`
	CompareType          string          `json:"compare_type,omitempty"`
	SecondSignal         *Signal         `json:"second_signal,omitempty"`
	SecondSignalNetworks *SignalNetworks `json:"second_signal_networks,omitempty"`
	Constant             int             `json:"constant"`
}

type Signal struct {
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Quality string `json:"quality,omitempty"`
}

type SignalNetworks struct {
	Red   bool `json:"red"`
	Green bool `json:"green"`
}

// CombinatorOutput represents output configuration for decider conditions
type DeciderConditionOutputs struct {
	Signal             *Signal `json:"signal"`
	CopyCountFromInput bool    `json:"copy_count_from_input"`
	Constant           int     `json:"constant"`
}

type ConstantCombinator struct {
	Sections []*ConstantCombinatorSection `json:"sections,omitempty"`
}

type ConstantCombinatorSection struct {
	Index   int                                `json:"index"`
	Filters []*ConstantCombinatorSectionFilter `json:"filters"`
}

// ConstantFilter represents a filter for a constant combinator.
type ConstantCombinatorSectionFilter struct {
	Index      int    `json:"index"`
	Name       string `json:"name,omitempty"`
	Quality    string `json:"quality,omitempty"`
	Comparator string `json:"comparator,omitempty"`
	Count      int    `json:"count"`
}
