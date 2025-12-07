# Setting Up Remote Repository

## Option 1: GitHub

```bash
# Create a new repository on GitHub (e.g., github.com/davidl71/devwisdom-go)
# Then add remote:
cd /Users/davidl/Projects/devwisdom-go
git remote add origin https://github.com/davidl71/devwisdom-go.git
git branch -M main
git push -u origin main
```

## Option 2: GitLab

```bash
# Create a new repository on GitLab
git remote add origin https://gitlab.com/davidl71/devwisdom-go.git
git branch -M main
git push -u origin main
```

## Option 3: Self-hosted

```bash
# Add your self-hosted git remote
git remote add origin <your-git-server-url>
git branch -M main
git push -u origin main
```

## Verify Remote

```bash
git remote -v
git status
```
