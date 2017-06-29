package main

import "errors"

func (ms *Members) get(id string) (Member, error) {
	for _, m := range *ms {
		if m.ID == id {
			return m, nil
		}
	}
	return Member{}, errors.New("not found")
}

func (ms *Members) add(m Member) error {
	_, err := ms.get(m.ID)
	if err == nil {
		return errors.New("already exists")
	}
	*ms = append(*ms, m)
	return nil
}

func (ms *Members) update(m Member) error {
	_, err := ms.get(m.ID)
	if err != nil {
		return err
	}
	for i, member := range *ms {
		if member.ID == m.ID {
			(*ms)[i].Name = m.Name
		}
	}
	return nil
}

func (ms *Members) delete(id string) error {
	ret := Members{}
	_, err := ms.get(id)
	if err != nil {
		return err
	}
	for _, m := range *ms {
		if m.ID != id {
			ret = append(ret, m)
		}
	}
	*ms = ret
	return nil
}
