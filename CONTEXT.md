# Gomander

A desktop application for launching, monitoring and organising the shell commands
a developer runs on a project.

## Language

**Command**:
A shell command line saved under a Project, together with the working directory
it runs in and the patterns that mark its output as failing.
_Avoid_: script, task, job

**Command Group**:
A named bundle of a Project's Commands, launched as one. A Group ceases to exist
when its last Command is removed.
_Avoid_: bundle, batch, collection

**Project**:
The unit the user opens to work: a name, a working directory, and the Commands
and Command Groups that belong to it.
_Avoid_: workspace, folder, repo

**Opened Project**:
The single Project the user currently has open. Having none open is a legitimate
state, not a failure.
_Avoid_: current project, active project, selected project

**Execution Environment**:
The paths and base working directory a Command runs in, resolved from the Opened
Project before the Runner is involved. Bare "execution" is ambiguous — it reads
as the act of running; use this term or say Process.

**Runner**:
The engine that spawns a Command as an operating-system process and streams its
output back.
_Avoid_: executor, launcher, engine

**Process**:
One running instance of a Command, owned by the Runner from spawn to
termination.
_Avoid_: execution, run, instance

**Running Commands**:
The set of Commands the Runner currently holds a Process for. It is what answers
whether a Command is running and how many of a Command Group's Commands are; no
consumer decides that for itself.
_Avoid_: active commands, live commands, running state

**Error Pattern**:
A substring that marks a line of a Command's output as failing when the line
contains it. A Command carries a list of them.
_Avoid_: error matcher, error regex — the match is a plain substring, not a regex

**Cascade**:
What a set of Command Groups becomes once a Command they held is gone: the ones
still standing, and the ones that ceased to exist because that Command was their
last.

**Position**:
Where a thing sits among its siblings — a Command among its Project's Commands,
a Command Group among its Project's Groups, a Command among the Commands its
Command Group holds. Positions are dense: 0, 1, 2 … with no gaps.
_Avoid_: order, index, sort order

**Unit of Work**:
One transaction spanning several repositories: everything an operation writes
inside it lands, or none of it does. It is what makes atomicity a property of
the operation rather than of a single repository.
_Avoid_: transaction (that is the database's word for the mechanism)
