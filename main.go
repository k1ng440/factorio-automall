package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/k1ng440/factorio-automall/blueprint"
	"github.com/k1ng440/factorio-automall/recipe"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename> <quality=true/false>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

    useQuality, err := strconv.ParseBool(os.Args[2])
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: quality must be 'true' or 'false', got %s\n", os.Args[2])
        os.Exit(1)
    }

	data, err := recipe.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	if len(data.Recipes) == 0 {
		fmt.Println("No recipes found in this file")
	}

	deciderConditions := &blueprint.DeciderConditions{
		Conditions: []*blueprint.DeciderConditionsCondition{ },
		Outputs: []*blueprint.DeciderConditionOutputs{
			{
				Signal: &blueprint.Signal{
					Type: "virtual",
					Name: "signal-each",
				},
				Constant:           -1,
				CopyCountFromInput: false,
			},
		},
	}

	qualities := []string{"normal"}
	if useQuality {
		qualities = []string{
			"normal",
			"uncommon",
			"rare",
			"epic",
			"legendary",
		}
	}

	constantCombinator := &blueprint.ConstantCombinator{}
	constantCombinatorSections := make(map[string]*blueprint.ConstantCombinatorSection)
	for idx, quality := range qualities {
		constantCombinatorSections[quality] = &blueprint.ConstantCombinatorSection{
			Index:   idx + 1,
			Filters: make([]*blueprint.ConstantCombinatorSectionFilter, 0),
		}
	}

	bp := &blueprint.Factorio{
		Blueprint: &blueprint.Blueprint{
			Description: "Auto generated with script created by k1ng440",
			Icons: []*blueprint.Icon{
				{
					Index: 1,
					Signal: &blueprint.Signal{
						Name: "decider-combinator",
					},
				},
				{
					Index: 2,
					Signal: &blueprint.Signal{
						Name: "constant-combinator",
					},
				},
				{
					Index: 3,
					Signal: &blueprint.Signal{
						Type: "virtual",
						Name: "signal-any-quality",
					},
				},
			},
			Entities: []*blueprint.Entity{
				{
					EntityNumber: 1,
					Name:         "decider-combinator",
					Position: blueprint.Position{
						X: 278.5,
						Y: 63,
					},
					ControlBehavior: &blueprint.ControlBehavior{
						DeciderConditions: deciderConditions,
					},
					PlayerDescription: "Auto generated with script created by k1ng440",
				},
				{
					EntityNumber: 2,
					Name:         "constant-combinator",
					Position: blueprint.Position{
						X: 278.5,
						Y: 64.5,
					},
					ControlBehavior: &blueprint.ControlBehavior{
						Sections: constantCombinator,
					},
					PlayerDescription: "Auto generated with script created by k1ng440",
				},
			},
			Wires: [][]int{
				{1, 1, 2, 1},
			},
			Item:    "blueprint",
			Version: "562949957025792",
		},
	}

	fludType := make(map[string]struct{})
	for _, item := range data.Items {
		if item.Type == "fluid" {
			fludType[item.Key] = struct{}{}
		}
	}


	// itemValue represents the unique value set in the constant combinator, used by the decider combinator to output the item can be crafted.
	itemValue := -5000000

	idx := map[string]int{
		"legendary": 0,
		"epic":      0,
		"rare":      0,
		"uncommon":  0,
		"normal":    0,
	}

	for _, rec := range data.Recipes {
		if shouldIgnore(rec) {
			continue
		}

		// Ignore fluid
		if _, ok := fludType[rec.Key]; ok {
			continue
		}

		for _, quality := range qualities {
			itemValue -= 5000

			deciderConditions.Conditions = append(deciderConditions.Conditions,
				&blueprint.DeciderConditionsCondition{
					FirstSignal: &blueprint.Signal{
						Type: "virtual",
						Name: "signal-each",
					},
					SecondSignal: &blueprint.Signal{
						Name:    rec.Key,
						Quality: quality,
					},
					Comparator: "=",
					FirstSignalNetworks: &blueprint.SignalNetworks{
						Red:   true,
						Green: false,
					},
					SecondSignalNetworks: &blueprint.SignalNetworks{
						Red:   true,
						Green: false,
					},
				},
			)

			idx[quality]++
			constantCombinatorSections[quality].Filters = append(constantCombinatorSections[quality].Filters, &blueprint.ConstantCombinatorSectionFilter{
				Index:      idx[quality],
				Name:       rec.Key,
				Quality:    quality,
				Comparator: "=",
				Count:      itemValue, // negative value
			})

			for _, ingredient := range rec.Ingredients {
				itemType := ""
				lQuality := quality
				if _, ok := fludType[ingredient.Name]; ok {
					itemType = "fluid"
					lQuality = ""
				}

				deciderConditions.Conditions = append(deciderConditions.Conditions, &blueprint.DeciderConditionsCondition{
					FirstSignal: &blueprint.Signal{
						Type:    itemType,
						Name:    ingredient.Name,
						Quality: lQuality,
					},
					Constant:   int(ingredient.Amount),
					Comparator: "≥",
					FirstSignalNetworks: &blueprint.SignalNetworks{
						Red:   false,
						Green: true,
					},
					CompareType: "and",
				})
			}
		}
	}

	for _, section := range constantCombinatorSections {
		bp.Blueprint.Entities[1].ControlBehavior.Sections.Sections = append(bp.Blueprint.Entities[1].ControlBehavior.Sections.Sections, section)
	}

	j, err := json.MarshalIndent(bp, " ", " ")
	if err != nil {
		fmt.Printf("Failed to marshal json. Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(j))
}

func shouldIgnore(rec *recipe.Recipe) bool {
	if rec == nil {
		return true
	}

	ignoreKey := map[string]bool{
		"copper-bacteria-cultivation": true,
		"infinity-pipe": true,
		"infinity-chest": true,
		"ice": true,
		"ice-platform": true,
		"heat-interface": true,
	}

	ignoreCategories := map[string]bool{
		"recycling":                  true,
		"recycling-or-hand-crafting": true,
		"crushing":                   true,
		"centrifuging":               true,
		"metallurgy":                 true,
		"smelting":                   true,
		"rocket-building":            true,
		"organic-or-hand-crafting":   true,
	}

	ignoreSubgroups := map[string]bool{
		"fluid-recipes":         true,
		"space-processing":      true,
		"agriculture":           true,
		"aquilo-processes":      true,
		"agriculture-products":  true,
		"raw-material":          true,
		"nauvis-agriculture":    true,
		"agriculture-processes": true,
	}

	return ignoreKey[rec.Key] || ignoreCategories[rec.Category] || ignoreSubgroups[rec.Subgroup]
}
