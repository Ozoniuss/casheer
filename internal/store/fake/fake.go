package fake

import (
	"slices"

	"github.com/Ozoniuss/casheer/internal/domain/model"
	"github.com/Ozoniuss/casheer/internal/store"
)

type FakeStore struct {
	Entries  []model.Entry
	Expenses []model.Expense
	Debts    []model.Debt
}

func (f *FakeStore) AddEntry(new model.Entry) (model.Entry, error) {
	if slices.ContainsFunc(f.Entries, func(e model.Entry) bool {
		return new.Id == e.Id
	}) {
		return model.Entry{}, store.ErrAlreadyExists{}
	}

	if slices.ContainsFunc(f.Entries, func(e model.Entry) bool {
		return new.Category == e.Category &&
			new.Subcategory == e.Subcategory &&
			new.Month == e.Month &&
			new.Year == e.Year
	}) {
		return model.Entry{}, store.ErrInvalidConstrains{}
	}

	f.Entries = append(f.Entries, new)
	return new, nil
}

func (f *FakeStore) DeleteEntry(id model.Id) (model.Entry, error) {
	toDelete, err := f.GetEntry(id)
	if err != nil {
		return model.Entry{}, err
	}
	f.Entries = slices.DeleteFunc(f.Entries, func(e model.Entry) bool {
		return e.Id == toDelete.Id
	})
	return toDelete, nil
}

func (f *FakeStore) GetEntry(id model.Id) (model.Entry, error) {
	idx := slices.IndexFunc(f.Entries, func(e model.Entry) bool {
		return id == e.Id
	})
	if idx < 0 {
		return model.Entry{}, store.ErrNotFound{}
	}
	return f.Entries[id], nil
}

func (f *FakeStore) ListEntries() ([]model.Entry, error) {
	return f.Entries, nil
}

func (f *FakeStore) UpdateEntry(new model.Entry) (model.Entry, error) {
	idx := slices.IndexFunc(f.Entries, func(e model.Entry) bool {
		return new.Id == e.Id
	})
	if idx < 0 {
		return model.Entry{}, store.ErrNotFound{}
	}

	f.Entries[idx] = new
	return new, nil
}

func basicAdd[T model.Entry | model.Debt | model.Expense](new T, objectlist *[]T, eqfuncs ...func(t1, t2 T) bool) (T, error) {
	var empty T

	// Verify the equality constraints.
	for _, eqfunc := range eqfuncs {
		if slices.ContainsFunc(*objectlist, func(e T) bool {
			return eqfunc(new, e)
		}) {
			return empty, store.ErrAlreadyExists{}
		}
	}

	*objectlist = append(*objectlist, new)
	return new, nil
}

func (f *FakeStore) AddEntry(new model.Entry) (model.Entry, error) {
	if slices.ContainsFunc(f.Entries, func(e model.Entry) bool {
		return new.Id == e.Id
	}) {
		return model.Entry{}, store.ErrAlreadyExists{}
	}

	if slices.ContainsFunc(f.Entries, func(e model.Entry) bool {
		return new.Category == e.Category &&
			new.Subcategory == e.Subcategory &&
			new.Month == e.Month &&
			new.Year == e.Year
	}) {
		return model.Entry{}, store.ErrInvalidConstrains{}
	}

	f.Entries = append(f.Entries, new)
	return new, nil
}

func (f *FakeStore) DeleteEntry(id model.Id) (model.Entry, error) {
	toDelete, err := f.GetEntry(id)
	if err != nil {
		return model.Entry{}, err
	}
	f.Entries = slices.DeleteFunc(f.Entries, func(e model.Entry) bool {
		return e.Id == toDelete.Id
	})
	return toDelete, nil
}

func (f *FakeStore) GetEntry(id model.Id) (model.Entry, error) {
	idx := slices.IndexFunc(f.Entries, func(e model.Entry) bool {
		return id == e.Id
	})
	if idx < 0 {
		return model.Entry{}, store.ErrNotFound{}
	}
	return f.Entries[id], nil
}

func (f *FakeStore) ListEntries() ([]model.Entry, error) {
	return f.Entries, nil
}

func (f *FakeStore) UpdateEntry(new model.Entry) (model.Entry, error) {
	idx := slices.IndexFunc(f.Entries, func(e model.Entry) bool {
		return new.Id == e.Id
	})
	if idx < 0 {
		return model.Entry{}, store.ErrNotFound{}
	}

	f.Entries[idx] = new
	return new, nil
}

func (f *FakeStore) AddEntry(new model.Entry) (model.Entry, error) {
	if slices.ContainsFunc(f.Entries, func(e model.Entry) bool {
		return new.Id == e.Id
	}) {
		return model.Entry{}, store.ErrAlreadyExists{}
	}

	if slices.ContainsFunc(f.Entries, func(e model.Entry) bool {
		return new.Category == e.Category &&
			new.Subcategory == e.Subcategory &&
			new.Month == e.Month &&
			new.Year == e.Year
	}) {
		return model.Entry{}, store.ErrInvalidConstrains{}
	}

	f.Entries = append(f.Entries, new)
	return new, nil
}

func (f *FakeStore) DeleteEntry(id model.Id) (model.Entry, error) {
	toDelete, err := f.GetEntry(id)
	if err != nil {
		return model.Entry{}, err
	}
	f.Entries = slices.DeleteFunc(f.Entries, func(e model.Entry) bool {
		return e.Id == toDelete.Id
	})
	return toDelete, nil
}

func (f *FakeStore) GetEntry(id model.Id) (model.Entry, error) {
	idx := slices.IndexFunc(f.Entries, func(e model.Entry) bool {
		return id == e.Id
	})
	if idx < 0 {
		return model.Entry{}, store.ErrNotFound{}
	}
	return f.Entries[id], nil
}

func (f *FakeStore) ListEntries() ([]model.Entry, error) {
	return f.Entries, nil
}

func (f *FakeStore) UpdateEntry(new model.Entry) (model.Entry, error) {
	idx := slices.IndexFunc(f.Entries, func(e model.Entry) bool {
		return new.Id == e.Id
	})
	if idx < 0 {
		return model.Entry{}, store.ErrNotFound{}
	}

	f.Entries[idx] = new
	return new, nil
}
