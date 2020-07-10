package klayout

import (
	"reflect"
	"strings"

	"github.com/eaciit/toolkit"
)

type UIModel struct {
	Name           string
	FormTitle      string
	FormNameFields []string
	SearchFields   []string
	KeyFields      []string

	model      interface{}
	fieldMetas map[string]*FieldMeta
	fieldNames []string
}

func NewUIModel(obj interface{}) (*UIModel, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	t := v.Type()

	ui := new(UIModel)
	ui.Name = t.Name()

	fieldCount := t.NumField()
	ui.fieldMetas = make(map[string]*FieldMeta, fieldCount)
	for idx := 0; idx < fieldCount; idx++ {
		tv := v.Field(idx)
		tf := t.Field(idx)
		getMetaFromStructField(tv, tf, ui)
	}
	ui.model = obj

	return ui, nil
}

func getMetaFromStructField(v reflect.Value, tf reflect.StructField, mdl *UIModel) *FieldMeta {
	if !v.CanSet() {
		return nil
	}

	tag := tf.Tag
	tagFieldName := tag.Get(toolkit.TagName())
	if tagFieldName == "-" {
		return nil
	}

	defaultGridShow := "show"
	fm := new(FieldMeta)
	fm.Field = tf.Name
	fm.Label = tagDefault(tag, "name", fm.Field)

	if formTitle := tagDefault(tag, "kf-title", ""); formTitle != "" && mdl.FormTitle == "" {
		mdl.FormTitle = formTitle
	}

	if formNames := tagDefault(tag, "kf-name", ""); formNames != "" && len(mdl.FormNameFields) == 0 {
		fields := strings.Split(formNames, ",")
		for _, field := range fields {
			mdl.FormNameFields = append(mdl.FormNameFields, strings.Trim(field, " "))
		}
	}

	label := tag.Get("label")
	if label == "" {
		label = tf.Name
	}

	var kind reflect.Kind
	var ft reflect.Type
	if tf.Type.Kind() == reflect.Ptr {
		ft = tf.Type.Elem()
		kind = tf.Type.Elem().Kind()
	} else {
		ft = tf.Type
		kind = tf.Type.Kind()
	}

	fieldType := ""
	if kind == reflect.Int || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 ||
		kind == reflect.Int8 || kind == reflect.Float32 || kind == reflect.Float64 {
		fieldType = "numeric"
	} else if kind == reflect.Struct {
		if ft.String() == "time.Time" {
			timeMode := tagDefault(tag, "time-mode", "date")
			fieldType = timeMode
		} else {
			//-- is a struct
		}
	} else if kind == reflect.Slice {
		fieldType = "array"
	} else {
		fieldType = strings.ToLower(tf.Type.String())
	}

	fieldName := strings.ToLower(tagDefault(tag, toolkit.TagName(), tf.Name))
	fm.ID = tf.Name
	fm.Field = fieldName
	fm.GridShow = tagMemberDefault(tag, "grid-show", []string{"show", "include", "hide"}, defaultGridShow)
	fm.FormShow = tagMemberDefault(tag, "form-show", []string{"show", "hide"}, "show")
	fm.Label = tagDefault(tag, "label", tf.Name)
	fm.Align = tagDefault(tag, "align", "left")
	fm.Control = tagDefault(tag, "kf-control", "")
	if fm.Control == "" {
		fm.FieldType = fieldType
	} else {
		fm.FieldType = fm.Control
	}
	fm.MultiRow = toolkit.ToInt(tagDefault(tag, "kf-multirow", "1"), toolkit.RoundingAuto)
	fm.Format = tagDefault(tag, "format", "")
	fm.Masked = tagDefault(tag, "masked", "false") == "true"
	fm.Required = strings.ToLower(tagDefault(tag, "required", "false")) == "true"
	fm.ReadOnly = strings.ToLower(tagDefault(tag, "readonly", "false")) == "true"
	fm.GridWidth = tag.Get("grid-width")
	fm.MinLength = toolkit.ToInt(tagDefault(tag, "min-length", "0"), toolkit.RoundingAuto)
	fm.MaxLength = toolkit.ToInt(tagDefault(tag, "max-length", "0"), toolkit.RoundingAuto)
	fm.MinValue = toolkit.ToInt(tagDefault(tag, "min-value", "0"), toolkit.RoundingAuto)
	fm.MaxValue = toolkit.ToInt(tagDefault(tag, "max-value", "0"), toolkit.RoundingAuto)
	if fm.FieldType == "date" && fm.Format == "" {
		fm.Format = DefaultDateFormat
	}
	formPos := strings.Split(tagDefault(tag, "kf-pos", ","), ",")
	if len(formPos) >= 1 {
		if formPos[0] == "" {
			fm.Row = 999
		} else {
			fm.Row = toolkit.ToInt(formPos[0], toolkit.RoundingAuto)
		}
	}

	if len(formPos) >= 2 {
		fm.Col = toolkit.ToInt(formPos[1], toolkit.RoundingAuto)
	}

	//-- list
	if listItems := strings.Split(tagDefault(tag, "kf-list", ""), "|"); len(listItems) > 1 {
		fm.UseList = true
		fm.ListItems = listItems
	} else {
		fm.UseList = false
		fm.ListItems = []string{}
	}
	if useList := tagDefault(tag, "kf-use-list", ""); !fm.UseList && (useList == "1" || useList == "true") {
		fm.UseList = true
	}

	allowAdd := tagDefault(tag, "kf-allow-add", "")
	if allowAdd == "1" || allowAdd == "true" {
		fm.AllowAdd = true
	}

	mdl.fieldMetas[tf.Name] = fm
	mdl.fieldNames = append(mdl.fieldNames, tf.Name)
	//am.fieldMetas = append(am.fieldMetas, fm)

	//-- lookup
	if lookupTxt := tagDefault(tag, "kf-lookup", ""); lookupTxt != "" {
		lookupParts := strings.Split(lookupTxt, ",")
		if len(lookupParts) >= 3 {
			fm.LookupURL = lookupParts[0]
			fm.LookupKey = lookupParts[1]
			fm.LookupFields = lookupParts[2:]
		}
	}

	//-- key fields
	keyFields := mdl.KeyFields
	if isKey := tagDefault(tag, "key", "0"); isKey == "1" || strings.ToLower(isKey) == "true" {
		fieldName := tagDefault(tag, toolkit.TagName(), fm.ID)
		keyFields = append(keyFields, fieldName)
	}
	mdl.KeyFields = keyFields

	return fm
}
