/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package busi

import (
	"fmt"
	"time"

	"github.com/dtm-labs/dtm/dtmcli"
	"github.com/dtm-labs/dtm/dtmcli/logger"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

// 启动命令：go run app/main.go qs

// 事务参与者的服务地址
const qsBusiAPI = "/api/busi_start"
const qsBusiPort = 8082

var qsBusi = fmt.Sprintf("http://localhost:%d%s", qsBusiPort, qsBusiAPI)

// QsStartSvr quick start: start server
func QsStartSvr() {
	app := dtmutil.GetGinApp()
	qsAddRoute(app)
	logger.Infof("quick start examples listening at %d", qsBusiPort)
	go func() {
		_ = app.Run(fmt.Sprintf(":%d", qsBusiPort))
	}()
	time.Sleep(100 * time.Millisecond)
}

// QsFireRequest quick start: fire request
func QsFireRequest() string {
	req := &gin.H{"amount": 30} // 微服务的载荷
	// DtmServer为DTM服务的地址
	saga := dtmcli.NewSaga(dtmutil.DefaultHTTPServer, dtmcli.MustGenGid(dtmutil.DefaultHTTPServer)).
		// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransOutCompensate"
		Add(qsBusi+"/TransOut", qsBusi+"/TransOutCompensate", req).
		// 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransInCompensate"
		Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// 等待事务全部完成后再返回，可选
	saga.WaitResult = true
	// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	err := saga.Submit()
	logger.FatalIfError(err)
	return saga.Gid
}

func qsAddRoute(app *gin.Engine) {
	app.POST(qsBusiAPI+"/TransIn", dtmutil.WrapHandler2(func(c *gin.Context) interface{} {
		logger.Infof("TransIn")
		return nil
	}))
	app.POST(qsBusiAPI+"/TransInCompensate", dtmutil.WrapHandler2(func(c *gin.Context) interface{} {
		logger.Infof("TransInCompensate")
		return nil
	}))
	app.POST(qsBusiAPI+"/TransOut", dtmutil.WrapHandler2(func(c *gin.Context) interface{} {
		logger.Infof("TransOut")
		return nil
	}))
	app.POST(qsBusiAPI+"/TransOutCompensate", dtmutil.WrapHandler2(func(c *gin.Context) interface{} {
		logger.Infof("TransOutCompensate")
		return nil
	}))
}
