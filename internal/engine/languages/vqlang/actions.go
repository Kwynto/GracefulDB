package vqlang

// Manage tables of simbols

func (tos *TTableOfSimbols) Record(dest *TArgument) bool {
	if _, okV := tos.Variables[dest.Simbol]; okV {
		tos.Variables[dest.Simbol] = TVariableData{
			Type:  dest.Type,
			Value: dest.Value,
		}
	} else if pointer, okP := tos.Pointers[dest.Simbol]; okP {
		pointer.Link.Variables[pointer.Simbol] = TVariableData{
			Type:  dest.Type,
			Value: dest.Value,
		}
	} else if tos.Transparent {
		return tos.Parent.Record(dest)
	} else {
		return false
	}

	return true
}

func (tos *TTableOfSimbols) Set(dest *TArgument, source *TArgument) bool {
	switch source.Term {
	case 0:
		dest.Type = source.Type
		dest.Value = source.Value
	case 1:
		ok, resultArgument := source.Productions[0].Exec(tos)
		if !ok && len(resultArgument) > 1 {
			return false
		}
		dest.Value = resultArgument[0].Value
	case 2:
		dest.Term = source.Term
		dest.Type = source.Type
		dest.Value = source.Value
	case 3:
		dest.Term = source.Term
		dest.Value = source.Value
		dest.Link = source.Link
	}

	if dest.Simbol != "" {
		if !tos.Record(dest) {
			switch dest.Term {
			case 2:
				tos.Variables[dest.Simbol] = TVariableData{
					Type:  dest.Type,
					Value: dest.Value,
				}
			case 3:
				tos.Pointers[dest.Simbol] = TPointer{
					Simbol: dest.Value,
					Link:   dest.Link,
				}
			}
		}
	} else {
		return false
	}

	return true
}

// дополнение таблицы символов набором новых переменных
func (tos *TTableOfSimbols) AddTOS(input *[]TArgument) {
	for _, source := range *input {
		switch source.Term {
		case 2:
			stVar := TVariableData{
				Type:  source.Type,
				Value: source.Value,
			}
			tos.Variables[source.Simbol] = stVar
		case 3:
			stPoint := TPointer{
				Simbol: source.Value,
				Link:   source.Link,
			}
			tos.Pointers[source.Simbol] = stPoint
		}
	}
}

func (tos *TTableOfSimbols) Extact(template TArgument) (bool, TArgument) {
	if varData, okV := tos.Variables[template.Simbol]; okV {
		resArg := TArgument{
			Term:   2,
			Simbol: template.Simbol,
			Type:   varData.Type,
			Value:  varData.Value,
		}
		return true, resArg
	} else if pointData, okP := tos.Pointers[template.Simbol]; okP {
		resArg := TArgument{
			Term:  3,
			Value: pointData.Simbol,
			Link:  pointData.Link,
		}
		return true, resArg
	} else if tos.Transparent {
		return tos.Parent.Extact(template)
	}

	return false, TArgument{}
}

// Production

func (actions TActions) RunCode(input TMapVariables) bool {
	pRootTOS := &TTableOfSimbols{
		Variables:   input,
		Transparent: false,
	}

	for _, production := range actions {
		if ok, _ := production.Exec(pRootTOS); !ok {
			return ok
		}
	}
	return true
}

func (parentProduction TProduction) Exec(parentTOS *TTableOfSimbols) (bool, []TArgument) {
	switch parentProduction.Type {
	case 1:
		// -- 1: НЕкод или комментарий
	case 2, 3, 4:
		// пока пропускаем
	case 11:
		// -- 11: простая самостоятельная продукция
		return false, []TArgument{}
	case 12:
		// -- 12: блок области видимости
		pLocalTOS := &TTableOfSimbols{
			Parent:      parentTOS,
			Transparent: true,
		}
		actions := parentProduction.LocalCode
		for _, production := range actions {
			if ok, res := production.Exec(pLocalTOS); !ok {
				return false, []TArgument{}
			} else if len(res) != 0 {
				return true, res
			}
		}
	case 13:
		// -- 13: функция
		// правые аргументы - это входящие аргументы
		// левые аргументы - это выходные аргументы
		pLocalTOS := &TTableOfSimbols{
			// Parent:      parentTOS,
			Transparent: false,
		}
		pLocalTOS.AddTOS(&parentProduction.Right)
		actions := parentProduction.LocalCode
		for _, production := range actions {
			if ok, res := production.Exec(pLocalTOS); !ok {
				return false, []TArgument{}
			} else if len(res) != 0 {
				return true, res
			}
		}
	case 14:
		// -- 14: возвратная операция
		if len(parentProduction.Right) == 0 {
			return true, []TArgument{}
		} else {
			var slResArg []TArgument
			for _, templateArg := range parentProduction.Right {
				if ok, resArg := parentTOS.Extact(templateArg); ok {
					slResArg = append(slResArg, resArg)
				} else {
					return false, slResArg
				}
			}
			return true, slResArg
		}
	case 15:
		// -- 15: условная операция
	case 16:
		// -- 16: классический цикл
	case 17:
		// -- 17: цикл по диапазону
	case 41:
		// -- 41: присваивание :=
		var slUnpackedArguments []TArgument
		for _, argVal := range parentProduction.Right {
			switch argVal.Term {
			case 0, 2, 3:
				slUnpackedArguments = append(slUnpackedArguments, argVal)
			case 1:
				for _, prodVal := range argVal.Productions {
					if okProd, resProd := prodVal.Exec(parentTOS); okProd {
						slUnpackedArguments = append(slUnpackedArguments, resProd...)
					} else {
						return false, []TArgument{}
					}
				}
			}
		}
		parentProduction.Right = slUnpackedArguments
		if len(parentProduction.Left) == len(parentProduction.Right) {
			for indRA, rightArg := range parentProduction.Right {
				if !parentTOS.Set(&parentProduction.Left[indRA], &rightArg) {
					return false, []TArgument{}
				}
			}
		} else {
			return false, []TArgument{}
		}
	}

	return true, []TArgument{}
}
