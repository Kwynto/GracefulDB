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
	} else if !tos.IsRoot {
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
		if !ok {
			return false
		}
		dest.Value = resultArgument.Value
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

// Production

func (actions TActions) Main(input TMapVariables) bool {
	pRootTOS := &TTableOfSimbols{
		Input:  input,
		IsRoot: true,
	}

	return actions.Run(pRootTOS)
}

func (actions TActions) Run(parentTOS *TTableOfSimbols) bool {
	pSelfTOS := &TTableOfSimbols{
		Parent: parentTOS,
		IsRoot: false,
	}
	for _, production := range actions {
		if ok, _ := production.Exec(pSelfTOS); !ok {
			return ok
		}
	}
	return true
}

func (parentProduction TProduction) Exec(parentTOS *TTableOfSimbols) (bool, *TArgument) {
	switch parentProduction.Type {
	case 1:
		// -- 1: НЕкод или комментарий
	case 2, 3, 4:
		// пока пропускаем
	case 11:
		// -- 11: простая вложенная продукция
		zeroArgument := parentProduction.Left[0]
		if zeroArgument.Term == 1 {
			if ok, _ := zeroArgument.Productions[0].Exec(parentTOS); !ok {
				return false, &TArgument{}
			}
		}
	case 12:
		// -- 12: блок области видимости
		pLocalTOS := &TTableOfSimbols{
			Parent: parentTOS,
			IsRoot: false,
		}
		return parentProduction.Left[0].Productions.Run(pLocalTOS), &TArgument{}
	case 13, 14, 15, 16, 17:
		// пока пропустить
	case 41:
		// -- 41: присваивание :=
		if len(parentProduction.Left) == len(parentProduction.Right) {
			for indRA, rightArg := range parentProduction.Right {
				if !parentTOS.Set(&parentProduction.Left[indRA], &rightArg) {
					return false, &TArgument{}
				}
			}
		} else {
			return false, &TArgument{}
		}
	}

	return true, &TArgument{}
}
