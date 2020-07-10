package klayout

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"git.kanosolution.net/kano/kaos"
	"git.kanosolution.net/koloni/crowd"
)

type mod struct {
}

func NewMod() *mod {
	return new(mod)
}

func (m *mod) Name() string {
	return "gridform"
}

func (m *mod) MakeGlobalRoute(svc *kaos.Service) ([]*kaos.ServiceRoute, error) {
	return []*kaos.ServiceRoute{}, nil
}

func (m *mod) MakeModelRoute(svc *kaos.Service, model *kaos.ServiceModel) ([]*kaos.ServiceRoute, error) {
	routes := []*kaos.ServiceRoute{}

	alias := model.Name

	sr := new(kaos.ServiceRoute)
	sr.Path = filepath.Join(svc.BasePoint(), alias, "formconfig")
	sr.Fn = reflect.ValueOf(func(ctx *kaos.Context, name string) (interface{}, error) {
		hr := ctx.Data().Get("http-request", nil).(*http.Request)
		name = hr.URL.Query().Get("name")
		uim := getDataModel(model)
		return FormConfig(uim, name)
	})
	routes = append(routes, sr)

	sr = new(kaos.ServiceRoute)
	sr.Path = filepath.Join(svc.BasePoint(), alias, "gridconfig")
	sr.Fn = reflect.ValueOf(func(ctx *kaos.Context, name string) (interface{}, error) {
		hr := ctx.Data().Get("http-request", nil).(*http.Request)
		name = hr.URL.Query().Get("name")
		uim := getDataModel(model)
		return GridConfig(uim, name)
	})
	routes = append(routes, sr)

	return routes, nil
}

func (m *mod) MakeEvent(s *kaos.Service, model *kaos.ServiceModel, ev kaos.EventHub, _ ...string) error {
	return nil
}

func getDataModel(sm *kaos.ServiceModel) *UIModel {
	m, _ := NewUIModel(sm.Model)
	return m
}

type gridConf struct {
	Field     string   `json:"field"`
	FieldType string   `json:"fieldType"`
	Format    string   `json:"format"`
	Label     string   `json:"label"`
	Align     string   `json:"align"`
	Width     string   `json:"width"`
	Show      string   `json:"show"`
	UseList   bool     `json:"useList"`
	ListItems []string `json:"listItems"`
	Control   string   `json:"control"`
}

type GridConfigResult struct {
	SearchFields []string    `json:"searchFields"`
	KeyField     string      `json:"keyField"`
	Fields       []*gridConf `json:"fields"`
}

func GridConfig(a *UIModel, name string) (*GridConfigResult, error) {
	res := new(GridConfigResult)

	confs := []*gridConf{}
	if name == "" || name == "default" {
		for _, fn := range a.fieldNames {
			if fm, ok := a.fieldMetas[fn]; ok {
				if fm.GridShow == "show" || fm.GridShow == "include" {
					confs = append(confs, &gridConf{
						Field:     fm.Field,
						FieldType: fm.FieldType,
						Format:    fm.Format,
						Label:     fm.Label,
						Align:     fm.Align,
						Width:     fm.GridWidth,
						Show:      fm.GridShow,
						UseList:   fm.UseList,
						ListItems: fm.ListItems,
						Control:   fm.Control,
					})
				}
			}
		}
	} else {
		rv := reflect.ValueOf(a.model)
		fnk := rv.MethodByName("KGrid" + name)
		if fnk.Kind() != reflect.Func {
			return nil, fmt.Errorf("invalid Grid Config setting: "+name+" .Kind: "+fnk.Kind().String(), http.StatusInternalServerError)
		}

		outs := fnk.Call([]reflect.Value{})
		metas := outs[0].Interface().([]*FieldMeta)
		for _, fm := range metas {
			if fm.GridShow == "show" || fm.GridShow == "include" {
				confs = append(confs, &gridConf{
					Field:     fm.Field,
					FieldType: fm.FieldType,
					Format:    fm.Format,
					Label:     fm.Label,
					Align:     fm.Align,
					Width:     fm.GridWidth,
					Show:      fm.GridShow,
					UseList:   fm.UseList,
					ListItems: fm.ListItems,
					Control:   fm.Control,
				})
			}
		}
	}

	if len(a.SearchFields) == 0 {
		res.SearchFields = a.KeyFields
	} else {
		res.SearchFields = a.SearchFields
	}
	res.KeyField = a.KeyFields[0]
	res.Fields = confs
	return res, nil
}

type FormConfigResult struct {
	Title       string     `json:"title"`
	KeyFields   []string   `json:"keyfields"`
	TitleFields []string   `json:"titlefields"`
	Rows        []*FormRow `json:"rows"`
}

func FormConfig(a *UIModel, name string) (*FormConfigResult, error) {
	for _, fm := range a.fieldMetas {
		if fm.Row == 0 {
			fm.Row = 999
		}
	}

	config := new(FormConfigResult)
	if name == "" || name == "default" {
		mapRows := []*FormRow{}
		if err := crowd.FromMap(a.fieldMetas).
			Filter(func(fm *FieldMeta) bool {
				return fm.FormShow == "show"
			}).
			Group(func(fm *FieldMeta) int {
				return fm.Row
			}).
			Map(func(key int, fcols []*FieldMeta) *FormRow {
				cols, _ := crowd.FromSlice(fcols).Sort(func(f1, f2 *FieldMeta) bool {
					if f1.Col != f2.Col {
						return f1.Col < f2.Col
					} else {
						return strings.Compare(f1.Label, f2.Label) < 0
					}
				}).Collect().Exec()
				return &FormRow{key, cols.([]*FieldMeta)}
			}).
			Sort(func(f1, f2 *FormRow) bool {
				return f1.RowIndex < f2.RowIndex
			}).
			Collect().
			Run(&mapRows); err != nil {
			return nil, err
		}

		config.Rows = mapRows
	}

	//if len(a.FormNameFields) > 0 {
	//config.KeyFields = []string{a.FormNameFields[0]}
	//}
	config.KeyFields = a.KeyFields
	config.Title = a.FormTitle
	config.TitleFields = a.FormNameFields
	return config, nil
}
