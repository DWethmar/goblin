package aggr

type MockCommand struct {
	aggregateID   string
	aggregateType string
	commandType   string
}

func (c MockCommand) CommandType() string   { return c.commandType }
func (c MockCommand) AggregateID() string   { return c.aggregateID }
func (c MockCommand) AggregateType() string { return c.aggregateType }
