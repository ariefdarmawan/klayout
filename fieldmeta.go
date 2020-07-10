package klayout

type FieldMeta struct {
	ID        string `json:"id"`
	Field     string `json:"field"`
	Label     string `json:"label"`
	GridShow  string `json:"gridShow"`
	GridWidth string `json:"gridWidth"`
	Format    string `json:"format"`
	Align     string `json:"align"`
	MultiRow  int    `json:"multiRow"`

	Required bool `json:"required"`
	ReadOnly bool `json:"readOnly"`

	MinLength int         `json:"minLength"`
	MaxLength int         `json:"maxLength"`
	MinValue  interface{} `json:"minValue"`
	MaxValue  interface{} `json:"maxValue"`

	FieldType string `json:"fieldType"`
	Row       int    `json:"row"`
	Col       int    `json:"col"`
	FormShow  string `json:"formShow"`
	Masked    bool   `json:"masked"`

	Control string `json:"control"`

	//-- List
	AllowAdd     bool     `json:"allowAdd"`
	UseList      bool     `json:"useList"`
	ListItems    []string `json:"listItems"`
	LookupURL    string   `json:"lookupUrl"`
	LookupKey    string   `json:"lookupKey"`
	LookupFields []string `json:"lookupFields"`
}
