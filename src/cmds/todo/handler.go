package main

import "galuma.net/todo"

var (
	activeJournal *todo.TaskJournal
	activeArchive *todo.TaskJournal
)

func loadJournal(filepath string) (*todo.TaskJournal, error) {
	var journal todo.TaskJournal
	err := journal.LoadOrCreate(filepath)
	if err != nil {
		return nil, err
	}
	return &journal, nil
}

func getActiveJournal() (*todo.TaskJournal, error) {
	if activeJournal != nil {
		return activeJournal, nil
	}
	cfg, err := todo.GetConfig()
	if err != nil {
		return nil, err
	}
	return loadJournal(cfg.GetActiveContext().JournalPath())
}

func getActiveArchive() (*todo.TaskJournal, error) {
	if activeArchive != nil {
		return activeArchive, nil
	}
	cfg, err := todo.GetConfig()
	if err != nil {
		return nil, err
	}
	return loadJournal(cfg.GetActiveContext().ArchivePath())
}
