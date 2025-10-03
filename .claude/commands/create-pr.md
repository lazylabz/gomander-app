---
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git branch:*), Bash(git push:*), Bash(gh pr create:*)
description: Create a PR following the project template with iterative description editing
---

Create a pull request following these steps:

1. **Analyze the current branch changes:**
   - Run `git status` to see uncommitted changes
   - Run `git diff main...HEAD --stat` to see all changes from main
   - Run `git log main..HEAD` to see all commits in this branch
   - Check if branch tracks remote with `git rev-parse --abbrev-ref --symbolic-full-name @{u}`

2. **Draft PR content based on `.github/pull_request_template.md`:**
   - **Title**: Follow conventional commit format (feat/fix/docs/refactor/test/chore)
   - **Description sections**:
     - Follow `.github/pull_request_template.md`

3. **Present the draft to me:**
   - Show the proposed PR title
   - Show the complete PR description
   - Wait for my feedback

4. **Iterate based on feedback:**
   - If I request edits, update specific sections
   - Show the updated version
   - Repeat until I approve

5. **Submit the PR:**
   - Only after I explicitly approve
   - Push branch to remote if needed (`git push -u origin <branch>`)
   - Create PR with `gh pr create --title "..." --body "..."`
   - Use HEREDOC for the body to ensure proper formatting

**Important:**
- Do NOT push or create the PR until I explicitly approve
- Ensure the title follows conventional commit format
- Fill out all relevant sections from the template
- Check all commits in the branch, not just the latest one
