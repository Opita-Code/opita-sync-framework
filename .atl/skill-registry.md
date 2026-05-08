# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

See `_shared/skill-resolver.md` for the full resolution protocol.

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| responsive design, mobile layouts, breakpoints, viewport adaptation, cross-device compatibility | adapt | ~/.agents/skills/adapt/SKILL.md |
| animation, transitions, micro-interactions, motion design, hover effects, making UI feel alive | animate | ~/.agents/skills/animate/SKILL.md |
| accessibility check, performance audit, technical quality review | audit | ~/.agents/skills/audit/SKILL.md |
| SAM deploy, CloudFormation, Route53, ACM, SSM, S3/CloudFront, infrastructure, AWS | aws-cli | ~/.config/opencode/skills/aws-cli/SKILL.md |
| design looks bland, generic, too safe, lacks personality, wants more visual impact | bolder | ~/.agents/skills/bolder/SKILL.md |
| creating a pull request, opening a PR, preparing changes for review | branch-pr | ~/.config/opencode/skills/branch-pr/SKILL.md |
| chained PRs, stacked PRs, PR exceeds 400 lines, reviewable slices | chained-pr | ~/.config/opencode/skills/chained-pr/SKILL.md |
| confusing text, unclear labels, bad error messages, hard-to-follow instructions, UX writing | clarify | ~/.agents/skills/clarify/SKILL.md |
| tareas ClickUp, comments, estados, listas, espacios, URLs/IDs de ClickUp | clickup | ~/.config/opencode/skills/clickup/SKILL.md |
| guides, READMEs, RFCs, onboarding docs, architecture docs, review-facing documentation | cognitive-doc-design | ~/.config/opencode/skills/cognitive-doc-design/SKILL.md |
| design looking gray, dull, lacking warmth, needing more color, vibrant palette | colorize | ~/.agents/skills/colorize/SKILL.md |
| drafting feedback, review comments, maintainer replies, Slack messages, GitHub comments | comment-writer | ~/.config/opencode/skills/comment-writer/SKILL.md |
| review, critique, evaluate, give feedback on a design or component | critique | ~/.agents/skills/critique/SKILL.md |
| add polish, personality, animations, micro-interactions, delight, make UI feel fun | delight | ~/.agents/skills/delight/SKILL.md |
| simplify, declutter, reduce noise, remove elements, make UI cleaner | distill | ~/.agents/skills/distill/SKILL.md |
| find a skill, how do I do X, is there a skill for, extending capabilities | find-skills | ~/.agents/skills/find-skills/SKILL.md |
| runs, checks, logs, reruns, artifacts, diagnosticar CI fallida en GitHub | gh-actions | ~/.config/opencode/skills/gh-actions/SKILL.md |
| login/logout, auth status, token, scopes, permisos, selección de host en GitHub CLI | gh-auth | ~/.config/opencode/skills/gh-auth/SKILL.md |
| crear, buscar, editar, etiquetar, triagear, cerrar o reabrir issues en GitHub | gh-issue | ~/.config/opencode/skills/gh-issue/SKILL.md |
| crear, abrir, editar, preparar, publicar o mergear un pull request como autor | gh-pr | ~/.config/opencode/skills/gh-pr/SKILL.md |
| crear tags, preparar release notes, adjuntar assets o publicar releases en GitHub | gh-release | ~/.config/opencode/skills/gh-release/SKILL.md |
| revisar un PR, aprobar, solicitar cambios o responder comentarios de revisión | gh-review | ~/.config/opencode/skills/gh-review/SKILL.md |
| buscar correos, leer la bandeja, organizar Gmail, archivar mensajes, IMAP | gmail-imap | ~/.config/opencode/skills/gmail-imap/SKILL.md |
| writing Go tests, using teatest, adding test coverage | go-testing | ~/.config/opencode/skills/go-testing/SKILL.md |
| build web components, pages, artifacts, posters, applications | impeccable | ~/.agents/skills/impeccable/SKILL.md |
| creating a GitHub issue, reporting a bug, requesting a feature | issue-creation | ~/.config/opencode/skills/issue-creation/SKILL.md |
| "judgment day", review adversarial, dual review, doble review, juzgar | judgment-day | ~/.config/opencode/skills/judgment-day/SKILL.md |
| layout feeling off, spacing issues, visual hierarchy, crowded UI, alignment problems | layout | ~/.agents/skills/layout/SKILL.md |
| slow, laggy, janky, performance, bundle size, load time, faster rendering | optimize | ~/.agents/skills/optimize/SKILL.md |
| wow, impress, go all-out, shaders, spring physics, 60fps animations | overdrive | ~/.agents/skills/overdrive/SKILL.md |
| polish, finishing touches, pre-launch review, something looks off, good to great | polish | ~/.agents/skills/polish/SKILL.md |
| prospeccion, Opita Prospección, leads Neiva, pipeline prospección | prospeccion | ~/.config/opencode/skills/prospeccion/SKILL.md |
| too bold, too loud, overwhelming, aggressive, garish, calmer aesthetic | quieter | ~/.agents/skills/quieter/SKILL.md |
| plan UX and UI for a feature before writing code | shape | ~/.agents/skills/shape/SKILL.md |
| create a new skill, add agent instructions, document patterns for AI | skill-creator | ~/.config/opencode/skills/skill-creator/SKILL.md |
| update skills, skill registry, actualizar skills, update registry | skill-registry | ~/.config/opencode/skills/skill-registry/SKILL.md |
| fonts, type, readability, text hierarchy, sizing looks off, typography | typeset | ~/.agents/skills/typeset/SKILL.md |
| implementing a change, preparing commits, splitting PRs, chained or stacked PRs | work-unit-commits | ~/.config/opencode/skills/work-unit-commits/SKILL.md |

