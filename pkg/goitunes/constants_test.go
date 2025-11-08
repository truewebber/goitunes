package goitunes

import "testing"

func TestGenre_String(t *testing.T) {
	tests := []struct {
		name  string
		genre Genre
		want  string
	}{
		{"All", GenreAll, "36"},
		{"Games", GenreGames, "6014"},
		{"Business", GenreBusiness, "6000"},
		{"Games Action", GenreGamesAction, "7001"},
		{"Kids Less 5", GenreKidsLess5, "KIDS_LESS_5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.genre.String(); got != tt.want {
				t.Errorf("Genre.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenre_Name(t *testing.T) {
	tests := []struct {
		name  string
		genre Genre
		want  string
	}{
		{"All", GenreAll, "All Categories"},
		{"Games", GenreGames, "Games"},
		{"Business", GenreBusiness, "Business"},
		{"Games Action", GenreGamesAction, "Action"},
		{"Social Networking", GenreSocialNetworking, "Social Networking"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.genre.Name(); got != tt.want {
				t.Errorf("Genre.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenre_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		genre Genre
		want  bool
	}{
		{"Valid - All", GenreAll, true},
		{"Valid - Games", GenreGames, true},
		{"Valid - Kids Less 5", GenreKidsLess5, true},
		{"Invalid - Empty", Genre(""), false},
		{"Invalid - Unknown ID", Genre("9999"), false},
		{"Invalid - Random string", Genre("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.genre.IsValid(); got != tt.want {
				t.Errorf("Genre.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenre_CustomGenre(t *testing.T) {
	// Test creating custom genre from string
	customGenre := Genre("7001")

	if !customGenre.IsValid() {
		t.Error("Expected custom genre '7001' to be valid")
	}

	if customGenre.Name() != "Action" {
		t.Errorf("Expected genre name 'Action', got '%s'", customGenre.Name())
	}

	if customGenre.String() != "7001" {
		t.Errorf("Expected genre string '7001', got '%s'", customGenre.String())
	}
}

func TestGenreNames_AllGenresHaveNames(t *testing.T) {
	// Test that all defined genre constants have names in the map
	allGenres := []Genre{
		// Main
		GenreAll, GenreGames, GenreShopping, GenreMagazinesNewspapers,
		GenreEducation, GenreBusiness, GenreKids, GenreFoodDrink,
		GenreHealthFitness, GenreCatalogs, GenreBooks, GenreMedical,
		GenreMusic, GenreNavigation, GenreNews, GenreLifestyle,
		GenreWeather, GenreProductivity, GenreTravel, GenreEntertainment,
		GenreSocialNetworking, GenreSports, GenreReference, GenreUtilities,
		GenreFinance, GenrePhotoVideo,
		// Games
		GenreGamesAction, GenreGamesAdventure, GenreGamesArcade,
		GenreGamesBoard, GenreGamesCard, GenreGamesCasino,
		GenreGamesDice, GenreGamesEducational, GenreGamesFamily,
		GenreGamesMusic, GenreGamesPuzzle, GenreGamesRacing,
		GenreGamesRolePlaying, GenreGamesSimulation, GenreGamesSports,
		GenreGamesStrategy, GenreGamesTrivia, GenreGamesWord,
		// Kids
		GenreKidsLess5, GenreKids6To8, GenreKids9To11,
		// Magazines (sampling)
		GenreMagazinesArtsPhotography, GenreMagazinesAutomotive,
		GenreMagazinesNewsPolitics, GenreMagazinesFashionStyle,
	}

	for _, genre := range allGenres {
		if !genre.IsValid() {
			t.Errorf("Genre %s should be valid", genre.String())
		}

		name := genre.Name()
		if name == "" {
			t.Errorf("Genre %s should have a name", genre.String())
		}
	}
}

func TestGenreNames_Count(t *testing.T) {
	// We should have exactly 73 genre constants
	// (72 unique IDs + GenreAll="36" which is also in the main categories list)
	expectedCount := 73
	actualCount := len(genreNames)

	if actualCount != expectedCount {
		t.Errorf("Expected %d genres, but found %d", expectedCount, actualCount)
	}
}
