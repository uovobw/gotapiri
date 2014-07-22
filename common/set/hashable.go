package set

// Hashable types can be hashed to a string via the
// Hash function. the set package expects to work on these types only
type Hashable interface {
	Hash() string
}
