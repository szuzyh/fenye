package controllers

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/statistics/models"
	"strconv"
	"strings"
)

// StatisticsController operations for Statistics
type StatisticsController struct {
	beego.Controller
}

// URLMapping ...
func (c *StatisticsController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Statistics
// @Param	body		body 	models.Statistics	true		"body for Statistics content"
// @Success 201 {int} models.Statistics
// @Failure 403 body is empty
// @router / [post]
func (c *StatisticsController) Post() {
	var v models.Statistics
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddStatistics(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Statistics by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Statistics
// @Failure 403 :id is empty
// @router /:id [get]
func (c *StatisticsController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetStatisticsById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetIpc ...
// @Title Get IpcAll 获取摄像头
// @Description get IpcAll
// Param	id		path 	string	true		"The key for ipc id, 摄像头id,获取单个摄像头(暂不开放)"
// @Success 200 {object} models.Statistics
// @Failure 403 :id is empty
// @router /ipcs
func (c *StatisticsController) GetIpcAll() {
	//idStr := c.Ctx.Input.Param(":id")
	//id, _ := strconv.Atoi(idStr)
	v, err := models.GetIpcs()
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Statistics
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	begin	query	string	false	"Start position of time. Must be an integer.TimeStamp!!"
// @Param	end	    query	string	false	"End position of time. Must be an integer.TimeStamp!!"
// @Success 200 {object} models.Statistics
// @Failure 403
// @router / [get]
func (c *StatisticsController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 5
	var offset int64
	var begin int64
	var end int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// begin: timestamp
	if v, err := c.GetInt64("begin"); err == nil {
		begin = v
	}
	// end: timestamp
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllStatistics(query, fields, sortby, order, offset, limit, begin, end)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Statistics
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Statistics	true		"body for Statistics content"
// @Success 200 {object} models.Statistics
// @Failure 403 :id is not int
// @router /:id [put]
func (c *StatisticsController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Statistics{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateStatisticsById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Statistics
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *StatisticsController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteStatistics(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
