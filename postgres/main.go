package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"scdata/public"
	"sync"
)

const (
	HOST        = "gphost"
	DB_USER     = "gpadmin"
	DB_PASSWORD = "oracle"
	DB_NAME     = "postgres"
)

var Wg8 sync.WaitGroup

//var ErrNoRows = errors.New("sql: no rows in result set")

func checkErr(err error) {
	if err != nil {
		panic(err)
	}

}

var Chlineinsert chan string = make(chan string, 10000)
var Chlinedelete chan string = make(chan string, 10000)
var Chlineupdateold chan string = make(chan string, 10000)
var Chlineupdatenew chan string = make(chan string, 10000)
var Chlineselect chan string = make(chan string, 10000)

var db *sql.DB
var err error
var datadir string
var datadirnew string
var cpu int
var usge string = `
postgres -datadir "/data/dat" -f select -cpu 4
postgres -datadir "/data/dat" -f insert -cpu 4
postgres -datadir "/data/dat" -f delete -cpu 4
postgres -datadir "/data/dat"  -datadirnew "/u01/data" -f update  -cpu 4
`
var f string

func init() {
	flag.StringVar(&datadir, "datadir", "/data/dat", "")
	flag.StringVar(&datadirnew, "datadirnew", "/u01/data", "")

	flag.StringVar(&f, "f", "insertq", "")
	flag.IntVar(&cpu, "cpu", 4, "")
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, DB_USER, DB_PASSWORD, DB_NAME)
	db, err = sql.Open("postgres", dbinfo)
	checkErr(err)
	db.SetConnMaxLifetime(500000000000)
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	db.SetConnMaxIdleTime(500000000000)
	db.Stats()
}

var wg sync.WaitGroup

func main() {
	flag.Parse()
	defer db.Close()
	if len(datadir) > 2 && f == "insert" {
		fmt.Println("insert ........")
		go public.ReadfileToChan(datadir, Chlineinsert)
		for i := 0; i < cpu; i++ {
			wg.Add(1)
			go insert()
		}
		wg.Wait()
	}
	if len(datadir) > 2 && len(datadirnew) > 2 && f == "update" {
		fmt.Println("update ........")
		go public.ReadfileToChan(datadir, Chlineupdateold)
		go public.ReadfileToChan(datadirnew, Chlineupdatenew)
		for i := 0; i < cpu; i++ {
			wg.Add(1)
			go Update()
		}
		wg.Wait()
	}

	if len(datadir) > 2 && f == "delete" {
		fmt.Println("delete ........")
		go public.ReadfileToChan(datadir, Chlinedelete)
		for i := 0; i < cpu; i++ {
			wg.Add(1)
			go delete()
		}
		wg.Wait()

	}
	if len(datadir) > 2 && f == "select" {
		fmt.Println("select ........")
		go public.ReadfileToChan(datadir, Chlineselect)
		for i := 0; i < cpu; i++ {
			wg.Add(1)
			go selectx(db)
		}
		wg.Wait()

	} else {
		fmt.Println(usge)
		return
	}

	//将文本文件读取到读取到内存bloom,并启动对外服务
	fmt.Println("insert begin...")
	//insert

	//delete
	///delete from tab_colall using (values (3),(4),(5)) as tmp(columall) where tab_colall.columall=tmp.columall;
	fmt.Println("delete begin...")

	///////////////

	fmt.Println("update begin...")

}
