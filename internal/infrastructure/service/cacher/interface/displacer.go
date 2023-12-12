package cacher_interface

type Displacer interface {
	Run(storage Storage)
	Stop()
}
