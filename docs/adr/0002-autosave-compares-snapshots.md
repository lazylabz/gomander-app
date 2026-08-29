# Autosaved forms compare snapshots, not dirty fields

`useAutosavedForm` decides whether a form changed by comparing a JSON snapshot of
its values against the last values it saved.

The obvious alternative is `formState.dirtyFields`, which react-hook-form already
maintains. It cannot work here: dirtiness is measured against the form's
*defaults*, and a field only leaves `dirtyFields` when its value matches that
default again. Autosave needs the comparison against the *last saved* values, so
once a save has moved those away from the defaults, a field returning to what was
last saved stays dirty for good. Making it work would require a `form.reset` after
every save — which re-keys any `useFieldArray` on screen and steals focus from the
user mid-typing, on a form that saves while they type.

Recorded because `dirtyFields` appears nowhere in the codebase: there is no code
to read that explains why the idiomatic API was passed over.
