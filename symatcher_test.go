package symatcher

import (
	"fmt"
	"math/rand"
	"testing"

	"gotest.tools/assert"
)

var (
	rnd = rand.New(rand.NewSource(int64(1)))
)

func TestBasic(t *testing.T) {
	entities := testEntities(10)
	training := NewTraining(entities)

	tagCounters := training.SelectedTagCounters
	assert.Equal(t, tagCounters["a"], 0)
	assert.Equal(t, tagCounters["b"], 0)
	assert.Equal(t, tagCounters["c"], 0)
	assert.Equal(t, tagCounters["d"], 0)
	assert.Equal(t, tagCounters["e"], 0)

	training.Select([]string{entities[0].Name}, []string{entities[2].Name})
	// fmt.Printf("%v\n", tagCounters)
	assert.Equal(t, tagCounters["a"], 0)
	assert.Equal(t, tagCounters["b"], 1)
	assert.Equal(t, tagCounters["c"], 1)
	assert.Equal(t, tagCounters["d"], 0)
	assert.Equal(t, tagCounters["e"], -1)

	training.Select([]string{entities[2].Name}, []string{entities[4].Name})
	// fmt.Printf("%v\n", entities[0:4])

	// fmt.Printf("%v\n", tagCounters)
	es := training.BestMatches(-99999)
	assert.Equal(t, 4, es[0].Score)
	// fmt.Printf("%v\n", es)

	// candidates := training.NextCandidates(2)
	// fmt.Printf("CANDIDATES 111 %v\n", candidates)
	// training.Select([]string{candidates[0].Name}, []string{candidates[1].Name})
	// fmt.Printf("TAG COUNTERS 111 %v\n", tagCounters)
	// // es = training.BestMatches(-99999)
	// // fmt.Printf("%v\n", es)

	// candidates = training.NextCandidates(2)
	// fmt.Printf("CANDIDATES 222 %v\n", candidates)
	// training.Select([]string{candidates[0].Name}, []string{candidates[1].Name})
	// fmt.Printf("TAG COUNTERS 222 %v\n", tagCounters)

	// candidates = training.NextCandidates(2)
	// fmt.Printf("CANDIDATES 333 %v\n", candidates)
	// training.Select([]string{candidates[0].Name}, []string{candidates[1].Name})
	// fmt.Printf("TAG COUNTERS 333 %v\n", tagCounters)
	// // es = training.BestMatches(-99999)
	// // fmt.Printf("BEST 222 %v\n", es)

}

func testEntities(qtty int) []Entity {
	entities := make([]Entity, 0)
	for i := 0; i < qtty; i++ {
		name := fmt.Sprintf("%d", rnd.Int())
		tags := make([]string, 0)
		if rnd.Intn(2) == 1 {
			tags = append(tags, "a")
		}
		if rnd.Intn(2) == 1 {
			tags = append(tags, "b")
		}
		if rnd.Intn(2) == 1 {
			tags = append(tags, "c")
		}
		if rnd.Intn(2) == 1 {
			tags = append(tags, "d")
		}
		if rnd.Intn(2) == 1 {
			tags = append(tags, "e")
		}
		et := Entity{
			Name: name,
			URL:  fmt.Sprintf("http://test.com/img/%s.jpg", name),
			Tags: tags,
		}
		entities = append(entities, et)
	}
	return entities
}
