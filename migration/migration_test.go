package migration

import "testing"

func TestLoop(t *testing.T) {
	t.Log("m1 -> m2, m2 -> m3, m3 -> m4, m4 -> m2: loop m2->m3->m4->m2")

	m1 := Meta{Name: "m1"}
	m2 := Meta{Name: "m2"}
	m3 := Meta{Name: "m3"}
	m4 := Meta{Name: "m4"}

	m1.AddRequirement(&m2)
	m2.AddRequirement(&m3)
	m3.AddRequirement(&m4)
	m4.AddRequirement(&m2)

	if err := m1.checkLoop(4, 0); err == nil {
		t.Error("Loop not found but it exists!!!")
	}

	m := Migrations{Items: []*Meta{&m1, &m2, &m3, &m4}}

	if err := m.CheckLoop(); err == nil {
		t.Error("Loop not found but it exists!!!")
	}
}

func TestFindTops(t *testing.T) {
	t.Log("m4 -> [m3,m2], m2 -> m1, m3 -> m1")

	m1 := Meta{Name: "m1"}
	m2 := Meta{Name: "m2"}
	m3 := Meta{Name: "m3"}
	m4 := Meta{Name: "m4"}

	m4.AddRequirement(&m3)
	m4.AddRequirement(&m2)
	m2.AddRequirement(&m1)
	m3.AddRequirement(&m1)

	migrations := Migrations{Items: []*Meta{&m1, &m2, &m3, &m4}}

	tops := migrations.FindTops()

	if len(tops.Items) != 1 || (len(tops.Items) == 1 && tops.Items[0].Name != "m1") {
		t.Errorf("Only m1 has no requirements!!!")
	}

	t.Log("m1 is applied")
	m1.Applied = true

	tops = migrations.FindTops()

	if len(tops.Items) != 2 || (len(tops.Items) == 2 && tops.Items[0].Name != "m2" && tops.Items[1].Name != "m3") {
		t.Errorf("Only m2 and m3 has no unapllied requirements!!! %v", tops)
	}

	t.Log("m2 is applied")
	m2.Applied = true

	tops = migrations.FindTops()

	if len(tops.Items) != 1 || (len(tops.Items) == 1 && tops.Items[0].Name != "m3") {
		t.Errorf("Only m3 has no requirements!!!")
	}

}

func TestFindBottoms(t *testing.T) {
	t.Log("m4 -> [m3,m2], m2 -> m1, m3 -> m1")

	m1 := Meta{Name: "m1", Pending: true}
	m2 := Meta{Name: "m2", Pending: true}
	m3 := Meta{Name: "m3", Pending: true}
	m4 := Meta{Name: "m4", Pending: true}

	m4.AddRequirement(&m3)
	m4.AddRequirement(&m2)
	m2.AddRequirement(&m1)
	m3.AddRequirement(&m1)

	migrations := Migrations{Items: []*Meta{&m1, &m2, &m3, &m4}}

	bottoms := migrations.FindBottoms()

	if len(bottoms.Items) != 1 || (len(bottoms.Items) == 1 && bottoms.Items[0].Name != "m4") {
		t.Errorf("Only m4 has no child!!!")
	}
}

func TestRemoveDuplicates(t *testing.T) {
	t.Log("[m1,m2,m2,m3,m4,m2,m4]")

	m1 := Meta{Name: "m1"}
	m2 := Meta{Name: "m2"}
	m3 := Meta{Name: "m3"}
	m4 := Meta{Name: "m4"}

	migrations := Migrations{Items: []*Meta{&m1, &m2, &m2, &m3, &m4, &m2, &m4}}

	plan := migrations.RemoveDuplicates()

	if len(plan.Items) != 4 {
		t.Errorf("Duplicates found!!! %v", plan)
	} else {
		if plan.Items[0].Name != "m1" || plan.Items[1].Name != "m2" || plan.Items[2].Name != "m3" || plan.Items[3].Name != "m4" {
			t.Errorf("Duplicates found!!! %v", plan)
		}
	}

}
