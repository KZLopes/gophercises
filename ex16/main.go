package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"main/internal/twitter"
	"math/rand"
	"os"
	"time"
)

type TwitterAPI interface {
	Get(string) (string, error)
}

func main() {
	usersFile := flag.String("file", "users.csv", "file to persist")
	pickWinner := flag.Bool("end", false, "")
	flag.Parse()

	client := twitter.MockedAPI{}
	rts, err := client.GetParsedRetweets()
	if err != nil {
		panic(err)
	}

	newUsernames := make([]string, 0, len(rts))
	for _, rt := range rts {
		newUsernames = append(newUsernames, rt.User.ScreenName)
	}
	oldUsernames, err := existingUsernames(*usersFile)
	if err != nil {
		panic(err)
	}
	allUsernames := merge(oldUsernames, newUsernames)

	err = persistUsernames(*usersFile, allUsernames)
	if err != nil {
		panic(err)
	}
	if *pickWinner {
		winner, err := runGiveaway(*usersFile)
		if err != nil {
			panic(err)
		}
		fmt.Println("The winner is:", winner)
	}

}

func runGiveaway(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	rng := rand.New(rand.NewSource(time.Now().Unix() + int64(rand.Int())))
	winner := rng.Intn(len(records))

	return records[winner][0], nil
}

func persistUsernames(file string, usernames []string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for _, username := range usernames {
		err = w.Write([]string{username})
		if err != nil {
			return err
		}
	}
	w.Flush()
	if err = w.Error(); err != nil {
		return err
	}
	return nil
}

func existingUsernames(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return []string{}, nil
	}
	defer f.Close()
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(lines))
	for _, line := range lines {
		ret = append(ret, line[0])
	}
	return ret, nil
}
func merge(old, new []string) []string {
	usernames := make(map[string]struct{}, 0)

	for _, name := range old {
		usernames[name] = struct{}{}
	}
	for _, name := range new {
		usernames[name] = struct{}{}
	}

	var ret []string
	for k := range usernames {
		ret = append(ret, k)
	}

	return ret

	// check if file exists
	// if it does, open and copy the current contetn/username

	// if doesn't create the file

	//write to file merging the old an new usernames
	// usernames cannt be duplicated

	//		f, err := os.Open(*usersFile)
	//		if err != nil {
	//			if err == os.ErrNotExist {
	//				f, _ = os.Create(*usersFile)
	//			} else {
	//				panic(err)
	//			}
	//		}
	//
	//		r := csv.NewReader(f)
	//		_ = r
	//	}
}
