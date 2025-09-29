package IAdapters

import "Kaspersky_Go/ModeLevel/Structures"

type Queue interface {
	Enqueue(job Structures.Job) error
	Dequeue() (Structures.Job, bool)
	Close()
}
