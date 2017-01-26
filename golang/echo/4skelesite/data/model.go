package data

//Users ..
var Users = make(map[string]User)
var UserFlicks = make(map[string][]Flick)

//User ..
type User struct {
	Name string
}

//Flick ..
type Flick struct {
	Name   string
	Rating float64
}

//Movie ..
type Movie struct {
	Flick
}

//Series ..
type Series struct {
	Flick
	Seasons int
}
