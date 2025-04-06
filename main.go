package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const COOKIES = ""

const SAMPLE = `
{
    "response":
    {
        "character_league_infos":
        [
            {
                "character_id": 1,
                "is_played": false,
                "league_info":
                {
                    "league_point": -1,
                    "league_rank": 39,
                    "master_league": 0,
                    "master_rating": 0,
                    "master_rating_ranking": 0
                },
                "character_name": "Ryu",
                "character_alpha": "RYU",
                "character_tool_name": "ryu",
                "character_sort": 12
            },
            ...
        ]
	}
}
`

type AccountLeagueInfo struct {
	Response struct {
		CharacterLeagueInfos []struct {
			ID            int    `json:"character_id"`
			IsPlayed      bool   `json:"is_played"`
			CharacterName string `json:"character_name"`
			LeagueInfo    struct {
				LeaguePoints        int `json:"league_point"`
				LeagueRank          int `json:"league_rank"`
				MasterLeague        int `json:"master_league"`
				MasterRating        int `json:"master_rating"`
				MasterRatingRanking int `json:"master_rating_ranking"`
			} `json:"league_info"`
		} `json:"character_league_infos"`
	} `json:"response"`
}

func getHighestCharacterAndLP(accountLeagueInfos AccountLeagueInfo) {
	highestLPCharacter := ""
	highestLP := -2
	highestMRCharacter := ""
	highestMR := 0

	for _, character := range accountLeagueInfos.Response.CharacterLeagueInfos {
		if character.IsPlayed && character.LeagueInfo.LeaguePoints > highestLP {
			highestLP = character.LeagueInfo.LeaguePoints
			highestLPCharacter = character.CharacterName
		}
		if character.LeagueInfo.MasterRating > highestMR {
			highestMR = character.LeagueInfo.MasterRating
			highestMRCharacter = character.CharacterName
		}
	}

	fmt.Printf("Highest character is %s with LP of %d\n", highestLPCharacter, highestLP)
	if highestLP > 25000 {
		fmt.Printf("Master character detected. ")
		if highestMR > 0 {
			fmt.Printf("Highest MR rated character this season is %s, with MR of %d\n", highestMRCharacter, highestMR)
		} else {
			fmt.Printf("However, they have not played any games on their Master Ranked characters this season.")
		}
	}
}

func getCharacterLeagueInfo(capcomId int) (AccountLeagueInfo, error) {
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"targetShortId":%d,"targetSeasonId":-1,"locale":"en","peak":true}`, capcomId))
	// var data = strings.NewReader(`{"targetShortId":1681090405,"targetSeasonId":-1,"locale":"en","peak":true}`)
	req, err := http.NewRequest("POST", "https://www.streetfighter.com/6/buckler/api/profile/play/act/leagueinfo", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:136.0) Gecko/20100101 Firefox/136.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Referer", fmt.Sprintf("https://www.streetfighter.com/6/buckler/profile/%d/play", capcomId))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://www.streetfighter.com")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", COOKIES)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=0")
	req.Header.Set("TE", "trailers")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var accountInfo AccountLeagueInfo
	err = json.Unmarshal([]byte(bodyText), &accountInfo)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return accountInfo, err
}

func main() {
	capcomIdStr := flag.String("cid", "1681080405", "The Capcom ID to look up.")
	flag.Parse()
	fmt.Printf("Looking up highest character and rank belonging to account %s!\n", *capcomIdStr)
	capcomId, err := strconv.Atoi(*capcomIdStr)
	if err != nil {
		log.Fatal(err)
	}
	accountInfo, err := getCharacterLeagueInfo(capcomId)
	if err != nil {
		log.Fatal(err)
	}
	getHighestCharacterAndLP(accountInfo)
}
