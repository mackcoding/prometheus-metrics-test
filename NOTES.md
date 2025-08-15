# My Setup
I’m running this in a kind cluster on WSL, with Docker Desktop on Windows 11 and kind running in Debian within WSL.

# Running the project

# Adjustments
- When I ran kubectl apply -f postgres.yaml, it failed due to a cgroup version mismatch, likely caused by a WSL2 kernel issue. I resolved the problem by commenting out the resource specifications in postgres.yaml.
- I'm deploying my custom app to the same namespace as postgres, I typically wouldn't do that
