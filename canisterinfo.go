package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CanisterInfo struct {
	Repository   string `json"repo"`
	Controller   string `json:"controller"`
	CanisterName string `json:"name"`
	CanisterID   string `json:"id"`
	CanisterType string `json:"type"`
	Framework    string `json:"framework"`
	Network      string `json:"network"`
}

func handleCanisterListAPI(r *gin.Engine) {
	r.GET("public/canisterlist", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, TRACE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type,Token,Accept, Connection, User-Agent, Cookie")
		c.Header("Access-Control-Max-Age", "3628800")

		ret, err := readCanisterList(context.TODO())
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		res, err := json.Marshal(ret)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "",
			"result":  string(res),
		})
	})
}

func handleCanisterInfoAPI(r *gin.Engine) {
	r.GET("public/canisterinfo", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, TRACE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type,Token,Accept, Connection, User-Agent, Cookie")
		c.Header("Access-Control-Max-Age", "3628800")

		canisterid := c.Query("canisterid")
		if canisterid == "" {
			c.String(http.StatusBadRequest, "canister id must provide")
			return
		}

		info, err := ReadCanisterInfo(context.TODO(), canisterid)
		if err != nil {
			c.String(http.StatusBadRequest, "read canister info err: %s", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "",
			"data":    string(info),
		})
	})
}