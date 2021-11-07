# Tasks

Tasks are scheduled in UTC.

Example natural language temporal expressions:

- `"every 2nd hour"`
- `"Tue/Thu at 4am"`

See https://github.com/pnelson/te for documentation on valid formats.

A task will not be run if its previous invocation is still running at the next
scheduled run time.
