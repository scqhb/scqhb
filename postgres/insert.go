package main

import "strings"
import "fmt"

func insert() {
	defer wg.Done()
	inserthead := "insert into bwallcol01 (columall ) values "
	var line20 string
	times := 0
lab002:
	for {
		select {
		case line, ok := <-Chlineinsert:
			if !ok {
				if len(line20) > 1 {
					line20 = strings.TrimRight(line20, ",")
					line20 = fmt.Sprintf("%v %v", inserthead, line20)
					fmt.Println(line20)
					_, err = db.Exec(line20)
					checkErr(err)
					line20 = ""
					fmt.Println("################################################################")
				}
				break lab002
			}
			line20 += fmt.Sprintf("('%v'),", line)
			times++
			if times%1000 == 0 {
				line20 = strings.TrimRight(line20, ",")
				line20 = fmt.Sprintf("%v %v", inserthead, line20)
				//fmt.Println(line20)
				tx, err := db.Begin()
				checkErr(err)
				_, err = tx.Exec(line20)
				checkErr(err)
				tx.Commit()
				//	db.Close()

				times = 0
				line20 = ""
			}
		}

	}

}
