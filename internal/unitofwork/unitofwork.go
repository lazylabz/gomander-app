// Package unitofwork is how an operation that writes through more than one
// repository lands as a single piece. Each repository keeps its own writes
// consistent, but nothing spanned them: importing a Project wrote the Project,
// then its Commands, then its Command Groups, and a failure partway through
// left the user with the half that had already landed. A Unit of Work hands the
// operation repositories bound to one transaction, so atomicity is a property
// of the operation rather than of whichever repository it happened to use.
package unitofwork

import (
	commanddomain "gomander/internal/command/domain"
	commandgroupdomain "gomander/internal/commandgroup/domain"
	projectdomain "gomander/internal/project/domain"
)

// Repositories are the repositories a Unit of Work binds to its transaction.
// They answer exactly like the ones wired into the app, so a caller reads and
// writes through them the way it would outside one.
type Repositories struct {
	Projects      projectdomain.Repository
	Commands      commanddomain.Repository
	CommandGroups commandgroupdomain.Repository
}

// UnitOfWork runs change against Repositories bound to a single transaction:
// every write inside it lands, or none of them does. An error returned by
// change undoes the writes it had already made and is reported back.
type UnitOfWork interface {
	Do(change func(Repositories) error) error
}
