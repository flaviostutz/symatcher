package symatcher

import (
	"fmt"
	"sort"
)

//Entity entities selected and predicted
type Entity struct {
	Name string
	URL  string
	Tags []string
}

//EntityScore best entities returned along with score
type EntityScore struct {
	ScoredEntity Entity
	Score        int
}

//Training will offer selections, register selections and perform predictions
type Training struct {
	Entities            map[string]Entity
	SelectedTagCounters map[string]int
}

//NewTraining Creates a new training context
func NewTraining(entities []Entity) Training {
	ment := make(map[string]Entity)
	stc := make(map[string]int)
	for _, ent := range entities {
		ment[ent.Name] = ent
		for _, tag := range ent.Tags {
			stc[tag] = 0
		}
	}
	return Training{Entities: ment, SelectedTagCounters: stc}
}

//BestMatches get entities that are closest to the selected entities during training
func (t *Training) BestMatches(minScore int) []EntityScore {
	return t.calculateScores(minScore, t.computeBestTagPoints())
}

func (t *Training) calculateScores(minScore int, tagPoints map[string]int) []EntityScore {
	//calculate the score of how well each entity is fit
	//to the training tag counters
	entityScores := make([]EntityScore, 0)
	for _, ent := range t.Entities {

		//calculate entity score based on its tag points
		score := 0
		for _, entityTag := range ent.Tags {
			score += tagPoints[entityTag]
		}

		if score >= minScore {
			mr := EntityScore{ScoredEntity: ent, Score: score}
			entityScores = append(entityScores, mr)
		}
	}
	//order results
	sort.Slice(entityScores, func(i, j int) bool {
		return entityScores[i].Score < entityScores[j].Score
	})
	return entityScores
}

//compute tag points based on each tag's order of importance
//according to the tag counter
func (t *Training) computeBestTagPoints() map[string]int {

	//calculate the order of relevance of each tag counter
	tagCounters := make([]int, 0)
	positiveCounters := 0
	negativeCounters := 0
	for _, tagCounter := range t.SelectedTagCounters {
		tagCounters = append(tagCounters, tagCounter)
		if tagCounter > 0 {
			positiveCounters++
		} else if tagCounter < 0 {
			negativeCounters++
		}
	}
	//decrescent order
	sort.Slice(tagCounters, func(i, j int) bool {
		//TODO VALIDATE ORDER HERE
		return tagCounters[i] < tagCounters[j]
	})

	tagPoints := make(map[string]int)
	for selectedTag, selectedCounter := range t.SelectedTagCounters {
		for i, orderedCounter := range tagCounters {
			if selectedCounter == orderedCounter {
				if orderedCounter > 0 {
					tagPoints[selectedTag] = positiveCounters - i
				} else if orderedCounter < 0 {
					tagPoints[selectedTag] = -negativeCounters + (len(tagCounters) - i)
				} else {
					tagPoints[selectedTag] = 0
				}
				break
			}
		}
	}

	return tagPoints
}

//NextCandidates get the next candidates that will be used in selection
func (t *Training) NextCandidates(qtty int) ([]Entity, error) {
	if qtty < 2 {
		return []Entity{}, fmt.Errorf("qtty must be >= 2")
	}
	//determine which candidates will help increase or decrease training
	//quality because they have tags that are still close to zero in training

	nearZeroTagPoints := make(map[string]int)
	for selectedTag, selectedCounter := range t.SelectedTagCounters {
		nearZeroTagPoints[selectedTag] = -selectedCounter
	}

	ent := t.calculateScores(-9999, nearZeroTagPoints)

	result := make([]Entity, 0)
	for i, es := range ent {
		if i > qtty {
			break
		}
		result = append(result, es.ScoredEntity)
	}
	return result, nil
}

//Select Tell which entities were selected and which weren't so that we learn
func (t *Training) Select(positiveEntityNames []string, negativeEntityNames []string) error {
	err1 := t.countTags(positiveEntityNames, 1)
	if err1 != nil {
		return err1
	}
	err2 := t.countTags(negativeEntityNames, -1)
	if err2 != nil {
		return err2
	}
	return nil
}

func (t *Training) countTags(entityNames []string, increment int) error {
	for _, en := range entityNames {
		ent, exists := t.Entities[en]
		if !exists {
			return fmt.Errorf("Entity name %s doesn't exist", en)
		}
		for _, tag := range ent.Tags {
			t.SelectedTagCounters[tag] = t.SelectedTagCounters[tag] + increment
		}
	}
	return nil
}
