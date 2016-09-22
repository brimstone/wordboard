package main

import (
	"encoding/json"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// merge merges the two JSON-marshalable values x1 and x2,
// preferring x1 over x2 except where x1 and x2 are
// JSON objects, in which case the keys from both objects
// are included and their values merged recursively.
//
// It returns an error if x1 or x2 cannot be JSON-marshaled.
// https://groups.google.com/forum/#!topic/golang-nuts/nLCy75zMlS8
func merge(x1, x2 interface{}) (interface{}, error) {
	data1, err := json.Marshal(x1)
	if err != nil {
		return nil, err
	}
	data2, err := json.Marshal(x2)
	if err != nil {
		return nil, err
	}
	var j1 interface{}
	err = json.Unmarshal(data1, &j1)
	if err != nil {
		return nil, err
	}
	var j2 interface{}
	err = json.Unmarshal(data2, &j2)
	if err != nil {
		return nil, err
	}
	return merge1(j1, j2), nil
}

func merge1(x1, x2 interface{}) interface{} {
	switch x1 := x1.(type) {
	case map[string]interface{}:
		x2, ok := x2.(map[string]interface{})
		if !ok {
			return x1
		}
		for k, v2 := range x2 {
			if v1, ok := x1[k]; ok {
				x1[k] = merge1(v1, v2)
			} else {
				x1[k] = v2
			}
		}
	case nil:
		// merge(nil, map[string]interface{...}) -> map[string]interface{...}
		x2, ok := x2.(map[string]interface{})
		if ok {
			return x2
		}
	}
	return x2
}

func main() {
	Mux := gin.Default()

	var data interface{}

	// POST is an update
	Mux.POST("/", func(c *gin.Context) {
		var body interface{}
		err := c.BindJSON(&body)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		data, err = merge(data, body)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.JSON(200, data)
	})
	// PUT is an overwrite
	Mux.PUT("/", func(c *gin.Context) {
		err := c.BindJSON(&data)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.JSON(200, data)
	})

	// Get does things
	Mux.GET("/data.json", func(c *gin.Context) {
		c.JSON(200, data)
	})

	Mux.StaticFile("/", "static/index.html")

	Mux.NoRoute(func(c *gin.Context) {
		static.ServeRoot("/", "static")(c)
	})
	Mux.Run(":8080")
}
