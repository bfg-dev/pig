package migration

// Migrations collection of migrations
type Migrations struct {
	Items []*Meta
}

// Prepare - do preparation of all migrations
func (m *Migrations) Prepare() error {
	if err := m.CheckLoop(); err != nil {
		return err
	}
	return nil
}

//

// CheckLoop - check for loops
func (m *Migrations) CheckLoop() error {
	for _, i := range m.Items {
		if err := i.checkLoop(len(m.Items), 0); err != nil {
			return err
		}
	}

	return nil
}

// GetByName - get migration by name
func (m *Migrations) GetByName(name string) *Meta {
	for _, i := range m.Items {
		if i.Name == name {
			return i
		}
	}

	return nil
}

// GetByGITinfo - get migration by name
func (m *Migrations) GetByGITinfo(gitinfo string) *Meta {
	for _, i := range m.Items {
		if i.GITinfo == gitinfo {
			return i
		}
	}

	return nil
}

// GetByNote - get migration by name
func (m *Migrations) GetByNote(note string) *Meta {
	for _, i := range m.Items {
		if i.Note == note {
			return i
		}
	}

	return nil
}

// FindTops - find unapplied migrations without [unapplied] requrements
func (m *Migrations) FindTops() *Migrations {
	var ans Migrations

	for _, i := range m.Items {
		if i.Applied {
			continue
		}

		if len(i.getUnappliedRequirements()) == 0 {
			ans.Items = append(ans.Items, i)
		}
	}

	return &ans
}

// FindBottoms - find unapplied migrations without any [unapplied] "child"
func (m *Migrations) FindBottoms() *Migrations {
	var ans Migrations

	for _, i := range m.Items {
		if i.Pending && !i.Applied && len(i.Children) == 0 {
			ans.Items = append(ans.Items, i)
		}
	}

	return &ans
}

// RemoveDuplicates - remove duplicated (from plan)
func (m *Migrations) RemoveDuplicates() *Migrations {
	var ans Migrations

	keys := make(map[string]bool)

	for _, i := range m.Items {
		if _, v := keys[i.Name]; !v {
			keys[i.Name] = true
			ans.Items = append(ans.Items, i)
		}
	}

	return &ans
}
