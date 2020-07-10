package klayout

type Form struct {
	Rows []FormRow `json:"rows"`
}

type FormRow struct {
	RowIndex int          `json:"rowIndex"`
	Columns  []*FieldMeta `json:"columns"`
}
