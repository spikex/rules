#!/bin/bash

# Script to fix commits made with cursoragent@cursor.com email
# Replaces cursoragent@cursor.com with current git email
# Renames current branch from cursor/name-of-feature to nate/name-of-feature

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository. Please run this script from a git repository."
    exit 1
fi

# Get current git config
CURRENT_NAME=$(git config user.name)
CURRENT_EMAIL=$(git config user.email)

print_info "Current git config:"
print_info "  Name: $CURRENT_NAME"
print_info "  Email: $CURRENT_EMAIL"

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)
print_info "Current branch: $CURRENT_BRANCH"

# Check if current branch starts with cursor/
if [[ "$CURRENT_BRANCH" == cursor/* ]]; then
    NEW_BRANCH=$(echo "$CURRENT_BRANCH" | sed 's/^cursor\//nate\//')
    print_info "Will rename branch: $CURRENT_BRANCH -> $NEW_BRANCH"
    RENAME_BRANCH=true
else
    print_info "Current branch does not start with 'cursor/', no branch renaming needed."
    RENAME_BRANCH=false
fi

# Find commits with cursoragent@cursor.com
print_info "Finding commits with cursoragent@cursor.com on current branch..."

CURSOR_COMMITS=$(git log --oneline --author="cursoragent@cursor.com" "$CURRENT_BRANCH")

if [ -z "$CURSOR_COMMITS" ]; then
    print_warning "No commits found with cursoragent@cursor.com email on current branch."
    COMMIT_COUNT=0
else
    print_info "Found the following commits with cursoragent@cursor.com on current branch:"
    echo "$CURSOR_COMMITS"
    echo ""

    # Count commits
    COMMIT_COUNT=$(echo "$CURSOR_COMMITS" | wc -l)
    print_info "Total commits to fix: $COMMIT_COUNT"
fi

# Confirm before proceeding
if [ $COMMIT_COUNT -gt 0 ] || [ "$RENAME_BRANCH" = true ]; then
    print_warning "This will:"
    if [ $COMMIT_COUNT -gt 0 ]; then
        print_warning "  - Change the author email from 'cursoragent@cursor.com' to '$CURRENT_EMAIL' for commits on current branch"
    fi
    if [ "$RENAME_BRANCH" = true ]; then
        print_warning "  - Rename current branch from '$CURRENT_BRANCH' to '$NEW_BRANCH'"
    fi
    
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo ""

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Operation cancelled."
        exit 0
    fi
else
    print_warning "No commits or branch renaming needed."
    exit 0
fi

# Create backup branch
BACKUP_BRANCH="backup-$(date +%Y%m%d-%H%M%S)"
print_info "Creating backup branch: $BACKUP_BRANCH"
git branch "$BACKUP_BRANCH"

# Use git filter-branch to change commits on current branch only (only if there are commits to fix)
if [ $COMMIT_COUNT -gt 0 ]; then
    print_info "Updating commit author emails on current branch..."

    git filter-branch --env-filter '
if [ "$GIT_AUTHOR_EMAIL" = "cursoragent@cursor.com" ]; then
    export GIT_AUTHOR_EMAIL="'"$CURRENT_EMAIL"'"
    export GIT_AUTHOR_NAME="'"$CURRENT_NAME"'"
fi
if [ "$GIT_COMMITTER_EMAIL" = "cursoragent@cursor.com" ]; then
    export GIT_COMMITTER_EMAIL="'"$CURRENT_EMAIL"'"
    export GIT_COMMITTER_NAME="'"$CURRENT_NAME"'"
fi
' --tag-name-filter cat -- "$CURRENT_BRANCH"

    # Clean up the backup refs
    print_info "Cleaning up backup refs..."
    git for-each-ref --format="%(refname)" refs/original/ | xargs -n 1 git update-ref -d

    print_success "Successfully updated commits with cursoragent@cursor.com on current branch!"
fi

# Rename current branch from cursor/ to nate/
if [ "$RENAME_BRANCH" = true ]; then
    print_info "Renaming current branch from '$CURRENT_BRANCH' to '$NEW_BRANCH'..."
    
    # Rename the current branch locally only
    git branch -m "$NEW_BRANCH"
    
    print_success "Successfully renamed branch: $CURRENT_BRANCH -> $NEW_BRANCH"
fi

print_info "Backup branch created: $BACKUP_BRANCH"
print_warning "If you're satisfied with the changes, you may want to delete the backup branch: git branch -d $BACKUP_BRANCH"
print_warning "If you need to push these changes, you may need to force push: git push --force-with-lease"