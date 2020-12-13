package gofactory

// func TestDependQueue(t *testing.T) {
// 	q := dependQueue{q: &Queue{}}
// 	q.Enqueue(&dependency{field: "test1", fix: false})
// 	q.Enqueue(&dependency{field: "test2", fix: true})
// 	q.Enqueue(&dependency{field: "test3", fix: false})
// 	q.Enqueue(&dependency{field: "test4", fix: true})
// 	q.Enqueue(&dependency{field: "test5", fix: true})
// 	q.Enqueue(&dependency{field: "test6", fix: false})
// 	q.Enqueue(&dependency{field: "test7", fix: false})

// 	scanner := q.Scan()
// 	count := 0
// 	for depend := scanner(); depend != nil; depend = scanner() {
// 		count++
// 	}
// 	assert.Equal(t, count, 7)
// 	q.clear()
// 	scanner = q.Scan()
// 	count = 0
// 	for depend := scanner(); depend != nil; depend = scanner() {
// 		count++
// 	}
// 	assert.Equal(t, count, 3)
// }
