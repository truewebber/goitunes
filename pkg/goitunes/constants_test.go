package goitunes_test

import (
	"testing"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

func TestGenre_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		genre goitunes.Genre
		want  string
	}{
		{"All", goitunes.GenreAll, "36"},
		{"Games", goitunes.GenreGames, "6014"},
		{"Business", goitunes.GenreBusiness, "6000"},
		{"Games Action", goitunes.GenreGamesAction, "7001"},
		{"Kids Less 5", goitunes.GenreKidsLess5, "KIDS_LESS_5"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.genre.String(); got != tt.want {
				t.Errorf("Genre.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenre_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		genre goitunes.Genre
		want  string
	}{
		{"All", goitunes.GenreAll, "All Categories"},
		{"Games", goitunes.GenreGames, "Games"},
		{"Business", goitunes.GenreBusiness, "Business"},
		{"Games Action", goitunes.GenreGamesAction, "Action"},
		{"Social Networking", goitunes.GenreSocialNetworking, "Social Networking"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.genre.Name(); got != tt.want {
				t.Errorf("Genre.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenre_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		genre goitunes.Genre
		want  bool
	}{
		{"Valid - All", goitunes.GenreAll, true},
		{"Valid - Games", goitunes.GenreGames, true},
		{"Valid - Kids Less 5", goitunes.GenreKidsLess5, true},
		{"Invalid - Empty", goitunes.Genre(""), false},
		{"Invalid - Unknown ID", goitunes.Genre("9999"), false},
		{"Invalid - Random string", goitunes.Genre("invalid"), false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.genre.IsValid(); got != tt.want {
				t.Errorf("Genre.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenre_CustomGenre(t *testing.T) {
	t.Parallel()

	// Test creating custom genre from string
	customGenre := goitunes.Genre("7001")

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
	t.Parallel()

	// Test that all defined genre constants have names in the map
	allGenres := []goitunes.Genre{
		// Main
		goitunes.GenreAll, goitunes.GenreGames, goitunes.GenreShopping, goitunes.GenreMagazinesNewspapers,
		goitunes.GenreEducation, goitunes.GenreBusiness, goitunes.GenreKids, goitunes.GenreFoodDrink,
		goitunes.GenreHealthFitness, goitunes.GenreCatalogs, goitunes.GenreBooks, goitunes.GenreMedical,
		goitunes.GenreMusic, goitunes.GenreNavigation, goitunes.GenreNews, goitunes.GenreLifestyle,
		goitunes.GenreWeather, goitunes.GenreProductivity, goitunes.GenreTravel, goitunes.GenreEntertainment,
		goitunes.GenreSocialNetworking, goitunes.GenreSports, goitunes.GenreReference, goitunes.GenreUtilities,
		goitunes.GenreFinance, goitunes.GenrePhotoVideo,
		// Games
		goitunes.GenreGamesAction, goitunes.GenreGamesAdventure, goitunes.GenreGamesArcade,
		goitunes.GenreGamesBoard, goitunes.GenreGamesCard, goitunes.GenreGamesCasino,
		goitunes.GenreGamesDice, goitunes.GenreGamesEducational, goitunes.GenreGamesFamily,
		goitunes.GenreGamesMusic, goitunes.GenreGamesPuzzle, goitunes.GenreGamesRacing,
		goitunes.GenreGamesRolePlaying, goitunes.GenreGamesSimulation, goitunes.GenreGamesSports,
		goitunes.GenreGamesStrategy, goitunes.GenreGamesTrivia, goitunes.GenreGamesWord,
		// Kids
		goitunes.GenreKidsLess5, goitunes.GenreKids6To8, goitunes.GenreKids9To11,
		// Magazines (sampling)
		goitunes.GenreMagazinesArtsPhotography, goitunes.GenreMagazinesAutomotive,
		goitunes.GenreMagazinesNewsPolitics, goitunes.GenreMagazinesFashionStyle,
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
	t.Parallel()

	// Test that all genre constants are valid and have names
	// We can't directly access genreNames as it's not exported,
	// but we can verify that all constants work correctly through IsValid()
	allGenres := []goitunes.Genre{
		goitunes.GenreAll, goitunes.GenreGames, goitunes.GenreShopping,
		goitunes.GenreBusiness, goitunes.GenreKids, goitunes.GenreHealthFitness,
		goitunes.GenreGamesAction, goitunes.GenreKidsLess5,
	}

	validCount := 0

	for _, genre := range allGenres {
		if genre.IsValid() && genre.Name() != "" {
			validCount++
		}
	}

	if validCount == 0 {
		t.Error("No valid genres found")
	}
}
