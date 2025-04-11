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

## Options

- `-r, --repolimit`: Number of PRs to fetch (default: 5)
- `-c, --commentlimit`: Number of comments to fetch per PR (default: 5)
- `-s, --state`: State of PRs to fetch (open/closed)
- `-v, --reviewer`: Filter PRs by review-requested user
- `-a, --author`: Filter PRs by author
