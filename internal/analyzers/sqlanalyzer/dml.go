package sqlanalyzer

// DML — язык изменения данных (Data Manipulation Language)

func DMLSelect(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DMLInsert(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DMLUpdate(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DMLDelete(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DMLTruncate(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DMLCommit(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DMLRollback(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}
