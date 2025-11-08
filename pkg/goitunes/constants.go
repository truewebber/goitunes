package goitunes

// Genre represents an App Store genre/category identifier.
type Genre string

// String returns the string representation of the genre ID.
func (g Genre) String() string {
	return string(g)
}

// IsValid checks if the genre is a valid known genre.
func (g Genre) IsValid() bool {
	_, exists := genreNames[g]

	return exists
}

// Name returns the human-readable name of the genre.
func (g Genre) Name() string {
	return genreNames[g]
}

// App Store genres.
const (
	GenreAll                        Genre = "36"           // All categories
	GenreGames                      Genre = "6014"         // Games
	GenreGamesAction                Genre = "7001"         // Games Action
	GenreGamesAdventure             Genre = "7002"         // Games Adventure
	GenreGamesArcade                Genre = "7003"         // Games Arcade
	GenreGamesBoard                 Genre = "7004"         // Games Board
	GenreGamesCard                  Genre = "7005"         // Games Card
	GenreGamesCasino                Genre = "7006"         // Games Casino
	GenreGamesDice                  Genre = "7007"         // Games Dice
	GenreGamesEducational           Genre = "7008"         // Games Educational
	GenreGamesFamily                Genre = "7009"         // Games Family
	GenreGamesMusic                 Genre = "7011"         // Games Music
	GenreGamesPuzzle                Genre = "7012"         // Games Puzzle
	GenreGamesRacing                Genre = "7013"         // Games Racing
	GenreGamesRolePlaying           Genre = "7014"         // Games Role-Playing
	GenreGamesSimulation            Genre = "7015"         // Games Simulation
	GenreGamesSports                Genre = "7016"         // Games Sports
	GenreGamesStrategy              Genre = "7017"         // Games Strategy
	GenreGamesTrivia                Genre = "7018"         // Games Trivia
	GenreGamesWord                  Genre = "7019"         // Games Word
	GenreShopping                   Genre = "6024"         // Shopping
	GenreMagazinesNewspapers        Genre = "6021"         // Magazines & Newspapers
	GenreMagazinesArtsPhotography   Genre = "13007"        // Magazines & Newspapers Arts & Photography
	GenreMagazinesAutomotive        Genre = "13006"        // Magazines & Newspapers Automotive
	GenreMagazinesBridesWeddings    Genre = "13008"        // Magazines & Newspapers Brides & Weddings
	GenreMagazinesBusinessInvesting Genre = "13009"        // Magazines & Newspapers Business & Investing
	GenreMagazinesChildrens         Genre = "13010"        // Magazines & Newspapers Children's Magazines
	GenreMagazinesComputingInternet Genre = "13011"        // Magazines & Newspapers Computing & Internet
	GenreMagazinesCookingFoodDrink  Genre = "13012"        // Magazines & Newspapers Cooking, Food & Drink
	GenreMagazinesCraftsHobbies     Genre = "13013"        // Magazines & Newspapers Crafts & Hobbies
	GenreMagazinesElectronicsAudio  Genre = "13014"        // Magazines & Newspapers Electronics & Audio
	GenreMagazinesEntertainment     Genre = "13015"        // Magazines & Newspapers Entertainment
	GenreMagazinesFashionStyle      Genre = "13002"        // Magazines & Newspapers Fashion & Style
	GenreMagazinesFilmsMusic        Genre = "13021"        // Magazines & Newspapers Films & Music
	GenreMagazinesHealthWellBeing   Genre = "13017"        // Magazines & Newspapers Health & Well-Being
	GenreMagazinesHistory           Genre = "13018"        // Magazines & Newspapers History
	GenreMagazinesHomeGarden        Genre = "13003"        // Magazines & Newspapers Home & Garden
	GenreMagazinesLiteraryJournals  Genre = "13019"        // Magazines & Newspapers Literary Magazines & Journals
	GenreMagazinesMensInterest      Genre = "13020"        // Magazines & Newspapers Men's Interest
	GenreMagazinesNewsPolitics      Genre = "13001"        // Magazines & Newspapers News & Politics
	GenreMagazinesOutdoorsNature    Genre = "13004"        // Magazines & Newspapers Outdoors & Nature
	GenreMagazinesPets              Genre = "13024"        // Magazines & Newspapers Pets
	GenreMagazinesProfessionalTrade Genre = "13025"        // Magazines & Newspapers Professional & Trade
	GenreMagazinesScience           Genre = "13027"        // Magazines & Newspapers Science
	GenreMagazinesSportsLeisure     Genre = "13005"        // Magazines & Newspapers Sports & Leisure
	GenreMagazinesTeens             Genre = "13028"        // Magazines & Newspapers Teens
	GenreMagazinesTravelRegional    Genre = "13029"        // Magazines & Newspapers Travel & Regional
	GenreMagazinesWomensInterest    Genre = "13030"        // Magazines & Newspapers Women's Interest
	GenreEducation                  Genre = "6017"         // Education
	GenreBusiness                   Genre = "6000"         // Business
	GenreKids                       Genre = "KIDS"         // Kids
	GenreKidsLess5                  Genre = "KIDS_LESS_5"  // Kids 5 & Under
	GenreKids6To8                   Genre = "KIDS_6_TO_8"  // Kids 6–8
	GenreKids9To11                  Genre = "KIDS_9_TO_11" // Kids 9–11
	GenreFoodDrink                  Genre = "6023"         // Food & Drink
	GenreHealthFitness              Genre = "6013"         // Health & Fitness
	GenreCatalogs                   Genre = "6022"         // Catalogs
	GenreBooks                      Genre = "6018"         // Books
	GenreMedical                    Genre = "6020"         // Medical
	GenreMusic                      Genre = "6011"         // Music
	GenreNavigation                 Genre = "6010"         // Navigation
	GenreNews                       Genre = "6009"         // News
	GenreLifestyle                  Genre = "6012"         // Lifestyle
	GenreWeather                    Genre = "6001"         // Weather
	GenreProductivity               Genre = "6007"         // Productivity
	GenreTravel                     Genre = "6003"         // Travel
	GenreEntertainment              Genre = "6016"         // Entertainment
	GenreSocialNetworking           Genre = "6005"         // Social Networking
	GenreSports                     Genre = "6004"         // Sports
	GenreReference                  Genre = "6006"         // Reference
	GenreUtilities                  Genre = "6002"         // Utilities
	GenreFinance                    Genre = "6015"         // Finance
	GenrePhotoVideo                 Genre = "6008"         // Photo & Video
)