## Compact Rules

Pre-digested rules per skill. Delegators copy matching blocks into sub-agent prompts as `## Project Standards (auto-resolved)`.

### go-testing
- Use table-driven tests: `tests := []struct{name, input, expected, wantErr}`, iterate with `t.Run(tt.name, func...)`
- Mock dependencies with interfaces; use `httptest` for HTTP handlers, `go-sqlmock` for PostgreSQL
- Use `t.TempDir()` for temp files; `t.Helper()` in test helpers
- Test both success and error paths; use `t.Fatalf` for setup failures, `t.Errorf` for assertion failures
- Coverage command: `go test -cover ./...`; skip integration tests with `-short`
- Package naming: use `package foo_test` for black-box tests (in same folder as `foo`)

### work-unit-commits
- Commit by deliverable work unit, NOT by file type (no "add models" then "add tests")
- Tests belong in the same commit as the behavior they verify
- Docs belong with the user-visible change they explain
- Each commit must tell a story: one clear purpose, independently buildable
- If PR approaches 400 changed lines, promote commits into chained PRs
- Use conventional commit messages: `feat(domain): what changed`

### chained-pr
- MUST split when PR exceeds 400 changed lines (additions + deletions), unless maintainer-approved `size:exception`
- Each chained PR must be autonomous: CI green, one deliverable work unit, reviewable alone
- State start/end/prior/next dependencies in every PR description
- If SDD forecasts >400-line workload, honor `delivery_strategy`: ask, auto-chain, or size:exception
- Child PRs in Feature Branch Chain target the immediate previous PR branch
- Always include a dependency diagram marking the current PR

### branch-pr
- Every PR MUST link an approved issue (`Closes/Fixes/Resolves #N`)
- PR title MUST use conventional commit format: `type(scope): description`
- Branch naming: `type/description` — lowercase, no spaces, `a-z0-9._-` only
- Add exactly one `type:*` label to the PR
- Automated checks (issue-link + conventional title) must pass

