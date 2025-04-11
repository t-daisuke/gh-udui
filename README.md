# gh-utui

this is git hub extensiton for comments, issues, and reviews Update detection on Text User Interface

## Installation

```bash
gh extension install t-daisuke/gh-utui
```

## Usage

```bash
# View PR comments (default: shows 5 latest PRs with 5 latest comments each)
gh utui

# Customize the number of PRs and comments
gh utui -r 10 -c 3  # Show 10 PRs with 3 comments each

# Filter by state
gh utui -s open    # Show only open PRs
gh utui -s closed  # Show only closed PRs

# Filter by reviewer
gh utui -v "@me"   # Show PRs where you are requested as a reviewer

# Filter by author
gh utui -a "@me"   # Show your PRs
```

## Example Output

```
> gh utui
Fetching 5 PRs and 5 comments, repoState: all, reviewer: any, author: @me 

org/project#123 Fix plugin URL environment variable [OPEN] (updated: 2023-04-10T07:28:09Z)
reviewer1 (2023-04-10 07:26:42):
✅ [State: APPROVED]
LGTM! Ready to merge.

reviewer1 (2023-04-10 07:22:37):
💬 [State: COMMENTED]
extraEnv is already defined, please use that as discussed
---
org/project#122 Enable API endpoint for Dify plugin [OPEN] (updated: 2023-04-10T06:42:45Z)
reviewer1 (2023-04-10 06:40:39):
✅ [State: APPROVED]
Looks good
https://github.com/example/project/blob/main/api/.env.example#L455
---
org/admin-tools#87 Refactoring authentication logic [OPEN] (updated: 2023-04-10T02:44:48Z)
reviewer2 (2023-04-10 02:42:10):
❌ [State: CHANGES_REQUESTED]
Please add tests for the new authentication flow

author (2023-04-10 02:30:15):
💬 [State: COMMENTED]
I've updated the PR based on previous feedback
---
org/project#121 Revert "Fixed plugin endpoint issue" [CLOSED] (updated: 2023-04-09T10:47:04Z)
reviewer1 (2023-04-09 10:45:55):
✅ [State: APPROVED]
The environment variable didn't exist, reverting is the right approach
---
org/project#120 Fix plugin endpoint URL format [OPEN] (updated: 2023-04-09T09:40:33Z)
reviewer3 (2023-04-09 09:40:33):
🚫 [State: DISMISSED]
Dismissing my review as we're taking a different approach

reviewer1 (2023-04-09 09:40:32):
✅ [State: APPROVED]
Values look good 👍 
https://github.com/example/project/blob/master/charts/values.yaml#L606

reviewer2 (2023-04-09 09:38:24):
❌ [State: CHANGES_REQUESTED]
The URL format doesn't match our standards, please reference the documentation
```

### Legend

Review states are indicated by emojis:
- ✅ [APPROVED] - The changes have been approved
- ❌ [CHANGES_REQUESTED] - Changes have been requested before approval
- 💬 [COMMENTED] - Comments without explicit approval or request for changes
- 🚫 [DISMISSED] - The review has been dismissed

## Options

- `-r, --repolimit`: Number of PRs to fetch (default: 5)
- `-c, --commentlimit`: Number of comments to fetch per PR (default: 5)
- `-s, --state`: State of PRs to fetch (open/closed)
- `-v, --reviewer`: Filter PRs by review-requested user
- `-a, --author`: Filter PRs by author
