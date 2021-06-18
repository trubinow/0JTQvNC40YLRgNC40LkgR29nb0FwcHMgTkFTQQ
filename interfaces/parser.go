package interfaces

type Parser interface {
	Parse(targetDate string) (string, error)
}
