package cacher

type Displacer interface {
	Run(storage Storage)
	Stop()
}
