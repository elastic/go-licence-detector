# Agent guidance

## Repository purpose

See [README.md](README.md) for why this repository exists, who or what uses it, and its responsibility boundary.

## AI attribution

For every AI tool that materially contributes to code, tests, documentation, configuration, or the substance of a change:

- Resolve the tool's runtime identity and add one unformatted trailer to the commit message. Use the exact model or agent slug and reasoning effort whenever exposed; omit unavailable components rather than guessing:

  ```text
  Assisted-by: <tool name> (<most specific verified runtime identity>)
  ```

- Repeat the same trailer in the pull-request description.
- Preserve the spelling and specificity of exposed runtime values; do not shorten a specific model slug to a broader model family.
- Preserve valid tool-native attribution, such as `Made with [Cursor](https://cursor.com)` or a genuine `Co-authored-by` trailer, in addition to `Assisted-by`.
- Never invent a bot identity, model name, or email address.
- Keep trailers on their own lines without bullets, Markdown emphasis, or surrounding underscores.
- Keep the human author or committer accountable for understanding and verifying the change.

For a squash merge, verify that the final squash commit message contains every attribution trailer. GitHub may populate that message from the pull-request description, commit information, or only the pull-request title depending on repository settings, so putting attribution in the PR description improves preservation but does not guarantee it.

When preparing a commit or pull request, offer to create it with the correct attribution. If the user will create it manually, show the exact trailers to copy into both places.

## Secret handling

- Never place credentials, tokens, private keys, cookies, or production secret values in tracked files, examples, tests, prompts, logs, or generated output. This repository is **public**, so anything committed here is world-readable immediately and permanently.
- GitHub secret scanning and push protection are **enabled** on this repository. They are the enforcing controls here: push protection blocks a detected secret at push time.
- Before committing or pushing, scan your work with Gitleaks. It is not pinned in this repository, so use your own installation: `gitleaks dir --no-banner --redact=100 .`. For an incident check over history, run `GIT_CONFIG_GLOBAL=/dev/null gitleaks git --no-banner --redact=100 --log-opts=--all .` from a complete clone.
- The shared `elastic/gitleaks-hooks` pre-commit hook used elsewhere in Engineering Productivity is **not** adopted here. That hook lives in a private repository, and a public repository cannot fetch it: fork pull requests receive no secrets, and external contributors have no access. Do not work around this by embedding a token, vendoring the hook, or substituting a different Gitleaks hook.
- Treat any detected secret as exposed: stop, remove it, and arrange rotation or revocation before continuing. Do not print the value while triaging it. For a public repository, assume a committed secret is compromised the moment it is pushed.
- Never bypass push protection or secret scanning, select a GitHub bypass reason, use `--no-verify`, disable a scanner, or add an allowlist or suppression unless the user explicitly requests that exact override after reviewing the finding and consequence.
- A generic request to finish, commit, push, merge, or open a pull request does not authorize an override.
