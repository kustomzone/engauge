package local

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/JKhawaja/errors"
)

func appendCSV(filename string, line []string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.New(err, map[string]interface{}{
			"filename": filename,
		})
	}

	w := csv.NewWriter(f)
	err = w.Write(line)
	if err != nil {
		return errors.New(err, map[string]interface{}{
			"filename": filename,
		})
	}

	w.Flush()

	return f.Close()
}

func (c *Client) decodeList(resource string, filenames []string) (interface{}, error) {
	switch resource {
	case db.Endpoints:
		list := make([]*types.Endpoint, 0)
		for _, filename := range filenames {
			fullName := fmt.Sprintf("%s/%s/%s", c.basepath, resource, filename)
			data, err := ioutil.ReadFile(fullName)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			item, err := decodeFile(resource, data)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			list = append(list, item.(*types.Endpoint))
		}
		return list, nil
	case db.Origins:
		list := make([]*types.Origin, 0)
		for _, filename := range filenames {
			fullName := fmt.Sprintf("%s/%s/%s", c.basepath, resource, filename)
			data, err := ioutil.ReadFile(fullName)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			item, err := decodeFile(resource, data)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			list = append(list, item.(*types.Origin))
		}
		return list, nil
	case db.OriginStats, db.EndpointStats, db.EntityStats:
		list := make([]*types.IntervalStats, 0)
		for _, filename := range filenames {
			fullName := fmt.Sprintf("%s/%s/%s", c.basepath, resource, filename)
			data, err := ioutil.ReadFile(fullName)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			item, err := decodeFile(resource, data)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			list = append(list, item.(*types.IntervalStats))
		}
		return list, nil
	case db.Entities:
		list := make([]*types.Entity, 0)
		for _, filename := range filenames {
			fullName := fmt.Sprintf("%s/%s/%s", c.basepath, resource, filename)
			data, err := ioutil.ReadFile(fullName)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			item, err := decodeFile(resource, data)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			list = append(list, item.(*types.Entity))
		}
		return list, nil
	case db.Properties:
		list := make([]*types.Property, 0)
		for _, filename := range filenames {
			fullName := fmt.Sprintf("%s/%s/%s", c.basepath, resource, filename)
			data, err := ioutil.ReadFile(fullName)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			item, err := decodeFile(resource, data)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			list = append(list, item.(*types.Property))
		}
		return list, nil
	case db.PropertyStats:
		list := make([]*types.PropertyStats, 0)
		for _, filename := range filenames {
			fullName := fmt.Sprintf("%s/%s/%s", c.basepath, resource, filename)
			data, err := ioutil.ReadFile(fullName)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			item, err := decodeFile(resource, data)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			list = append(list, item.(*types.PropertyStats))
		}
		return list, nil
	case db.Summaries:
		list := make([]*types.Summary, 0)
		for _, filename := range filenames {
			fullName := fmt.Sprintf("%s/%s/%s", c.basepath, resource, filename)
			data, err := ioutil.ReadFile(fullName)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			item, err := decodeFile(resource, data)
			if err != nil {
				return nil, errors.New(err, map[string]interface{}{
					"resource": resource,
					"file":     filename,
				})
			}

			list = append(list, item.(*types.Summary))
		}
		return list, nil
	}

	return nil, nil
}

func decodeFile(resource string, b []byte) (interface{}, error) {
	var item interface{}
	switch resource {
	case db.Interactions:
		item = &types.Interaction{}
	case db.Endpoints:
		item = &types.Endpoint{}
	case db.Origins:
		item = &types.Origin{}
	case db.Entities:
		item = &types.Entity{}
	case db.EndpointStats, db.OriginStats, db.EntityStats:
		item = &types.IntervalStats{}
	case db.Properties:
		item = &types.Property{}
	case db.PropertyStats:
		item = &types.PropertyStats{}
	case db.Summaries:
		item = &types.Summary{}
	case db.Settings:
		item = &types.Settings{}
	}

	// decode
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err := dec.Decode(item)
	if err != nil {
		return nil, errors.New(err, map[string]interface{}{
			"resource": resource,
		})
	}

	return item, nil
}

func encode(item interface{}) ([]byte, error) {
	// encode
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(item)
	if err != nil {
		return []byte{}, errors.New(err, nil)
	}

	return buf.Bytes(), nil
}

func pint(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func pint64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
