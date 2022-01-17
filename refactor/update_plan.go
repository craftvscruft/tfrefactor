package refactor

type FileUpdate struct {
	beforeText string
	afterText  string
	filename   string
}

type UpdatePlan struct {
	fileUpdates []*FileUpdate
}

func (p *UpdatePlan) addFileUpdate(fileUpdate *FileUpdate) {
	p.fileUpdates = append(p.fileUpdates, fileUpdate)
}
func newUpdatePlan() UpdatePlan {
	return UpdatePlan{fileUpdates: []*FileUpdate{}}
}
