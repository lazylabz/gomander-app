package handlers_test

type FakeEvent struct{}

func (FakeEvent) GetName() string { return "fake" }
