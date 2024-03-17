package cacherinterface

type Displacer interface {
	Run(storage Storage)
	Stop()
}
