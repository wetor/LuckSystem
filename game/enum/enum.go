package enum

type VMRunMode int8

const (
	VMRun VMRunMode = iota
	VMRunExport
	VMRunImport
)
