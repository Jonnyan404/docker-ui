package log

import (
	"os"
	"path/filepath"

	"github.com/gohutool/log4go"
)

const defaultLog4goPath = "./config/log4go.xml"

const defaultLog4goXML = `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
	<appender enabled="true" name="console">
		<type>console</type>
		<pattern>[%D %T %m] [%L][%l] (%S) %M</pattern>
		<!-- level is (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) -->
	</appender>
	<appender enabled="true" name="file">
		<type>file</type>
		<pattern>[%D %T %m] [%L][%l] (%S) %M</pattern>
		<property name="filename">config/docker.ui.log</property>
		<property name="rotate">false</property>
		<property name="maxsize">0M</property>
		<property name="maxlines">0K</property>
		<property name="daily">true</property>
	</appender>

	<root>
		<level>info</level>
		<appender-ref ref="console" />
		<appender-ref ref="file" />
	</root>

</configuration>
`

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : logger.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/12 20:20
* 修改历史 : 1. [2022/5/12 20:20] 创建文件 by LongYong
*/

func init() {
	defer func() {
		if err := recover(); err != nil {
			log4go.LoggerManager.InitWithDefaultConfig()
		}
	}()

	if cfg := findOrCreateLog4goConfig(); cfg != "" {
		log4go.LoggerManager.InitWithXML(cfg)
		return
	}

	log4go.LoggerManager.InitWithDefaultConfig()

}

func findOrCreateLog4goConfig() string {
	// 0) explicit config path via env
	if v := os.Getenv("LOG4GO_XML"); v != "" {
		if _, err := os.Stat(v); err == nil {
			return v
		}
	}

	// 1) prefer ./config/log4go.xml (keep it together with ./config/data.db)
	if ensureDefaultConfigAt(defaultLog4goPath) {
		return defaultLog4goPath
	}

	// 1) current working directory
	if _, err := os.Stat("log4go.xml"); err == nil {
		return "log4go.xml"
	}

	// 2) executable directory
	exe, err := os.Executable()
	if err == nil {
		cfg := filepath.Join(filepath.Dir(exe), "log4go.xml")
		if _, err := os.Stat(cfg); err == nil {
			return cfg
		}
	}

	// 3) create a default config (best-effort): current directory.
	if writeDefaultIfPossible("log4go.xml") {
		return "log4go.xml"
	}

	// 4) if current directory is not writable, try executable directory.
	if err == nil {
		cfg := filepath.Join(filepath.Dir(exe), "log4go.xml")
		if writeDefaultIfPossible(cfg) {
			return cfg
		}
	}

	return ""
}

func ensureDefaultConfigAt(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false
	}

	return writeDefaultIfPossible(path)
}

func writeDefaultIfPossible(path string) bool {
	// Don't overwrite an existing file.
	if _, err := os.Stat(path); err == nil {
		return true
	}
	// Best-effort write.
	if err := os.WriteFile(path, []byte(defaultLog4goXML), 0644); err != nil {
		return false
	}
	return true
}

var Logger = log4go.LoggerManager.GetLogger("gohutool.docker4go.ui")
