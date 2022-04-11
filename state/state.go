package state

type State interface {
	SetValue() error
	SetAttribute() error
}
