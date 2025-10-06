package shell

import "strings"

type Command struct {
	Name string
	Args []string
}

type ConditionalCommand struct {
	Cmd         Command
	Conditional string
}

func ParseConditional(line string) []ConditionalCommand {
	var result []ConditionalCommand
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return result
	}

	var current []string
	var cond string
	for _, tok := range tokens {
		if tok == "&&" || tok == "||" {
			if len(current) > 0 {
				result = append(result, ConditionalCommand{
					Cmd: Command{
						Name: current[0],
						Args: current[1:],
					},
					Conditional: cond,
				})
				current = nil
			}
			cond = tok
		} else {
			current = append(current, tok)
		}
	}
	if len(current) > 0 {
		result = append(result, ConditionalCommand{
			Cmd: Command{
				Name: current[0],
				Args: current[1:],
			},
			Conditional: cond,
		})
	}

	return result
}
