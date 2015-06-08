package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/fusionrsrch/go-mysql-elasticsearch/river"
)

var configFile = flag.String("config", "./etc/river.yaml", "go-mysql-elasticsearch config file")
var my_addr = flag.String("my_addr", "", "MySQL addr")
var my_user = flag.String("my_user", "", "MySQL user")
var my_pass = flag.String("my_pass", "", "MySQL password")
var es_addr = flag.String("es_addr", "", "Elasticsearch addr")
var data_dir = flag.String("data_dir", "", "path for go-mysql-elasticsearch to save data")
var server_id = flag.Int("server_id", 0, "MySQL server id, as a pseudo slave")
var flavor = flag.String("flavor", "", "flavor: mysql or mariadb")
var execution = flag.String("exec", "", "mysqldump execution path")

func main() {

	fmt.Println("Go MySQL ElasticSearch")

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	fmt.Println("Loading config ...")
	cfg, err := river.NewConfigWithFile(*configFile)

	if err != nil {
		//println(err.Error())
		fmt.Println(err)
		return
	} else {
		fmt.Println(cfg)
	}

	if len(*my_addr) > 0 {
		cfg.MyAddr = *my_addr
	}

	if len(*my_user) > 0 {
		cfg.MyUser = *my_user
	}

	if len(*my_pass) > 0 {
		cfg.MyPassword = *my_pass
	}

	if *server_id > 0 {
		cfg.ServerID = uint32(*server_id)
	}

	if len(*es_addr) > 0 {
		cfg.ESAddr = *es_addr
	}

	if len(*data_dir) > 0 {
		cfg.DataDir = *data_dir
	}

	if len(*flavor) > 0 {
		cfg.Flavor = *flavor
	}

	if len(*execution) > 0 {
		cfg.DumpExec = *execution
	}

	fmt.Println("Start river.NewRiver() ...")
	r, err := river.NewRiver(cfg)
	if err != nil {
		//println(err.Error())
		fmt.Println(err)
		return
	}
	fmt.Println(r)

	go func() {
		<-sc
		r.Close()
	}()

	r.Run()
}
