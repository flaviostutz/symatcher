package symatcher

import (
	"fmt"
	"math"
	"math/rand"
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

//DiscriminationLevel calculates the difference between the [average score from N best elements]
//and [average score of remaining matches]
func (t *Training) DiscriminationLevel(bestMatchesCount int) int {
	bm := t.BestMatches(-99999)
	bm1 := bm[:bestMatchesCount]
	bm2 := bm[bestMatchesCount+1:]
	as1 := averageScore(bm1)
	as2 := averageScore(bm2)
	return as1 - as2
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
		return entityScores[i].Score > entityScores[j].Score
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
		return tagCounters[i] > tagCounters[j]
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
func (t *Training) NextCandidates(qtty int, randomRange int) []Entity {
	//get candidates based on how its attributes have zero counter in training
	//so that those attributes will be evaluated when those candidates are selected
	//and training will get stronger
	nearZeroTagPoints := make(map[string]int)
	for selectedTag, selectedCounter := range t.SelectedTagCounters {
		nearZeroTagPoints[selectedTag] = -selectedCounter
	}

	//get entities whose tags are near zero point in training
	ent := t.calculateScores(-9999, nearZeroTagPoints)

	r1 := make([]Entity, 0)
	for i, es := range ent {
		if i >= randomRange {
			break
		}
		r1 = append(r1, es.ScoredEntity)
	}

	if len(r1) < qtty {
		return r1
	}

	//get randomized elements to diversify training
	rnd := rand.New(rand.NewSource(1))
	resultIndexes := rnd.Perm(len(r1))

	r2 := make([]Entity, 0)
	for _, v := range resultIndexes[:qtty] {
		r2 = append(r2, r1[v])
	}
	return r2
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

func averageScore(entities []EntityScore) int {
	tv := 0
	for _, en := range entities {
		tv += en.Score
	}
	return int(math.Floor(float64(tv) / float64(len(entities))))
}
