# Frontend coverage is not uploaded to Codecov

`codecov.yml` sets a repo-wide project target of 80%, measured over the Go
backend only. The frontend suite is young and covers a small fraction of
`cmd/gomander/frontend/src`, so uploading it would drag the combined figure under
the target and fail every pull request — including ones that add frontend tests.

This is a sequencing decision, not a judgement that frontend coverage does not
matter: `pnpm run test:cov` works locally and is worth running. Revisit once the
suite is large enough to clear the target on its own, rather than by lowering it.
