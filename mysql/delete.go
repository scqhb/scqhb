package main

import "strings"
import "fmt"

func delete() {
	defer wg.Done()
	deleteheader := "delete from bwallcol01 using (values  "
	deletelast := ") as tmp(columall) where tab_colall.columall=tmp.columall;"
	var line20 string
	times := 0
lab002:
	for {
		select {
		case line, ok := <-Chlinedelete:
			if !ok {
				if len(line20) > 1 {
					line20 = strings.TrimRight(line20, ",")
					line20 = fmt.Sprintf("%v %v %v", deleteheader, line20, deletelast)
					fmt.Println(line20)
					_, err = mysqldb.Exec(line20)
					checkErr(err)

				}
				break lab002
			}
			line20 += fmt.Sprintf("('%v'),", line)
			times++
			if times%10000 == 0 {
				line20 = strings.TrimRight(line20, ",")
				line20 = fmt.Sprintf("%v %v %v", deleteheader, line20, deletelast)
				//fmt.Println(line20)
				_, err = mysqldb.Exec(line20)
				checkErr(err)
				times = 0
				line20 = line20[:0]
			}
		}

	}

}
