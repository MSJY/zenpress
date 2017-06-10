package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"github.com/insionng/zenpress/helper"
	"github.com/insionng/zenpress/module/core"
	"github.com/insionng/zenpress/module/hook"

	"github.com/insionng/makross"
	"github.com/insionng/makross/logger"
	"github.com/insionng/makross/pongor"
	"github.com/insionng/makross/static"

	"github.com/fsnotify/fsnotify"

	_ "qlang.io/lib/builtin"
)

var (
	wg   sync.WaitGroup
	port int = 8888
)

func init() {
	core.VmString(`
		hook.AddAction("init", fn {
			println("Q>init")
		})
	`)

	core.VmString(`
		hook.DoAction("init")
	`)
}

func hello() {
	fmt.Println("[hello main action]")
}

func coreHello() {
	fmt.Println("[core hello action]")
}

func bootstrap() {
	fmt.Println("[bootstrap action]")
}

func theloop() {
	fmt.Println("[theloop action]")
}

func quit() {
	fmt.Println("[quit action]")
}

func TheTitle() []byte {
	quote := "The bird is the word."
	return []byte(quote)
}

func changeTheQuote(quote []byte) []byte {
	quote = []byte(strings.Replace(string(quote), "bird", "nerd", -1))
	fmt.Println(string(quote))
	return quote
}

func bigTitle(quote []byte) []byte {
	quote = []byte("<<<" + string(quote) + ">>>")
	fmt.Println(string(quote))
	fmt.Println("CurrentFilter>", hook.CurrentFilter())
	return quote
}

func main() {
	defer func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			core.AddFunc("quit", quit)
			core.VmString(`
			hook.AddAction("quit", quit)
			hook.DoAction("quit")
`)
		}()
		wg.Wait()
	}()

	core.VmString(`
	hook.AddAction("bootstrap", bootstrap)
	hook.DoAction("bootstrap")
`)
	//////////////////////////////

	core.Plugins()

	/////////////////////////////////////////////
	core.VmString(`
		hook.AddAction("bootstrap", fn {
			println("Q>bootstrap")
		})
	`)

	core.AddFunc("theloop", theloop)

	core.VmString(`
		hook.AddAction("theloop", theloop)
	`)
	hook.DoAction("theloop")

	hook.AddFilter("peter_griffin_quote", changeTheQuote)
	hook.AddFilter("peter_griffin_quote", bigTitle)
	hook.DoFilter("peter_griffin_quote", TheTitle)

	///-------------------------------------------------///

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	//执行主题逻辑
	var theme = "default"
	appCode := core.Codes(theme)

	app := makross.New()
	app.Use(logger.Logger())
	app.Use(static.Static(fmt.Sprintf("content/theme/%s/public", theme)))
	app.SetRenderer(pongor.Renderor(pongor.Option{Directory: fmt.Sprintf("content/theme/%s/template", theme), Reload: false, Filter: true}))
	core.Q.SetVar("app", app)

	err = core.VmByte(appCode)
	if err != nil {
		fmt.Println("#1 执行Q出错：")
		panic(err)
	}

	m, okay := core.Q.GetVar("app")
	if !okay {
		panic("cannot find app object in code.")
	}

	app, okay = m.(*makross.Makross)

	/*------------------------------------*/
	/*
		if macaron.Env == macaron.DEV {
			m.Use(toolbox.Toolboxer(m, toolbox.Options{
				HealthCheckFuncs: []*toolbox.HealthCheckFuncDesc{
					&toolbox.HealthCheckFuncDesc{
						Desc: "Database connection",
						Func: models.Ping,
					},
				},
			}))
		}
	*/
	/*------------------------------------*/

	if okay {
		fmt.Printf("app.Listen(%v)\n", port)
		go app.Listen(port)
	} else {
		panic("cannot convert app to (*makross.Makross)")
	}

	// Process events
	go func() {
		for {
			select {
			case <-watcher.Events:
				// Loading App Logic
				fmt.Println("Reload Application")

				app.Close()

				appCode = core.Codes(theme)

				app = makross.New()
				app.Use(logger.Logger())
				app.Use(static.Static(fmt.Sprintf("content/theme/%s/public", theme)))
				app.SetRenderer(pongor.Renderor(pongor.Option{Directory: fmt.Sprintf("content/theme/%s/template", theme), Reload: false, Filter: true}))

				core.Q.SetVar("app", app)

				core.VmByte(appCode)

				m, okay = core.Q.GetVar("app")
				if !okay {
					panic("cannot find app object in code..")
				}

				app, okay = m.(*makross.Makross)
				if okay {
					fmt.Printf("app.Listen(%v)\n", port)
					go app.Listen(port)
				} else {
					panic("cannot convert app to (*makross.Makross)")
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	fmt.Println(".........................................................")
	fmt.Printf("Current application pid is %d\n", os.Getpid())
	fmt.Println(".........................................................")

	//热更新监控目录
	var applicationDir = "content/application"
	var applicationRootDir = fmt.Sprintf("%s/root", applicationDir)

	var themeAppDir = fmt.Sprintf("content/theme/%s/handler", theme)
	var themeAppRootDir = fmt.Sprintf("%s/root", themeAppDir)

	err = watcher.Add(applicationDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("watcher.Add applicationDir has error:%v", err))
	}

	err = watcher.Add(applicationRootDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("watcher.Add applicationRootDir has error:%v", err))
	}

	err = watcher.Add(themeAppDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("watcher.Add themeAppDir has error:%v", err))
	}

	if helper.IsExist(themeAppRootDir) {
		err = watcher.Add(themeAppRootDir)
		if err != nil {
			log.Fatal(fmt.Sprintf("watcher.Add themeAppRootDir has error:%v", err))
		}
	}

	// Hang so program doesn't exit
	<-done
	watcher.Close()

}
