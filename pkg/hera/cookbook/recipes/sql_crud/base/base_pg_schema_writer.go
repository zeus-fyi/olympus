package base

func (m *ModelTemplate) WritePgTableDefinition() error {
	err := m.GetTableData()
	if err != nil {
		return err
	}

	for tbl, s := range m.StructMapToCodeGen {
		m.Path.AddGoFn(tbl)
		err = m.CreateTemplateFromStruct(s)
		if err != nil {
			return err

		}
	}
	return err
}