### issue-creation
- Create issues before starting work; link approved issues in PRs
- Use issue templates: bug_report.yml, feature_request.yml
- Issues need `status:approved` label before a PR can be opened

### gh-pr
- Use `gh pr create` with body passed as PowerShell here-string: `gh pr create --title "..." --body @'...'@`
- Create branch, push with `-u`, then create PR
- Verify issue linkage and conventional title before opening

### gh-issue
- Use `gh issue create` with labels, assignees, milestones
- Search issues: `gh issue list --label "..." --state open`
- Edit and close: `gh issue edit`, `gh issue close`

### gh-review
- Use `gh pr review` with `--approve`, `--request-changes`, or `--comment`
- Add inline comments: `gh pr comment`
- Check review status: `gh pr view --json reviews`

### gh-actions
- List runs: `gh run list --workflow "name"`
- View logs: `gh run view RUN_ID --log`
- Rerun: `gh run rerun RUN_ID`

### cognitive-doc-design
- Progressive disclosure: summary first, details after, use collapsible sections
- Chunk information into scannable sections with clear headers
- Use signposting: "you'll learn X, Y, Z" at top of sections
- Prefer tables and checklists over prose for reference material
- Recognition over recall: include concrete examples, CLI commands

### comment-writer
- Write warm, direct, human comments — no corporate-speak
- Be specific about what was good and what needs work
- Include code suggestions with `suggestion` blocks when relevant
- Keep feedback constructive: problem → impact → solution

### adapt
- Design mobile-first: start at smallest viewport, add complexity for larger
- Use fluid layouts (CSS Grid, flexbox), avoid fixed pixel widths
- Touch targets: minimum 44x44px for interactive elements
- Test at 320px, 768px, 1024px, 1440px breakpoints

### animate
- Prefer `transform` and `opacity` for animations (GPU-composited, no layout thrash)
- Keep animations 200-400ms; micro-interactions 100-200ms
- Use `prefers-reduced-motion` media query to respect accessibility
- Easing: ease-out for entering, ease-in for exiting

### audit
- Check accessibility: color contrast, keyboard navigation, ARIA labels
- Check performance: Largest Contentful Paint, Cumulative Layout Shift, bundle size
- Check theming: dark mode, consistent token usage
- Generate scored report with P0-P3 severity and actionable plan

### overdrive
- Use CSS `@property` for typed custom properties with spring animations
- Prefer scroll-driven animations with `animation-timeline: scroll()`
- Keep 60fps: use `will-change` sparingly, measure with Performance API
- Spring physics: use `linear()` easing or JS spring libraries

### polish
- Fix inconsistent spacing (use design system tokens, not arbitrary values)
- Align elements to a 4px or 8px grid
- Ensure consistent border-radius, shadow, and typography scale
- Review hover, focus, active, disabled states for all interactive elements

### quiet
- Reduce color saturation and contrast of non-primary elements
- Use lighter font weights for secondary text
- Increase whitespace around primary content
- Mute background accents; keep CTAs distinctive

### bolder
- Increase color contrast and saturation on primary elements
- Use heavier font weights for headings
- Add distinctive borders, shadows, or accents to key sections
- Make CTAs immediately recognizable with size and color

### colorize
- Add accent colors to break monochromatic monotony
- Use a 60-30-10 rule: 60% neutral, 30% primary, 10% accent
- Ensure all colors meet WCAG AA contrast (4.5:1 for text)
- Define colors as design tokens for consistency

### distill
- Remove redundant elements: if two things do the same job, keep one
- Progressive disclosure: hide advanced options behind "show more"
- Use default sensible values; let users override when needed
- Every element must justify its existence with a user need

### delight
- Add subtle easter eggs or personality moments (not distracting)
- Use delightful loading states (meaningful skeletons, not spinners)
- Confetti or subtle celebration on meaningful achievements
- Keep delight moments brief and skip-able

