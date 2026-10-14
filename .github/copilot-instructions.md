<!-- configure-ai-review:baseline:start -->
# High-signal code review policy

These rules apply only when reviewing a pull request. They do not apply to Copilot Chat or cloud agent tasks. When reviewing, default to leaving no comment. Comment only when the changed code introduces a concrete correctness, security, data-integrity, compatibility, reliability, rollout, or operability defect, or when changed material behavior has no test that would catch a regression (see "Missing tests for material behavior" below).

When multiple issues compete for attention, prioritize in this order: (1) concrete correctness, security, or data-integrity defects; (2) architectural or design problems with a concrete, describable consequence; (3) missing test coverage for changed behavior. Leave the highest-priority comments first.

Do not comment on wording, spelling, formatting, naming preferences, syntax, style, import ordering, or other concerns that a compiler, formatter, linter, schema validator, test, or CI check enabled in this repository already reports. Do not restate a failing automated check. Do not leave praise, summaries, optional refactors, generic best practices, or speculative suggestions.

Before leaving a comment, verify all of the following:

- The issue is introduced or made materially worse by this pull request.
- You can describe a specific failure scenario and its impact.
- The finding starts from the changed lines. You may read callers, tests, and configuration to confirm it, but do not report issues the diff does not introduce or make worse.
- The repository's current code and configuration support the claim.
- The comment identifies the affected path and gives an actionable correction or asks for concrete missing evidence.
- The expected value of the comment is high enough to interrupt a human reviewer.

If any condition is missing, remain silent. Express uncertainty only when the uncertainty itself is a material risk that the author can resolve with specific evidence. On a re-review after new commits, do not repeat a comment on a line that has not changed since the prior run.

Use the repository's `.github/skills/code-review/SKILL.md` for domain-specific review checks when it exists. Copilot review is advisory. Whether its approval counts toward required approvals is controlled by org and branch settings, not this file.

## Trust-boundary and supply-chain changes

Always inspect a diff that changes CI, build, release, dependency, credential,
network, or security configuration. Note: Copilot code review does not inspect
lockfiles, `requirements.txt`, or generated files — leave those to Dependabot
and human review. Leave a review comment when the change introduces or materially
broadens a capability without a concrete, repository-supported need and an
appropriate control. In particular, flag:

- broader token, role, cloud, repository, or secret permissions than the changed
  operation requires;
- new collection, copying, logging, caching, artifact upload, or persistence of
  credentials or secret-bearing worker state;
- privileged workflows that check out, download, or execute untrusted pull-
  request code, artifacts, build scripts, hooks, or dependency lifecycle scripts;
- new outbound destinations, telemetry, uploads, remote shells, or network
  services unrelated to the repository's documented responsibility;
- downloaded or remotely referenced code, actions, reusable workflows, images,
  plugins, tools, or binaries that are executed without immutable identity and
  integrity verification;
- new dependencies or build inputs that are not resolved immutably;
- weakened or bypassed tests, approval gates, secret scanning, signature or
  checksum verification, TLS verification, sandboxing, or least-privilege
  defaults;
- shared caches, workspaces, or non-ephemeral runners that allow an untrusted or
  earlier build to influence a privileged or later build;
- unexplained executable binaries, encoded or obfuscated payloads, generated
  code without a documented generator, or functionality unrelated to the
  repository's stated purpose;
- untrusted input, such as a pull-request title, branch name, or issue body,
  placed directly into a `run:` script through a `${{ }}` expression instead
  of through an environment variable.

For GitHub Actions, require third-party actions and reusable workflows to be
pinned to a full commit SHA from the intended upstream repository. Require
container images by digest when they execute in CI, build, release, or
production paths. Require downloaded executables or scripts to be pinned to an
immutable version and verified before execution with an expected cryptographic
digest, or with a signature or attestation from an expected signer. HTTPS or a
checksum downloaded from the same mutable source is not sufficient integrity
evidence. Reject pipe-to-shell execution of unverified remote content.

For package-manager dependencies, require the repository's lockfile or equivalent
immutable resolution to be committed and enforced. Do not comment on a version
range merely because it is a range when an enforced lockfile pins the resolved
artifact. Do comment when a new executable build or runtime dependency can float
without reviewed immutable resolution.

