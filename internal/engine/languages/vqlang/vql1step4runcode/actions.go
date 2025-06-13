package vql1step4runcode

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

func (tos *TTableOfSimbols) Set(dest TArgument, source *TArgument) bool {
	switch source.Term {
	case 0:
		dest.Type = source.Type
		dest.Value = source.Value
	case 1:
		ok, resultArgument := source.Productions[0].Exec(tos)
		if !ok.Result && len(resultArgument) > 1 {
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
		if !tos.Record(&dest) {
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

func (tos *TTableOfSimbols) SearchFunc(name string) (bool, TFunction) {
	if tos.Parent != nil {
		if stFunc, okFunc := tos.Parent.Function[name]; okFunc {
			return okFunc, stFunc
		} else {
			return tos.Parent.SearchFunc(name)
		}
	}

	return false, TFunction{}
}

// Production

func (actions TActions) RunCode(input TMapVariables) bool {
	pRootTOS := &TTableOfSimbols{
		Variables:   input,
		Transparent: false,
	}

	for _, production := range actions {
		if ok, _ := production.Exec(pRootTOS); !ok.Result {
			return ok.Result
		}
	}
	return true
}

func (parentProduction TProduction) Exec(parentTOS *TTableOfSimbols) (TReturn, []TArgument) {
	switch parentProduction.Type {
	case 1:
		// -- 1: НЕкод или комментарий
	case 2:
		// -- 2: инструкция компилятору
		switch parentProduction.Name {
		case "version":
			// пока пропускаем
		case "engine":
			// пока пропускаем
		}
	case 3, 4:
		// пока пропускаем
	case 11:
		// -- 11: простая самостоятельная продукция
		return TReturn{Result: false, Returned: false, Logic: false}, []TArgument{}
	case 12:
		// -- 12: блок области видимости
		pLocalTOS := &TTableOfSimbols{
			Parent:      parentTOS,
			Transparent: true,
		}
		actions := parentProduction.LocalCode[0]
		for _, production := range actions {
			if ok, res := production.Exec(pLocalTOS); !ok.Result {
				return ok, []TArgument{}
			} else if ok.Returned {
				return ok, res
			}
		}
	case 13:
		// -- 13: функция
		// правые аргументы - это входящие аргументы
		// левые аргументы - это выходные аргументы
		pLocalTOS := &TTableOfSimbols{
			Parent:      parentTOS,
			Transparent: false,
		}
		pLocalTOS.AddTOS(&parentProduction.Right)

		okFunc, stFunc := pLocalTOS.SearchFunc(parentProduction.Name)
		if !okFunc {
			return TReturn{Result: false, Returned: false}, []TArgument{}

		}

		iLenArg := len(stFunc.Input)
		if len(parentProduction.Right) != iLenArg {
			return TReturn{Result: false, Returned: false, Logic: false}, []TArgument{}
		}

		actions := stFunc.FuncCode
		for _, production := range actions {
			if ok, res := production.Exec(pLocalTOS); !ok.Result {
				return ok, []TArgument{}
			} else if ok.Returned {
				return ok, res
			}
		}
	case 14:
		// -- 14: возвратная операция
		if len(parentProduction.Right) == 0 {
			return TReturn{Result: true, Returned: true, Logic: false}, []TArgument{}
		} else {
			var slResArg []TArgument
			for _, templateArg := range parentProduction.Right {
				if ok, resArg := parentTOS.Extact(templateArg); ok {
					slResArg = append(slResArg, resArg)
				} else {
					return TReturn{Result: false, Returned: true, Logic: false}, slResArg
				}
			}
			return TReturn{Result: true, Returned: true, Logic: false}, slResArg
		}
	case 15:
		// -- 15: условная операция
		// конструкция "else" должна быть реализована в продукции "elseif 0 == 0 {}"
		if len(parentProduction.Right) > 0 {
			for iIf, stIf := range parentProduction.Right {
				ifRet, ifRes := stIf.Productions[0].Exec(parentTOS)
				if ifRet.Result && ifRes[0].Value == "true" {
					pLocalTOS := &TTableOfSimbols{
						Parent:      parentTOS,
						Transparent: true,
					}
					actions := parentProduction.LocalCode[iIf]
					for _, production := range actions {
						if ok, res := production.Exec(pLocalTOS); !ok.Result {
							return ok, []TArgument{}
						} else if ok.Returned {
							return ok, res
						}
					}
				}
			}
		} else {
			return TReturn{Result: false, Returned: false, Logic: false}, []TArgument{}
		}
	case 16:
		// -- 16: классический цикл
		actions := parentProduction.LocalCode[0]
		if parentProduction.Right[0].Term == 1 && parentProduction.Right[1].Term == 1 && parentProduction.Right[2].Term == 1 {
			pLocalTOS := &TTableOfSimbols{
				Parent:      parentTOS,
				Transparent: true,
			}
			parentProduction.Right[0].Productions[0].Exec(pLocalTOS)
		labelfor:
			if ifRet, ifRes := parentProduction.Right[1].Productions[0].Exec(pLocalTOS); ifRet.Result && ifRet.Logic && ifRes[0].Value == "true" {
				for _, production := range actions {
					if ok, res := production.Exec(pLocalTOS); !ok.Result {
						return ok, []TArgument{}
					} else if ok.Returned {
						return ok, res
					}
				}
			} else {
				return TReturn{Result: true, Returned: false, Logic: false}, []TArgument{}
			}
			parentProduction.Right[2].Productions[0].Exec(pLocalTOS)
			goto labelfor
		} else {
			return TReturn{Result: false, Returned: false, Logic: false}, []TArgument{}
		}
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
					if okProd, resProd := prodVal.Exec(parentTOS); okProd.Result {
						slUnpackedArguments = append(slUnpackedArguments, resProd...)
					} else {
						return TReturn{Result: false, Returned: false, Logic: false}, []TArgument{}
					}
				}
			}
		}
		if len(parentProduction.Left) == len(slUnpackedArguments) {
			for indRA, rightArg := range slUnpackedArguments {
				if !parentTOS.Set(parentProduction.Left[indRA], &rightArg) {
					return TReturn{Result: false, Returned: false, Logic: false}, []TArgument{}
				}
			}
		} else {
			return TReturn{Result: false, Returned: false, Logic: false}, []TArgument{}
		}
	case 42:
		// -- 42: разделение голов и хвостов |=
	}

	return TReturn{Result: true, Returned: false, Logic: false}, []TArgument{}
}
