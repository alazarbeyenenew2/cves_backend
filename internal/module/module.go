package module

type NVD interface {
	Start()
	Stop()
	TriggerNow()
}
