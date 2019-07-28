package main

import (
	"todogo/conf"
	"todogo/data"
)

var (
	activeJournal *data.TaskJournal
	activeArchive *data.TaskJournal
)

func loadJournal(filepath string) (*data.TaskJournal, error) {
	var journal data.TaskJournal
	err := journal.LoadOrCreate(filepath)
	if err != nil {
		return nil, err
	}
	return &journal, nil
}

func getActiveJournal() (*data.TaskJournal, error) {
	if activeJournal != nil {
		return activeJournal, nil
	}
	cfg, err := conf.GetConfig()
	if err != nil {
		return nil, err
	}
	return loadJournal(cfg.GetActiveContext().JournalPath())
}

func getActiveArchive() (*data.TaskJournal, error) {
	if activeArchive != nil {
		return activeArchive, nil
	}
	cfg, err := conf.GetConfig()
	if err != nil {
		return nil, err
	}
	return loadJournal(cfg.GetActiveContext().ArchivePath())
}
