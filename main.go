package main

import (
	"fmt"
	"github.com/rebasar/lunch/fsutil"
	"github.com/rebasar/lunch/lunchy"
	"log"
	"os"
	"strings"
)

const (
	LUNCH_FORMAT        = "%%d | %%-%ds | %%s\n"
	PLACE_HEADER        = "Place"
	DEFAULT_LUNCH_URL   = "https://rebworks.net/lunch/"
	URL_CONFIG_FILENAME = ".lunchy.url"
)

func calculatePlaceNameWidth(places []lunchy.Place) int {
	width := len(PLACE_HEADER)
	for _, place := range places {
		newLength := len(place.Name)
		if newLength > width {
			width = newLength
		}
	}
	return width
}

func getFormatString(width int) string {
	return fmt.Sprintf(LUNCH_FORMAT, width)
}

func printMainMenu(places []lunchy.Place) {
	format := getFormatString(calculatePlaceNameWidth(places))
	fmt.Printf(format, 0, PLACE_HEADER, "Aliases")
	for i, place := range places {
		fmt.Printf(format, i+1, place.Name, strings.Join(place.Aliases, ", "))
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printRestaurantMenu(menu lunchy.Menu) {
	for i, item := range menu.Items {
		fmt.Printf("%d- %s: %s (%d SEK)\n", i+1, item.Title, item.Description, item.Price)
	}
}

func getLunchUrl() string {
	config, err := fsutil.ReadConfigFile(URL_CONFIG_FILENAME)
	if err != nil {
		return DEFAULT_LUNCH_URL
	}
	return strings.TrimSpace(string(config))
}

func main() {
	lunchUrl := getLunchUrl()
	client, err := lunchy.NewClient(lunchUrl)
	handleError(err)
	if len(os.Args) > 1 {
		menu, err := client.GetMenu(os.Args[1])
		handleError(err)
		printRestaurantMenu(menu)
	} else {
		printMainMenu(client.GetPlaces())
	}
}