### typeset
- Define a type scale with 2-3 font sizes minimum (body, heading, display)
- Use a maximum of 2 font families (one for headings, one for body)
- Line height: 1.5-1.75 for body text, 1.1-1.3 for headings
- Line length: 45-75 characters for body text

### layout
- Use a consistent grid system (12-column or custom)
- Maintain visual hierarchy: primary content gets most space
- Group related elements with proximity; separate unrelated with whitespace
- Use alignment consistently: left-align text, consistent padding

### clarify
- Write error messages that say what happened, why, and what to do next
- Labels should describe outcome, not implementation: "Save changes" not "Submit"
- Use active voice, present tense, and direct language
- Avoid jargon; use terms your users actually use

### shape
- Run structured discovery before coding: goals, constraints, personas
- Produce a design brief with: problem, users, constraints, UX strategy
- Establish design direction before implementation begins

### critique
- Evaluate visual hierarchy: does the eye flow naturally to primary actions?
- Score information architecture: can users find what they need in 3 clicks?
- Test with personas: how would different user types interact?
- Score emotional resonance: does the design feel trustworthy, exciting, calm?

### skill-creator
- Create SKILL.md with YAML frontmatter: name, description, license, metadata
- Follow Agent Skills spec: `~/.config/opencode/skills/skill-creator/`
- Description must include "Trigger:" text for auto-activation
- Place skill in appropriate directory for the target tool

### gmail-imap
- Use IMAP with imapflow for secure connection
- Search: `imap.search({ criteria })` with IMAP search syntax
- Decode MIME with mailparser for HTML/plain text extraction
- Move/archive: `imap.messageMove(uid, targetFolder)`

### aws-cli
- Use `aws cloudformation describe-stacks` before any mutation
- Prefer `--no-execute-changeset` for preview of CloudFormation changes
- Always verify Route53/DNS changes with `aws route53 list-resource-record-sets`

### clickup
- Resolve task by URL: extract ID from `/v/vt/...` or `/t/...` patterns
- List tasks: `clickup.tasks.list(spaceId, { statuses: [...] })`
- Create tasks: `clickup.tasks.create(spaceId, { name, description })`
- Comments: `clickup.comments.list(taskId)` and `clickup.comments.create(taskId, { text })`

### prospeccion
- Pipeline: sourcing → auditoría → enriquecimiento → mockups → video → outreach
- Focus on negocios en Neiva, Huila
- Auditoría de sitios web con scoring automático
- Generar landing pages como mockups previo a outreach

### find-skills
- Search skill directories for matching SKILL.md files by trigger text
- Install via skill-creator patterns
- Present available skills with trigger phrases and brief descriptions

### judgment-day
- Launch two independent blind judge sub-agents simultaneously
- Synthesize findings after both complete
- Apply fixes and re-judge until both pass (max 2 iterations)
- Escalate on persistent disagreements

## Project Conventions

| File | Path | Notes |
|------|------|-------|
| README.md | ./README.md | Project overview, OSF boundaries, endpoints, corridor |
| REPO_WORKFLOW.md | ./.github/REPO_WORKFLOW.md | Repo governance: issue+PR for large changes, conventional commits |
| PULL_REQUEST_TEMPLATE.md | ./.github/PULL_REQUEST_TEMPLATE.md | PR template: linked issue, type, summary, changes table, validation |
| specs/_status.md | ./specs/_status.md | Current state: all phases closed, baseline reusable v1 |
| specs/master-plan.md | ./specs/master-plan.md | 5-phase master plan (A-E) |
| docs/ALPHA_SCOPE.md | ./docs/ALPHA_SCOPE.md | Alpha technical scope and boundaries |
| docs/RUNBOOK.md | ./docs/RUNBOOK.md | How to run the intent-service |

Read the convention files listed above for project-specific patterns and rules. All referenced paths have been extracted — no need to read index files to discover more.
