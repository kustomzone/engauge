package local

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/JKhawaja/errors"
)

// Client --
type Client struct {
	basepath string
}

// NewClient --
func NewClient(basepath string) (*Client, error) {
	c := &Client{
		basepath: basepath,
	}
	err := c.init()
	if err != nil {
		return nil, errors.New(err, nil)
	}

	c.InitCache()

	return c, nil
}

func (c *Client) init() error {
	err := os.MkdirAll(c.basepath, 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.Interactions), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.Endpoints), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.EndpointStats), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.Origins), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.OriginStats), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.Entities), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.EntityStats), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.Properties), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.PropertyStats), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s/", c.basepath, db.Summaries), 0644)
	if err != nil {
		return errors.New(err, nil)
	}

	return nil
}

// Do --
func (c *Client) Do(op *db.Op) db.Result {
	var result db.Result

	switch op.Type {
	case db.Create, db.Update:
		if op.Resource == db.Interactions {
			interaction := op.Item.(*types.Interaction)
			filename := fmt.Sprintf("%s/%s/%s.csv", c.basepath, db.Interactions, interaction.Date())
			appendCSV(filename, interaction.CSV())
			return result
		}

		filename := c.filename(op.Resource, op.Item)

		data, err := encode(op.Item)
		if err != nil {
			result.Error = errors.New(err, map[string]interface{}{
				"op":       op.Type,
				"resource": op.Resource,
				"file":     filename,
			})
			return result
		}

		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			result.Error = errors.New(err, map[string]interface{}{
				"op":       op.Type,
				"resource": op.Resource,
				"file":     filename,
			})
			return result
		}
	case db.Read:
		filename := c.filenameFromWhere(op.Resource, op.Where)
		data, err := ioutil.ReadFile(filename)
		if os.IsNotExist(err) {
			result.Error = types.ErrDNE
			return result
		} else if err != nil {
			result.Error = errors.New(err, map[string]interface{}{
				"op":       op.Type,
				"resource": op.Resource,
				"file":     filename,
			})
			return result
		}

		item, err := decodeFile(op.Resource, data)
		if err != nil {
			result.Error = errors.New(err, map[string]interface{}{
				"op":       op.Type,
				"resource": op.Resource,
				"file":     filename,
			})
			return result
		}

		result.Item = item
		return result
	case db.List:
		dir := fmt.Sprintf("%s/%s/", c.basepath, op.Resource)
		filenames, err := c.readDir(dir)
		if err != nil {
			result.Error = errors.New(err, map[string]interface{}{
				"op":       op.Type,
				"resource": op.Resource,
				"dir":      dir,
			})
			return result
		}

		limit := pint64(op.Limit)
		offset := pint64(op.Offset)
		if offset > 0 || limit > 0 {
			if limit == 0 {
				limit = int64(len(filenames)) - offset - 1
			}

			end := limit + offset
			if end > int64(len(filenames)-1) {
				end = int64(len(filenames) - 1)
			}

			if len(filenames) == 0 {
				end = 0
			}

			filenames = filenames[offset:end]
		}

		list, err := c.decodeList(op.Resource, filenames)
		if err != nil {
			result.Error = errors.New(err, map[string]interface{}{
				"op":       op.Type,
				"resource": op.Resource,
				"dir":      dir,
			})
			return result
		}
		result.Item = list
		return result
	case db.Count:
		dir := fmt.Sprintf("%s/%s/", c.basepath, op.Resource)

		filenames, err := c.readDir(dir)
		if err != nil {
			result.Error = errors.New(err, map[string]interface{}{
				"op":       op.Type,
				"resource": op.Resource,
				"dir":      dir,
			})
			return result
		}

		result.Item = int64(len(filenames))
		return result
	case db.Delete:
		filename := c.filenameFromWhere(op.Resource, op.Where)
		err := os.Remove(filename)
		if err != nil {
			result.Error = errors.New(err, nil)
		}
	}

	return result
}

func (c *Client) filenameFromWhere(resource string, where db.Where) string {
	if where == nil {
		return ""
	}

	wm := where.(db.WhereMap)
	i, ok := wm["item.id"]
	if !ok {
		switch resource {
		case db.Summaries:
			st, ok := wm["item.spanType"]
			if !ok {
				return ""
			}
			spanType, ok := st.(string)
			if !ok {
				return ""
			}

			return fmt.Sprintf("%s/%s/%s", c.basepath, resource, spanType)
		case db.Properties:
			n, ok := wm["item.name"]
			if !ok {
				return ""
			}
			name, ok := n.(string)
			if !ok {
				return ""
			}

			return fmt.Sprintf("%s/%s/%s", c.basepath, resource, name)
		case db.PropertyStats:
			n, ok := wm["item.name"]
			if !ok {
				return ""
			}
			name, ok := n.(string)
			if !ok {
				return ""
			}

			st, ok := wm["item.spanType"]
			if !ok {
				return ""
			}
			spanType, ok := st.(string)
			if !ok {
				return ""
			}

			return fmt.Sprintf("%s/%s/%s-%s", c.basepath, resource, name, spanType)
		}
	}

	id, ok := i.(*types.UUID)
	if !ok {
		return ""
	}

	switch resource {
	case db.EndpointStats, db.OriginStats, db.EntityStats:
		st, ok := wm["item.interval"]
		if !ok {
			return ""
		}

		interval, ok := st.(string)
		if !ok {
			return ""
		}

		return fmt.Sprintf("%s/%s/%s-%s", c.basepath, resource, id, interval)
	}

	return fmt.Sprintf("%s/%s/%s", c.basepath, resource, id)
}

// lists and counts
func (c *Client) readDir(dir string) ([]string, error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (c *Client) filename(resource string, item interface{}) string {
	var filename string
	switch resource {
	case db.Endpoints:
		i := item.(*types.Endpoint)
		filename = fmt.Sprintf("%s/%s/%s", c.basepath, resource, i.ID.String())
	case db.Origins:
		i := item.(*types.Origin)
		filename = fmt.Sprintf("%s/%s/%s", c.basepath, resource, i.ID.String())
	case db.Entities:
		i := item.(*types.Entity)
		filename = fmt.Sprintf("%s/%s/%s", c.basepath, resource, i.ID.String())
	case db.OriginStats, db.EntityStats, db.EndpointStats:
		i := item.(*types.IntervalStats)
		filename = fmt.Sprintf("%s/%s/%s-%s", c.basepath, resource, i.ID.String(), i.Interval)
	case db.Properties:
		i := item.(*types.Property)
		filename = fmt.Sprintf("%s/%s/%s", c.basepath, resource, i.Name)
	case db.PropertyStats:
		i := item.(*types.PropertyStats)
		filename = fmt.Sprintf("%s/%s/%s-%s", c.basepath, resource, i.Name, i.SpanType)
	case db.Summaries:
		i := item.(*types.Summary)
		filename = fmt.Sprintf("%s/%s/%s", c.basepath, resource, i.SpanType)
	case db.Settings:
		filename = fmt.Sprintf("%s/%s", c.basepath, resource)
	}

	return filename
}
