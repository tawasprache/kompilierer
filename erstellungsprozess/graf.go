package erstellungsprozess

import (
	"sort"

	"github.com/alecthomas/repr"
)

type Knote string
type Kante struct {
	Von Knote
	Zu  Knote
}

type Graf struct {
	knoten map[Knote]struct{}
	kanten map[Kante]struct{}

	vonKnoten map[Knote]map[Knote]struct{} // which nodes does this come from? n1 -> {(n2, n1), (n3, n1), (n4, n1), ...}
	zuKnoten  map[Knote]map[Knote]struct{} // which nodes does this go to? n1 -> {(n1, n2), (n1, n3), (n1, n4), ...}

	aktivNoten map[Knote]struct{}
}

func NeuGraf() *Graf {
	return &Graf{
		knoten: map[Knote]struct{}{},
		kanten: map[Kante]struct{}{},

		vonKnoten: map[Knote]map[Knote]struct{}{},
		zuKnoten:  map[Knote]map[Knote]struct{}{},

		aktivNoten: map[Knote]struct{}{},
	}
}

func (g *Graf) KnoteHinzufügen(a Knote) {
	g.knoten[a] = struct{}{}
	g.vonKnoten[a] = map[Knote]struct{}{}
	g.zuKnoten[a] = map[Knote]struct{}{}
}

func (g *Graf) KnoteEntfernen(a Knote) {
	delete(g.knoten, a)
	delete(g.vonKnoten, a)
	delete(g.vonKnoten, a)
}

func (g *Graf) KanteHinzufügen(von Knote, zu Knote) {
	kante := Kante{von, zu}

	g.kanten[kante] = struct{}{}
	g.vonKnoten[kante.Zu][kante.Von] = struct{}{}
	g.zuKnoten[kante.Von][kante.Zu] = struct{}{}
}

func (g *Graf) KanteEntfernen(von Knote, zu Knote) {
	kante := Kante{von, zu}

	delete(g.kanten, kante)
	delete(g.vonKnoten[kante.Zu], kante.Von)
	delete(g.zuKnoten[kante.Von], kante.Zu)
}

func (g *Graf) KnoteAktivieren(knote Knote) {
	g.aktivNoten[knote] = struct{}{}
}

func (g *Graf) KnoteDeaktivieren(knote Knote) {
	delete(g.aktivNoten, knote)
}

func hat(s []Knote, e Knote) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// topographicalSort Input: g: a directed acyclic graph with vertices number 1..n
// Output: a linear order of the vertices such that u appears before v
// in the linear order if (u,v) is an edge in the graph.
func topographicalSort(g map[Knote]map[Knote]struct{}) []Knote {
	linearOrder := []Knote{}

	// 1. Let inDegree[1..n] be a new array, and create an empty linear array of
	//    verticies
	inDegree := map[Knote]int{}

	// 2. Set all values in inDegree to 0
	for n := range g {
		inDegree[n] = 0
	}

	// 3. For each vertex u
	for _, adjacent := range g {
		// A. For each vertex *v* adjacent to *u*:
		for v := range adjacent {
			//  i. increment inDegree[v]
			inDegree[v]++
		}
	}

	// 4. Make a list next consisting of all vertices u such that
	//    in-degree[u] = 0
	next := []Knote{}
	for u, v := range inDegree {
		if v != 0 {
			continue
		}

		next = append(next, u)
	}

	// 5. While next is not empty...
	for len(next) > 0 {
		// A. delete a vertex from next and call it vertex u
		u := next[0]
		next = next[1:]

		// B. Add u to the end of the linear order
		linearOrder = append(linearOrder, u)

		// C. For each vertex v adjacent to u
		for v := range g[u] {
			// i. Decrement inDegree[v]
			inDegree[v]--

			// ii. if inDegree[v] = 0, then insert v into next list
			if inDegree[v] == 0 {
				next = append(next, v)
			}
		}
	}

	// 6. Return the linear order
	return linearOrder
}

func idxVon(k Knote, in []Knote) int {
	for idx, it := range in {
		if k == it {
			return idx
		}
	}
	return -1
}

func (g *Graf) AktivierteKnoten() []Knote {
	// TODO: cycles

	var (
		aktivierte []Knote
		todo       []Knote
		k          Knote
	)

	for it := range g.aktivNoten {
		todo = append(todo, it)
	}

	repr.Println(todo)

	for len(todo) != 0 {
		k, todo = todo[0], todo[1:]

		for knote := range g.zuKnoten[k] {
			if !hat(aktivierte, knote) {
				aktivierte = append(aktivierte, knote)
				todo = append(todo, knote)
			}
		}
	}

	ord := topographicalSort(g.zuKnoten)

	sort.Slice(aktivierte, func(i, j int) bool {
		// less

		ii := idxVon(aktivierte[i], ord)
		jj := idxVon(aktivierte[j], ord)

		return ii < jj
	})

	return aktivierte
}
