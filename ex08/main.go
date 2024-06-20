package main

import (
	"fmt"
	"phnorm/db"
	"phnorm/phone"

	_ "github.com/lib/pq"
)

func main() {
	db, err := db.InitDd()
	must(err)
	defer db.Close()

	phones, err := phone.ListPhones(db)
	must(err)

	for _, p := range phones {
		fmt.Println("Working on...", p.Value)
		n := phone.Normalize(p.Value)

		if p.Normalized == nil || n != *p.Normalized {
			fmt.Print("Updating...", n)
			p.Normalized = &n
			err := phone.UpdateNormalizedNumber(db, p)
			if err != nil {
				fmt.Printf(" - failed to update id:%d; error:%s\n", p.Id, err.Error())
			} else {
				fmt.Println(" - success!")
			}
		} else {
			fmt.Println("No changes necessary")
		}
	}

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
