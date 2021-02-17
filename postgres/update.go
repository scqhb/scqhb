package main

import "strings"
import "fmt"

func Update() {
	defer wg.Done()
	//update tab_colall set columall=tmp.idnew from (values ('222','222bbb'),('333','333bbb'),('111','111bbb')) as tmp (columall,idnew) where tb002.id=tmp.id;
	updateheader := "update bwallcol01 set columall=tmp.idnew from (values   "
	updatelast := ") as tmp (columall,idnew) where tab_colall.columall=tmp.id; "
	var line20 string
	times := 0
lab002:
	for {
		select {
		case linenew, ok := <-Chlineupdatenew:
			if !ok {
				if len(line20) > 1 {
					line20 = strings.TrimRight(line20, ",")
					line20 = fmt.Sprintf("%v %v %v", updateheader, line20, updatelast)
					fmt.Println(line20)
					_, err = db.Exec(line20)
					checkErr(err)

				}
				break lab002
			}
			line20 += fmt.Sprintf("('%v','%v'),", <-Chlineupdateold, linenew)
			times++
			if times%1000 == 0 {
				line20 = strings.TrimRight(line20, ",")
				line20 = fmt.Sprintf("%v %v %v", updateheader, line20, updatelast)
				//fmt.Println(line20)
				_, err = db.Exec(line20)
				checkErr(err)
				times = 0
				line20 = line20[:0]
			}
		}

	}

}
