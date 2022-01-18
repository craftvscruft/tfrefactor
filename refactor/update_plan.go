package refactor

type FileUpdate struct {
	Filename   string
	BeforeText string
	AfterText  string
}

type UpdatePlan struct {
	FileUpdates []*FileUpdate
}

func (p *UpdatePlan) addFileUpdate(fileUpdate *FileUpdate) {
	p.FileUpdates = append(p.FileUpdates, fileUpdate)
}
func newUpdatePlan() UpdatePlan {
	return UpdatePlan{FileUpdates: []*FileUpdate{}}
}
