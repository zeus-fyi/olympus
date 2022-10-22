package base

func (m *ModelTemplate) ReadPgTableDefinition() error {
	err := m.GetTableData()
	return err
}
