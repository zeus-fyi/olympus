package postgres

func (d *PgSchemaAutogen) GetTableData() error {
	tables, err := d.GetTables()
	if err != nil {
		return err
	}
	err = d.ProcessTables(d.Postgresql, d.Settings, tables...)
	return err
}