Describe the exact new trust path, the realistic abuse or compromise scenario,
and the smallest corrective control. Do not allege malicious intent; when code
appears unrelated or concealed, ask for concrete purpose and provenance and
recommend removal when neither is established.

## Missing tests for material behavior

Check whether changed behavior has meaningful deterministic coverage. Comment
when the diff adds or changes a material correctness, security, compatibility,
failure-handling, state-transition, permission, parsing, retry, concurrency,
build, or release behavior and no existing or new test would fail if that
behavior regressed. A bug fix should normally include a regression test that
fails for the original defect.

Before commenting, inspect the relevant existing tests, including higher-level
and generated-output tests. Do not request tests for documentation-only changes,
mechanical refactors with unchanged behavior, generated files whose source and
generation are already tested, or behavior already exercised at an appropriate
layer. Do not demand a unit test when an existing integration or end-to-end test
more directly covers the contract.

A missing-test comment must name the unverified behavior and a concrete scenario
and assertion that would catch a realistic regression. Do not leave a generic
"add tests" comment or prescribe test-file structure. When deterministic testing
is genuinely impractical, ask for specific reproducible verification evidence
and the reason it cannot be automated.

<!-- configure-ai-review:baseline:end -->

# Repository-specific review guidance

## Language and toolchain

Go module `go.elastic.co/go-licence-detector`; the Go version is pinned in `go.mod` and CI reads it through `go-version-file`. The tool reads `go list -m -json all` output and generates licence notices and dependency lists for Elastic projects. Compiler and `go test` results are authoritative.

## Focus areas

- Changes to `assets/rules.json` (allowlist or maybelist): `Rules.IsAllowed` in `detector/rules.go` treats both lists as allowed, so a new entry makes every dependency with that licence pass for users of the embedded default rules (callers can supply their own with `-rules`). A public policy reference or maintainer confirmation is enough as the basis.
- Changes to `assets/licenses.db`: it is an embedded binary produced by the `license_serializer` tool, so a hand-edited or unexplained replacement can change which licences the classifier detects. Ask for the source and regeneration steps in the PR description.
- Changes to override handling (`mkDepInfo` in `detector/detector.go`, `LoadOverrides` in `dependency/dependency.go`): overrides are matched by module name and win over module data through `coalesce`. A changed key match or precedence applies an override to the wrong module or drops it, so NOTICE files state a wrong licence type or URL. Check that a test covers the changed precedence.
- Changes to licence file path handling: `licenceTextOverrideFile` must stay inside the overrides file's directory (`LoadOverrides`), and `licenceFile` must stay inside the module directory (`doDetectLicences`). Review whether each containment is preserved; a path that escapes it can read arbitrary files into a NOTICE. Do not require one specific function.
- Changes to direct and indirect dependency selection (`parseDependencies`, `-includeIndirect`): a dependency that is silently dropped or wrongly classed as indirect leaves a missing entry in a legal notice.
- Changes to `replace` handling in `mkDepInfo`: the intended behaviour is that a remote replacement uses the replacement's path, version, and metadata, while a local replacement (path starting with `.`) keeps the original module's name and version and reads the replacement directory. Flag changes that break either contract, not the existing split.
- Changes to templates in `example/templates` or the `render` package: a change can alter the generated NOTICE and dependency list for every consumer. Ask for sample output or a render test only when the change creates a concrete regression risk, such as new logic or a changed data field.
- Changes to `validate`: it makes network requests to each dependency URL, so it must stay explicit behind `-validate` and must not become a required network dependency of a normal run or build.
- Changes to CLI flags or defaults in `main.go`: the flags are documented in `README.md` and used by downstream scripts. Flag a renamed or removed flag, or a changed default that alters output, that has no README update or compatibility note.

## Skip list

- `go test ./...` and `go build` run in `.github/workflows/build.yml`; do not comment on compile errors or failing tests.
- Renovate (`renovate.json`) keeps dependency versions up to date; do not comment on dependency version ranges or on routine version bumps it proposes.
