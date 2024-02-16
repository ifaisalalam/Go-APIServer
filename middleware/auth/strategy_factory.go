package auth

var supportedStrategies []Strategy

func GetStrategy(name StrategyName) Strategy {
	for _, strategy := range supportedStrategies {
		if strategy.Is(name) {
			return strategy
		}
	}
	return nil
}
