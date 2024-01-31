package sqlanalyzer

// DDL — язык определения данных (Data Definition Language)

func DDLCreate(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DDLAlter(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}

func DDLDrop(instruction *string, placeholder *[]string) (result *string, remainder *string, err error) {
	res := ""
	result = &res
	remainder = instruction
	return result, remainder, nil
}
