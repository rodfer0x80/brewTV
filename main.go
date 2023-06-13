package brewTV

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func updateRecords() int {
	fmt.Println("Updated records")
	return
}

func failedToUpdateRecordsWarning() {
	fmt.Println("Failed to update records")
	return
}

func getLibraryRecords() []string {
	var records []string
	return records
}

func getRecordSeasons() []string {
	var seasons []string
	return seasons
}

func getRecordSeasonEpisodes() []string {
	var episodes []string
	return episodes
}

func getEpisodePath(record string, season string, episode string) string {
	var path string
	return path
}

func playEpisode(record string, season string, episode string, path string) {
	save := Save
	save.record = record
	save.season = season
	save.episode = episode
	save.timestamp = 0
	fmt.Printf("Now playing: %s\n", path)
	return
}

type State struct {
	record  string
	season  string
	episode string
	depth   int // 0 :: records; 1 :: seasons // 2 :: episodes
}

type Save struct {
	record    string
	season    string
	episode   string
	timestamp string
}

func main() {
	var records []string
	var seasons []string
	var record int
	var season int
	var episode int
	var path string

	state := State

	var cmd string
	stdin := bufio.NewReader(os.Stdin)

	if updateRecords() != nil {
		failedToUpdateRecordsWarning()
	}
	records = getRecords()

	fmt.Println("<<<<<<<< brewTV >>>>>>>>")
	for i := len(records) - 1; i >= 0; i-- {
		fmt.Printf("%d - %s\n", i+1, records[i])
	}
	fmt.Printf(">>> ")
	cmd, _ = stdin.ReadString('\n')
	cmd = strings.Replace(cmd, "\n", "", -1)
	record, err := strconv.Atoi(cmd)
	if err != nil {
		record = 0
	}

	if cmd == "0" || cmd == "exit" || cmd == "quit" || cmd == "q" || cmd == "x" {
		os.Exit(0)
	} else if record == 0 {
		os.Exit(0)
	} else if 0 < record && record <= len(records) {
		record = record - 1
		state.record = records[record]
		seasons = getRecordSeasons(records[record])
		fmt.Println("<<<<<<<< brewTV >>>>>>>>")
		fmt.Println(records[record])
		for i := len(seasons) - 1; i >= 0; i-- {
			fmt.Printf("%d - %s\n", i+1, seasons[i])
		}
		fmt.Printf(">>> ")
		cmd, _ = stdin.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		season, err := strconv.Atoi(cmd)
		if err != nil {
			season = 0
		}
		if cmd == "0" || cmd == "exit" || cmd == "quit" || cmd == "q" || cmd == "x" {
			// go back
		} else if season == 0 {
			// go back
		} else if 0 < season && season <= len(seasons) {
			season = season - 1
			state.season = seasons[season]
			episodes = getSeasonEpisodes(records[record], seasons[season])
			fmt.Println("<<<<<<<< brewTV >>>>>>>>")
			fmt.Println(seasons[season])
			fmt.Printf(">>> ")
			cmd, _ = stdin.ReadString('\n')
			cmd = strings.Replace(cmd, "\n", "", -1)
			episode, err := strconv.Atoi(cmd)
			if err != nil {
				episode = 0
			}
			if cmd == "0" || cmd == "exit" || cmd == "quit" || cmd == "q" || cmd == "x" {
				// go back
			} else if episode == 0 {
				// go back
			} else if 0 < episode && episode <= len(episodes) {
				episode = episode - 1
				state.episode = episodes[episode]
				path = getEpisodePath(records[record], seasons[season], episodes[episode])
				playEpisode(records[record], seasons[season], episodes[episode])
			}
		}
	}
	return
}
