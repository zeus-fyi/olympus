package hera

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/fraenky8/tables-to-go/pkg/settings"
	"github.com/jmoiron/sqlx"

	// postgres database driver
	_ "github.com/lib/pq"
)

var (
	// dbTypeToDriverMap maps the database type to the driver names.
	dbTypeToDriverMap = map[settings.DBType]string{
		settings.DBTypePostgresql: "postgres",
		settings.DBTypeMySQL:      "mysql",
		settings.DBTypeSQLite:     "sqlite3",
	}
)

// Postgresql implements the Database interface with help of GeneralDatabase.
type Postgresql struct {
	*GeneralDatabase
	defaultUserName string
}

// Table has a name and a set (slice) of columns.
type Table struct {
	Name    string `db:"table_name"`
	Columns []Column
}

// Column stores information about a column.
type Column struct {
	OrdinalPosition        int            `db:"ordinal_position"`
	Name                   string         `db:"column_name"`
	DataType               string         `db:"data_type"`
	DefaultValue           sql.NullString `db:"column_default"`
	IsNullable             string         `db:"is_nullable"`
	CharacterMaximumLength sql.NullInt64  `db:"character_maximum_length"`
	NumericPrecision       sql.NullInt64  `db:"numeric_precision"`
	ColumnKey              string         `db:"column_key"`      // mysql specific
	Extra                  string         `db:"extra"`           // mysql specific
	ConstraintName         sql.NullString `db:"constraint_name"` // pg specific
	ConstraintType         sql.NullString `db:"constraint_type"` // pg specific
}

// Connect connects to the database by the given data source name (dsn) of the
// concrete database.
func (pg *Postgresql) Connect(dsn string) error {
	return pg.GeneralDatabase.Connect(dsn)
}

// GeneralDatabase represents a base "class" database - for all other concrete
// databases it implements partly the Database interface.
type GeneralDatabase struct {
	GetColumnsOfTableStmt *sqlx.Stmt
	*sqlx.DB
	*settings.Settings
	driver string
}

// Connect establishes a connection to the database with the given DSN.
// It pings the database to ensure it is reachable.
func (gdb *GeneralDatabase) Connect(dsn string) (err error) {
	gdb.DB, err = sqlx.Connect(gdb.driver, dsn)
	if err != nil {
		usingPswd := "no"
		if gdb.Settings.Pswd != "" {
			usingPswd = "yes"
		}
		return fmt.Errorf(
			"could not connect to database (type=%q, user=%q, database=%q, host='%v:%v', using password: %v): %w",
			gdb.DbType, gdb.User, gdb.DbName, gdb.Host, gdb.Port, usingPswd, err,
		)
	}

	return gdb.Ping()
}

// GetTables gets all tables for a given schema by name.
func (pg *Postgresql) GetTables() (tables []*Table, err error) {

	err = pg.Select(&tables, `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_type = 'BASE TABLE'
		AND table_schema = $1
		ORDER BY table_name
	`, pg.Schema)

	if pg.Verbose {
		if err != nil {
			fmt.Println("> Error at GetTables()")
			fmt.Printf("> schema: %q\r\n", pg.Schema)
		}
	}

	return tables, err
}

// PrepareGetColumnsOfTableStmt prepares the statement for retrieving the
// columns of a specific table for a given database.
func (pg *Postgresql) PrepareGetColumnsOfTableStmt() (err error) {

	pg.GetColumnsOfTableStmt, err = pg.Preparex(`
		SELECT
			ic.ordinal_position,
			ic.column_name,
			ic.data_type,
			ic.column_default,
			ic.is_nullable,
			ic.character_maximum_length,
			ic.numeric_precision,
			itc.constraint_name,
			itc.constraint_type
		FROM information_schema.columns AS ic
			LEFT JOIN information_schema.key_column_usage AS ikcu ON ic.table_name = ikcu.table_name
			AND ic.table_schema = ikcu.table_schema
			AND ic.column_name = ikcu.column_name
			LEFT JOIN information_schema.table_constraints AS itc ON ic.table_name = itc.table_name
			AND ic.table_schema = itc.table_schema
			AND ikcu.constraint_name = itc.constraint_name
		WHERE ic.table_name = $1
		AND ic.table_schema = $2
		ORDER BY ic.ordinal_position
	`)

	return err
}

// GetColumnsOfTable executes the statement for retrieving the columns of a
// specific table in a given schema.
func (pg *Postgresql) GetColumnsOfTable(table *Table) (err error) {

	err = pg.GetColumnsOfTableStmt.Select(&table.Columns, table.Name, pg.Schema)

	if pg.Verbose {
		if err != nil {
			fmt.Printf("> Error at GetColumnsOfTable(%v)\r\n", table.Name)
			fmt.Printf("> schema: %q\r\n", pg.Schema)
		}
	}

	return err
}

// IsPrimaryKey checks if the column belongs to the primary key.
func (pg *Postgresql) IsPrimaryKey(column Column) bool {
	return strings.Contains(column.ConstraintType.String, "PRIMARY KEY")
}

// IsAutoIncrement checks if the column is an auto_increment column.
func (pg *Postgresql) IsAutoIncrement(column Column) bool {
	return strings.Contains(column.DefaultValue.String, "nextval")
}

// GetStringDatatypes returns the string datatypes for the Postgresql database.
func (pg *Postgresql) GetStringDatatypes() []string {
	return []string{
		"character varying",
		"varchar",
		"character",
		"char",
		"uuid",
	}
}

// IsString returns true if colum is of type string for the Postgresql database.
func (pg *Postgresql) IsString(column Column) bool {
	return isStringInSlice(column.DataType, pg.GetStringDatatypes())
}

// GetTextDatatypes returns the text datatypes for the Postgresql database.
func (pg *Postgresql) GetTextDatatypes() []string {
	return []string{
		"text",
	}
}

// IsText returns true if colum is of type text for the Postgresql database.
func (pg *Postgresql) IsText(column Column) bool {
	return isStringInSlice(column.DataType, pg.GetTextDatatypes())
}

// GetIntegerDatatypes returns the integer datatypes for the Postgresql database.
func (pg *Postgresql) GetIntegerDatatypes() []string {
	return []string{
		"smallint",
		"integer",
		"bigint",
		"smallserial",
		"serial",
		"bigserial",
	}
}

// IsInteger returns true if colum is of type integer for the Postgresql database.
func (pg *Postgresql) IsInteger(column Column) bool {
	return isStringInSlice(column.DataType, pg.GetIntegerDatatypes())
}

// GetFloatDatatypes returns the float datatypes for the Postgresql database.
func (pg *Postgresql) GetFloatDatatypes() []string {
	return []string{
		"numeric",
		"decimal",
		"real",
		"double precision",
	}
}

// IsFloat returns true if colum is of type float for the Postgresql database.
func (pg *Postgresql) IsFloat(column Column) bool {
	return isStringInSlice(column.DataType, pg.GetFloatDatatypes())
}

// GetTemporalDatatypes returns the temporal datatypes for the Postgresql database.
func (pg *Postgresql) GetTemporalDatatypes() []string {
	return []string{
		"time",
		"timestamp",
		"time with time zone",
		"timestamp with time zone",
		"time without time zone",
		"timestamp without time zone",
		"date",
	}
}

// IsTemporal returns true if colum is of type temporal for the Postgresql database.
func (pg *Postgresql) IsTemporal(column Column) bool {
	return isStringInSlice(column.DataType, pg.GetTemporalDatatypes())
}

// isStringInSlice checks if needle (string) is in haystack ([]string).
func isStringInSlice(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
