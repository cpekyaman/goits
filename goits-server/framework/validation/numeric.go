package validation

func IntMin(min int32) validatorImpl {
	return validatorImpl{
		name:   "min",
		params: map[string]interface{}{"min": min},
		vFunc: intPredicateFunc(func(intVal int32) bool {
			return intVal >= min
		}),
	}
}

func IntMax(max int32) validatorImpl {
	return validatorImpl{
		name:   "max",
		params: map[string]interface{}{"max": max},
		vFunc: intPredicateFunc(func(intVal int32) bool {
			return intVal <= max
		}),
	}
}

func IntRange(min int32, max int32) validatorImpl {
	return validatorImpl{
		name:   "range",
		params: map[string]interface{}{"min": min, "max": max},
		vFunc: intPredicateFunc(func(intVal int32) bool {
			return min <= intVal && intVal <= max
		}),
	}
}

func intPredicateFunc(pred func(int32) bool) func(interface{}) bool {
	return func(value interface{}) bool {
		in := fieldValue(value)
		intVal, ok := in.(int32)
		if !ok {
			return false
		}
		return pred(intVal)
	}
}
