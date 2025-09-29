package IAdapters

import "Kaspersky_Go/ModeLevel/Structures"

type StateStore interface {
	Set(id string, st Structures.JobStatus)
	Get(id string) (Structures.JobStatus, bool)
}
