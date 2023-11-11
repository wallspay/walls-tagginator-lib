package wallstagginatorlib

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata" // Ensure this package is installed
)

// LookUpRepository is an interface for checking tag uniqueness.
type LookUpRepository interface {
	DataExists(tag string) bool
}

// WallstagGenerator struct holds the TagStorage and a local random source.
type WallstagGenerator struct {
	LookUp  LookUpRepository
	randSrc *rand.Rand
}

// NewWallstagGenerator creates a new WallstagGenerator.
func NewWallstagGenerator(storage LookUpRepository) *WallstagGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(source)

	return &WallstagGenerator{
		LookUp:  storage,
		randSrc: randGen,
	}
}

func (wg *WallstagGenerator) GenerateTags(name string, tagLen int, numTags int) []string {
	var tags []string
	generatedTags := make(map[string]bool)
	var attempts int

	// Generate a mixture of silly and base name tags
	for len(tags) < numTags && attempts < 100 {
		var tag string

		// Alternate between generating a silly name tag and a base name tag
		if attempts%2 == 0 {
			// Generate a silly name tag
			tag = wg.sillyNameWithRandWrapper()
		} else {
			// Generate a base name tag
			tag = generateNameBasedTag(name, tagLen, 1, wg.randSrc)[0]
		}

		// Truncate or pad the tag to ensure it has the fixed length
		if len(tag) > tagLen {
			tag = tag[:tagLen]
		} else {
			for len(tag) < tagLen {
				tag += fmt.Sprintf("%d", wg.randSrc.Intn(10)) // Append a random digit
			}
		}

		// Check for uniqueness and add to the list
		if !wg.LookUp.DataExists(tag) && !generatedTags[tag] {
			tags = append(tags, tag)
			generatedTags[tag] = true
		}

		attempts++
	}

	return tags
}

func (wg *WallstagGenerator) sillyNameWithRandWrapper() string {
	// Save the current global seed
	prevSeed := rand.Int63()
	rand.Seed(wg.randSrc.Int63()) // Set new seed
	name := randomdata.SillyName()
	rand.Seed(prevSeed) // Reset to previous seed
	return strings.ToLower(name)
}

func generateNameBasedTag(name string, tagLen int, numTags int, randSrc *rand.Rand) []string {
	combinedName := strings.ReplaceAll(strings.ToLower(name), " ", "")

	// Extend the combined name to the desired tag length with random digits
	for len(combinedName) < tagLen {
		combinedName += fmt.Sprintf("%d", randSrc.Intn(10)) // Append a random digit
	}

	// Generate unique tags from the combined name
	uniqueTags := make(map[string]bool)
	for i := 0; i <= len(combinedName)-tagLen; i++ {
		tag := combinedName[i : i+tagLen]
		uniqueTags[tag] = true
	}

	// Convert the unique tags map to a slice
	potentialTags := make([]string, 0, len(uniqueTags))
	for tag := range uniqueTags {
		potentialTags = append(potentialTags, tag)
	}

	// Shuffle the slice and select the first numTags elements
	randSrc.Shuffle(len(potentialTags), func(i, j int) {
		potentialTags[i], potentialTags[j] = potentialTags[j], potentialTags[i]
	})

	return potentialTags[:min(numTags, len(potentialTags))]
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
