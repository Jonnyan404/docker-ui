package main

import (
	"expvar"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/gohutool/boot4go-docker-ui/db"
	"github.com/gohutool/boot4go-docker-ui/handle"
	. "github.com/gohutool/boot4go-docker-ui/log"
	constants "github.com/gohutool/boot4go-docker-ui/model"
	prometheusfasthttp "github.com/gohutool/boot4go-prometheus/fasthttp"
	httputil "github.com/gohutool/boot4go-util/http"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
	"github.com/valyala/fasthttp/pprofhandler"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : docker.ui.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/12 20:19
* 修改历史 : 1. [2022/5/12 20:19] 创建文件 by LongYong
*/

var SERVER_VERSION = "docker-ui-dev"

const (
	SERVER_MAJOR = 1
	SERVER_MINOR = 0
	SERVER_BUILD = 0
)

func main() {
	app := kingpin.New("docker-ui", "A docker management with UI.")
	addr_flag := app.Flag("addr", "Addr: docker management listen addr.").Short('l').Default(":8999").String()
	issuer_flag := app.Flag("issuer", "Issuer: token's issuer.").Short('i').Default(constants.DEFAULT_ISSUER).String()
	expired_flag := app.Flag("token_expire", "Token_expire: many hour(s) token will expire.").Short('e').Default("24").Int()

	reset_user_flag := app.Flag("reset-user", "Reset password for username (creates user if not exists), then exit unless --reset-keep-running is set.").Default("").String()
	reset_password_flag := app.Flag("reset-password", "Password used with --reset-user.").Default("").String()
	reset_keep_running_flag := app.Flag("reset-keep-running", "Keep running the server after reset/create user.").Default("false").Bool()

	li_flag := app.Flag("license", "License: CubeUI License.").Default("").String()
	endpoint_flag := app.Flag("endpoint", "Endpoint: the endpoint of docker, default is unix, if tcp is open, like as 192.168.56.102:2375").Default("unix").String()

	app.HelpFlag.Short('h')
	app.Version(SERVER_VERSION)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	setLicenseJS(*li_flag)
	setEndpointJS(*endpoint_flag)

	db.InitDB()

	Logger.Info("Database is load.")

	if reset_user_flag != nil && reset_password_flag != nil && len(*reset_user_flag) > 0 && len(*reset_password_flag) > 0 {
		userID := db.UserIDFromUsername(*reset_user_flag)
		u := db.GetUser(userID)
		if u == nil {
			if err := db.CreateUser(*reset_user_flag, *reset_password_flag); err != nil {
				Logger.Critical("Create user failed: %v", err)
				return
			}
			Logger.Info("User created: %v", *reset_user_flag)
		} else {
			if err := db.UpdatePwd(userID, *reset_password_flag); err != nil {
				Logger.Critical("Reset password failed: %v", err)
				return
			}
			Logger.Info("Password reset ok: %v", *reset_user_flag)
		}

		if reset_keep_running_flag == nil || !*reset_keep_running_flag {
			return
		}
	}

	l, err := net.Listen("tcp", *addr_flag)
	if err != nil {
		fmt.Println("Start server error " + err.Error())
		return
	}

	if issuer_flag != nil && len(*issuer_flag) > 0 {
		constants.Issuer = *issuer_flag
	}

	if expired_flag != nil && *expired_flag > 0 {
		constants.TokenExpire = time.Duration(*expired_flag) * time.Hour
	}

	if err = initAdmin(); err != nil {
		panic("Init Admin User error " + err.Error())
	}

	Logger.Info("Start " + SERVER_VERSION + " now .... ")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	Logger.Debug("%v %v %v %v", *addr_flag, *issuer_flag, *expired_flag, *li_flag)

	go func(twg *sync.WaitGroup) {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		fmt.Println("signal service close")
		twg.Done()
	}(wg)

	startHttpServer(l)

	wg.Wait()

	l.Close()
	Logger.Info(SERVER_VERSION + " is close")
	time.Sleep(20 * time.Microsecond)
}

func initAdmin() error {
	return nil
}

func startHttpServer(listener net.Listener) {

	router := routing.New()

	v3Group := router.Group("/v3/api")

	handle.PrometheusHandler.InitRouter(router, v3Group)
	handle.UserHandler.InitRouter(router, v3Group)

	fsHandler := newStaticHandler()

	router.Get("/stats", func(ctx *routing.Context) error {
		expvarhandler.ExpvarHandler(ctx.RequestCtx)
		return nil
	})

	var POOLSIZE int = 16
	var testPool [][]byte = make([][]byte, POOLSIZE)

	mu := sync.Mutex{}
	var SIZE int = 1024 * 1024 * 12

	router.Get("/usages", func(ctx *routing.Context) error {
		v, _ := mem.VirtualMemory()
		percent, _ := cpu.Percent(time.Second, false)
		info2, _ := disk.Usage("/") //指定某路径的硬盘使用情况

		data := make(map[string]any)
		data["cpu"] = percent
		data["disk"] = info2
		data["memory"] = v

		httputil.Result.Success(data, "ok").Response(ctx.RequestCtx)

		mu.Lock()
		defer mu.Unlock()

		if len(testPool) >= POOLSIZE {
			testPool = testPool[0:0]
		}

		testPool = append(testPool, make([]byte, SIZE))

		return nil
	})

	router.Any("/debug/pprof/*", func(ctx *routing.Context) error {
		pprofhandler.PprofHandler(ctx.RequestCtx)
		return nil
	})

	//wsHanlder := &EchoWSHandler{}

	wsGroup := router.Group("/docker-api-ws")
	handle.InitWsRouter(router, wsGroup)
	//router.Any("/docker-api-ws/echo", EchoWSHandler)

	router.Any("/docker-api-stream/*", func(ctx *routing.Context) error {
		ri, err := handle.GetRequestInfo("/docker-api-stream", ctx)
		if err != nil {
			return err
		}

		ctx.Response.Header.Set(fasthttp.HeaderTransferEncoding, "chunked")

		if strings.Index(ri.Host, "unix") == 0 {
			if err := httputil.WithHttpReverseUnixSocketProxy("/var/run/docker.sock", ri.Path, ctx.QueryArgs().String(), ctx.RequestCtx); err != nil {
				Logger.Debug("Proxy stream error : %v", err)
			} else {
				return err
			}
		} else {
			if err := httputil.WithHttpReverseProxy(ri.Host, "http", ri.Path, ctx.QueryArgs().String(), ctx.RequestCtx); err != nil {
				Logger.Debug("Proxy stream error : %v", err)
			} else {
				return err
			}
		}

		return nil
	})

	router.Any("/docker-api/*", func(ctx *routing.Context) error {
		ri, err := handle.GetRequestInfo("/docker-api", ctx)
		if err != nil {
			return err
		}

		if strings.Index(ri.Host, "unix") == 0 {

			if err := httputil.UnixSocketProxy("/var/run/docker.sock", ri.Host, "http", ri.Path, "", ctx.RequestCtx, 0); err != nil {
				Logger.Debug("Proxy error : %v", err)
			} else {
				Logger.Debug("Proxy over : %v", string(ctx.RequestURI()))
				return err
			}

		} else {

			if err := httputil.Proxy(ri.Host, "http", ri.Path, "", ctx.RequestCtx, 0); err != nil {
				Logger.Debug("Proxy error : %v", err)
			} else {
				Logger.Debug("Proxy over : %v", string(ctx.RequestURI()))
				return err
			}

		}

		return nil
	})

	router.Any("/*", func(context *routing.Context) error {
		ctx := context.RequestCtx
		fsHandler(ctx)
		UpdateFSCounters(ctx)
		return nil
	})

	requestHandler := func(ctx *fasthttp.RequestCtx) {

		Logger.Debug("%v %v %v %v", string(ctx.Path()), ctx.URI().String(), string(ctx.Method()), ctx.QueryArgs().String())
		defer func() {
			if err := recover(); err != nil {
				Logger.Debug(err)
				// ctx.Error(fmt.Sprintf("%v", err), http.StatusInternalServerError)
				httputil.Error(ctx, httputil.Result.Fail(fmt.Sprintf("%v", err)).Json(), http.StatusInternalServerError)
			}

			ctx.Response.Header.Set("tick", time.Now().String())
			ctx.Response.Header.SetServer("Docker-UIManager")

			prometheusfasthttp.RequestCounterHandler(nil)(ctx)

			Logger.Debug("router.HandleRequest is finish")

		}()

		router.HandleRequest(ctx)
	}

	//go http.ListenAndServe("0.0.0.0:8887", nil)

	// Start HTTP server.
	Logger.Info("Starting HTTP server on %v", listener.Addr().String())
	go func() {

		s := &fasthttp.Server{
			Handler: requestHandler,
		}

		//fix bug error when serving connection "[::1]:8999"<->"[::1]:27982": body size exceeds the given limit
		maxBodySize := int64(4) * 1024 * 1024 * 1024
		maxInt := int64(^uint(0) >> 1)
		if maxBodySize > maxInt {
			maxBodySize = maxInt
		}
		s.MaxRequestBodySize = int(maxBodySize)

		if err := s.Serve(listener); err != nil {
			Logger.Critical("error in ListenAndServe: %v", err)
		}
	}()
}

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)

func UpdateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

func GetUserId(context *routing.Context) string {
	return context.Get("userid").(string)
}

func SetUserId(context *routing.Context, userid string) {
	context.Set("userid", userid)
}