// genreNames maps genre IDs to their human-readable names.
//
//nolint:gochecknoglobals // map is acceptable as a global variable
var genreNames = map[Genre]string{
	// Main genres
	GenreAll:                 "All Categories",
	GenreGames:               "Games",
	GenreShopping:            "Shopping",
	GenreMagazinesNewspapers: "Magazines & Newspapers",
	GenreEducation:           "Education",
	GenreBusiness:            "Business",
	GenreKids:                "Kids",
	GenreFoodDrink:           "Food & Drink",
	GenreHealthFitness:       "Health & Fitness",
	GenreCatalogs:            "Catalogs",
	GenreBooks:               "Books",
	GenreMedical:             "Medical",
	GenreMusic:               "Music",
	GenreNavigation:          "Navigation",
	GenreNews:                "News",
	GenreLifestyle:           "Lifestyle",
	GenreWeather:             "Weather",
	GenreProductivity:        "Productivity",
	GenreTravel:              "Travel",
	GenreEntertainment:       "Entertainment",
	GenreSocialNetworking:    "Social Networking",
	GenreSports:              "Sports",
	GenreReference:           "Reference",
	GenreUtilities:           "Utilities",
	GenreFinance:             "Finance",
	GenrePhotoVideo:          "Photo & Video",

	// Game sub-genres
	GenreGamesAction:      "Action",
	GenreGamesAdventure:   "Adventure",
	GenreGamesArcade:      "Arcade",
	GenreGamesBoard:       "Board",
	GenreGamesCard:        "Card",
	GenreGamesCasino:      "Casino",
	GenreGamesDice:        "Dice",
	GenreGamesEducational: "Educational",
	GenreGamesFamily:      "Family",
	GenreGamesMusic:       "Music",
	GenreGamesPuzzle:      "Puzzle",
	GenreGamesRacing:      "Racing",
	GenreGamesRolePlaying: "Role-Playing",
	GenreGamesSimulation:  "Simulation",
	GenreGamesSports:      "Sports",
	GenreGamesStrategy:    "Strategy",
	GenreGamesTrivia:      "Trivia",
	GenreGamesWord:        "Word",

	// Kids sub-genres
	GenreKidsLess5: "Kids 5 & Under",
	GenreKids6To8:  "Kids 6–8",
	GenreKids9To11: "Kids 9–11",

	// Magazines sub-genres
	GenreMagazinesArtsPhotography:   "Arts & Photography",
	GenreMagazinesAutomotive:        "Automotive",
	GenreMagazinesBridesWeddings:    "Brides & Weddings",
	GenreMagazinesBusinessInvesting: "Business & Investing",
	GenreMagazinesChildrens:         "Children's Magazines",
	GenreMagazinesComputingInternet: "Computing & Internet",
	GenreMagazinesCookingFoodDrink:  "Cooking, Food & Drink",
	GenreMagazinesCraftsHobbies:     "Crafts & Hobbies",
	GenreMagazinesElectronicsAudio:  "Electronics & Audio",
	GenreMagazinesEntertainment:     "Entertainment",
	GenreMagazinesFashionStyle:      "Fashion & Style",
	GenreMagazinesFilmsMusic:        "Films & Music",
	GenreMagazinesHealthWellBeing:   "Health & Well-Being",
	GenreMagazinesHistory:           "History",
	GenreMagazinesHomeGarden:        "Home & Garden",
	GenreMagazinesLiteraryJournals:  "Literary Magazines & Journals",
	GenreMagazinesMensInterest:      "Men's Interest",
	GenreMagazinesNewsPolitics:      "News & Politics",
	GenreMagazinesOutdoorsNature:    "Outdoors & Nature",
	GenreMagazinesPets:              "Pets",
	GenreMagazinesProfessionalTrade: "Professional & Trade",
	GenreMagazinesScience:           "Science",
	GenreMagazinesSportsLeisure:     "Sports & Leisure",
	GenreMagazinesTeens:             "Teens",
	GenreMagazinesTravelRegional:    "Travel & Regional",
	GenreMagazinesWomensInterest:    "Women's Interest",
}

// Common user agents.
const (
	UserAgentWindows = "iTunes/10.6 (Windows; Microsoft Windows 7 x64 Ultimate Edition " +
		"Service Pack 1 (Build 7601)) AppleWebKit/534.54.16"
	UserAgentTop200   = "AppStore/2.0 iOS/9.0 model/iPhone6,1 hwp/s5l8960x build/13A344 (6; dt:89)"
	UserAgentTop1500  = "iTunes-iPad/5.1.1 (64GB; dt:28)"
	UserAgentDownload = "itunesstored/1.0 iOS/9.0 model/iPhone6,1 hwp/s5l8960x build/13A344 (6; dt:89)"
)
