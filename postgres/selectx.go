package main

import "database/sql"
import "fmt"

func selectx(dbgreenplum *sql.DB) {
	fmt.Println("select ........................")
	defer wg.Done()

lab002:
	for {
		select {
		case line, ok := <-Chlineselect:
			if !ok {
				break lab002
			}
			var username string
			// line20 = fmt.Sprintf("%v $1 limit 1;", inserthead)
			//	 fmt.Println(line20)
			err = dbgreenplum.QueryRow("select true from bwallcol01 where columall=$1 limit 1;", line).Scan(&username)
			switch {
			case err == sql.ErrNoRows:
				//fmt.Println("no found")
			case err != nil:
				{
					checkErr(err)
				}

			}
			// fmt.Println(username)
			/* 			for rows.Next() {
			 				var username string
			 				err = rows.Scan( &username )
							checkErr(err)
			 			//	fmt.Printf("%8v",username )
						}
						rows.Close()*/

			//	dbgreenplum.Close()
		}

	}

}
