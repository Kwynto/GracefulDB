package sqlanalyzer

// DCL — язык управления данными (Data Control Language)

func DCLGrant(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DCLRevoke(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DCLUse(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}
