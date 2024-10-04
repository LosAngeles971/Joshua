package core

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Data struct {
	Data map[string]interface{}
}

func newData() *Data {
	return &Data{
		Data: map[string]interface{}{},
	}
}

func (d *Data) Clone() *Data {
	c := newData()
	for key, value := range d.Data {
		switch t := value.(type) {
		case int:
			c.Data[key] = d.Data[key].(int)
		case int64:
			c.Data[key] = d.Data[key].(int64)
		case float64:
			c.Data[key] = d.Data[key].(float64)
		case string:
			c.Data[key] = d.Data[key].(string)
		case bool:
			c.Data[key] = d.Data[key].(bool)
		case []string:
			c.Data[key] = d.Data[key].([]string)
		default:
			logrus.Warningf("unsupported data type ( %v )", t)
		}
	}
	return c
}

func (d *Data) SetData(key string, value interface{}) error {
	switch t := value.(type) {
	case int, int64, float64, string, bool, []string:
		d.Data[key] = value
	default:
		return fmt.Errorf("data not supported for key/value store ( %v  )", t)
	}
	return nil
}

func (d *Data) SetDatas(data map[string]interface{}) error {
	for k, v := range data {
		err := d.SetData(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Data) GetData(key string) (interface{}, bool) {
	if v, ok := d.Data[key]; ok {
		return v, true
	} else {
		return nil, false
	}
}

func (d *Data) GetString(key string) (string, bool) {
	if v, ok := d.GetData(key); ok {
		return v.(string), true
	} else {
		return "", false
	}
}

func (d *Data) GetInt(key string) (int, bool) {
	if v, ok := d.GetData(key); ok {
		return v.(int), true
	} else {
		return 0, false
	}
}

func (d *Data) GetInt64(key string) (int64, bool) {
	if v, ok := d.GetData(key); ok {
		return v.(int64), true
	} else {
		return 0, false
	}
}

func (d *Data) GetFloat64(key string) (float64, bool) {
	if v, ok := d.GetData(key); ok {
		return v.(float64), true
	} else {
		return 0, false
	}
}

func (d *Data) GetBool(key string) (bool, bool) {
	if v, ok := d.GetData(key); ok {
		return v.(bool), true
	} else {
		return false, false
	}
}

func (d *Data) GetStringArray(key string) ([]string, bool) {
	if v, ok := d.GetData(key); ok {
		return v.([]string), true
	} else {
		return nil, false
	}
}

func (d1 *Data) IsEqual(d2 *Data) bool {
	for k, v := range d1.Data {
		if vv, ok := d2.Data[k]; ok {
			t1 := reflect.TypeOf(v)
			t2 := reflect.TypeOf(vv)
			if t1 != t2 {
				return false
			}
			switch t1.Kind() {
			case reflect.Int:
				if v.(int) != vv.(int) {
					return false
				}
			case reflect.Int64:
				if v.(int64) != vv.(int64) {
					return false
				}
			case reflect.Float64:
				if v.(float64) != vv.(float64) {
					return false
				}
			case reflect.String:
				if v.(string) != vv.(string) {
					return false
				}
			case reflect.Bool:
				if v.(bool) != vv.(bool) {
					return false
				}
			}			
		} else {
			return false
		}
	}
	return true
}