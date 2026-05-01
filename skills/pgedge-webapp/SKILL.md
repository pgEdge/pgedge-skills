---
name: pgedge-webapp
description: "Use when the user asks to scaffold, set up, create, or bootstrap a new full-stack pgEdge web application or web client. Trigger on phrases like 'new web app', 'scaffold a frontend', 'set up a React app', 'new pgEdge web project', 'create a new client', or when working in an empty/new directory and the user mentions React, MUI, or login pages in the context of starting from scratch. Do NOT trigger for edits to an existing app, for adding a single component, or for purely backend Go work."
---

# pgEdge Web App Scaffolding Skill

This skill scaffolds a full-stack pgEdge web application from a single
template tree. The scaffolded project ships with locked-in design,
configuration, accessibility, and test-coverage conventions. Every future
pgEdge web app inherits the same UI patterns, the same form-label CSS,
the same auth flow, and the same 90% coverage gate.

## When This Skill Activates

Activates when the user asks to scaffold, bootstrap, or set up a new
full-stack pgEdge web application. Does not activate when editing an
existing app, when adding a single component, or for backend-only work.

## Project Setup from Template

The skill copies the entire `template/` tree wholesale into the target
directory, then walks the result to substitute placeholders.

### Inputs to Collect

Ask one question at a time. Validate as you go. Defaults are shown in
parentheses.

1. Target directory (must be empty or non-existent — abort if it has any
   contents, do not overwrite).
2. Project display name (e.g., "My App").
3. Project slug — kebab-case, must match `^[a-z][a-z0-9-]*$`.
4. Binary name — kebab-case (e.g., `my-app-server`).
5. Go module path — must contain at least one `/`.
6. HTTP listen port — integer in 1024–65535 (default 8080).
7. Initial admin password — at least 12 characters, prompted not echoed.
8. Include Ellie chat panel? (default no.)

### Setting Up a New Project

Execute the following steps in order. Surface any failure verbatim and
stop. Do not report success on a broken scaffold.

1. **Verify preconditions.** Target directory does not exist or is
   empty. Workbench reference path is reachable when needed.
2. **Copy the template tree** with a single `cp -r` operation.
3. **Replace placeholders.** Walk the tree and substitute the values
   from `PLACEHOLDERS.md`.
4. **Generate the seed admin hash.** Run
   `tools/genhash <password>` to produce the bcrypt hash, then
   substitute `<SEEDED_ADMIN_PASSWORD_HASH>` in
   `server/src/cmd/<binary>/seed.go`.
5. **Conditionally wire Ellie.** If the user said yes, edit
   `client/src/App.tsx` to import `EllieChatPanel` and `ChatFAB`,
   add `useState` for `chatOpen`, and render both inside the main
   layout. Concrete diff is in the Wiring Ellie section below.
6. **Verify the scaffold builds.** Run `go build ./...`,
   `go vet ./...`, `npm install`, `npm run typecheck`,
   `npm run lint`. Surface any error.
7. **Run the test suites.** `go test -race ./...` and
   `npm run test:coverage`. Confirm both report 90%+ coverage.
8. **Print the completion summary.** Files created (count by
   category), commands to run, URLs (`http://localhost:5173` for
   vite, `http://localhost:<port>` for the server), the seeded
   admin username (do not print the password), a reminder to
   rotate the seed password after first login.
9. **Offer pgedge-docs hand-off.** Ask: "Scaffold user-facing
   documentation now?" If yes, invoke `pgedge-docs`.

## Required Conventions

These rules are non-negotiable and apply during scaffolding and during
later edits to projects scaffolded from this template.

### Form Field Pattern

Every `<TextField>`, `<Select>`, and `<Autocomplete>` must use:

```tsx
<TextField
    label="..."
    InputLabelProps={{ shrink: true }}
    sx={{ ...textFieldSx, ...SELECT_FIELD_SX }}
/>
```

Use `SELECT_FIELD_SX` for fields on `background.paper` surfaces (dialogs,
cards). Use `SELECT_FIELD_DEFAULT_BG_SX` for fields on
`background.default` (full-page panels). For other surfaces, call
`getSelectFieldSx('your.bg.color')` rather than inventing styles. The
`textFieldSx` per-component object provides hover and focus border
colors; copy the pattern from `WelcomeCard.tsx`.

### Theme and Styling

- All colors come from `pgedgeTheme.ts`. No raw hex inline.
- All spacing uses MUI `theme.spacing()`. No raw `px` for spacing.
- All radii come from `theme.shape.borderRadius`.
- All typography uses theme variants.
- Light/dark mode is selected by the user via the Header toggle and
  persisted in `localStorage['theme-mode']`. Components check
  `theme.palette.mode === 'dark'` via `useTheme()`, not
  `prefers-color-scheme`.

### Style Constants

Every component file declares its `sx` objects as named constants at
the top of the file before the component body, with banner comments
delimiting the section. Theme-aware styles are factory functions:
`const getAppBarSx = (theme: Theme) => ({ ... })`.

### File and Directory Conventions

- One component per file; PascalCase filenames.
- Components with subcomponents live in a directory:
  `ComponentName/index.tsx` plus peer files plus `__tests__/`.
- Tests co-located in `__tests__/` directories.
- React contexts: provider in `ContextName.tsx`, consumer hook in
  `useContextName.ts` (separate files for fast-refresh).
- Type-only imports use `import type`.

### API Client Conventions

All HTTP requests go through `apiClient.ts`. All requests use
`credentials: 'include'`. Errors are `Error` instances and surface via
`<Alert severity="error">`. Response types are explicit generics — no
`any`.

### Backend Conventions

- Constructors take a logger and dependencies; no package-level state.
- Errors wrap with `fmt.Errorf("doing X: %w", err)`.
- HTTP handlers use `http.Handler` interfaces, not raw functions.
- No `os.Getenv` calls in `cmd/`, `internal/config/`, or
  `internal/auth/`. Configuration enters via CLI flags and YAML only.
- Logging through an injected `*slog.Logger`. No `log.Print*`, no
  global logger.

### Naming

- Go: lowercase package names, no underscores; PascalCase types and
  exports; camelCase unexported.
- TypeScript: PascalCase types and components; camelCase functions and
  variables; SCREAMING_SNAKE_CASE only for truly constant exported
  primitives like `SELECT_FIELD_SX`.

### Comments

No comments restating what code does. Comments earn their place by
explaining a non-obvious why: a workaround, an invariant, a load-bearing
constraint.

### Devcontainer

Scaffolded projects ship with `.devcontainer/devcontainer.json` plus
`post-create.sh`, mirroring the `pgedge-skills` repo's own pattern. The
container provisions Go 1.23 and Node 24, installs Helm via
`post-create.sh`, runs `go mod download` and `npm ci` so the dev loop
starts hot, and forwards the templated HTTP port. The
`overrideFeatureInstallOrder` block is load-bearing — it ensures Node
lands on `PATH` before later steps. Do not remove or reorder it.

## Accessibility Requirements

WCAG 2.1 AA at minimum. Tests assert zero AA violations via vitest-axe.
Hard rules:

- Every `<IconButton>` carries an explicit `aria-label`.
- Form labels associate with inputs via MUI `<TextField label>` or
  explicit `htmlFor`/`id`.
- Tooltips supplement, never replace, labels.
- Dialogs trap focus and have `aria-labelledby`.
- Color is never the only conveyor of meaning.
- Loading states announce themselves
  (`<CircularProgress aria-label="Loading application" />`).
- No `onClick` on a generic `<div>`/`<Box>` without role/tabIndex/
  onKeyDown — use `<Button>` or `<IconButton>` instead.
- Animations are wrapped with
  `@media (prefers-reduced-motion: no-preference)` so they resolve to
  static positions when the user prefers reduced motion.

## Test and Coverage Requirements

- Frontend: `vitest.config.js` enforces 90% on lines/branches/functions/
  statements. Excluded: `src/test/`, `src/main.tsx`, config files.
- Backend: `make coverage-check` runs the local
  `tools/coverage-check` and fails below 90%. Excluded:
  `cmd/<binary>/main.go`. Adding new exclusions requires reviewer
  approval — flag any such request.
- Use `renderWithTheme`, never raw `render` from RTL.
- Mock `fetch` globally via `vi.fn()` in `beforeEach`.
- Use `screen.getByRole` / `getByLabelText` over `getByTestId`.

## CI Requirements

The template ships with `.github/workflows/ci-client.yml` (Node 22 + 24
matrix) and `.github/workflows/ci-server.yml` (Go 1.23 + 1.24 matrix).
Both run lint, vet/typecheck, build, and the coverage gate. Coverage
artifacts are uploaded from the latest version of each runtime.

## Wiring Ellie (Optional)

If the user opted to include Ellie, edit `client/src/App.tsx` and:

1. Add the imports below the existing import block:

```tsx
import EllieChatPanel from './components/EllieChatPanel';
import ChatFAB from './components/EllieChatPanel/ChatFAB';
```

2. Inside `AppContent`, after the `useEffect` that persists `mode`,
   add:

```tsx
const [chatOpen, setChatOpen] = useState(false);
```

3. Inside the authenticated branch's `<Box sx={styles.mainLayoutBody}>`,
   add `<EllieChatPanel open={chatOpen} onClose={() => setChatOpen(false)} />`
   alongside `<WelcomeCard />`. After the closing `</Box>` for
   `mainLayoutBody`, before the closing `</Box>` for
   `mainLayoutRoot`, add:

```tsx
<ChatFAB onClick={() => setChatOpen(true)} isOpen={chatOpen} />
```

4. Set `ellie_wired: true` in `.pgedge-webapp.json`.

## Final Steps

### Hand-off to pgedge-docs

After scaffolding, ask: "Scaffold user-facing documentation now?" On
yes, invoke `pgedge-docs`. The docs skill recognizes the
`.pgedge-webapp.json` marker and applies its In-App Help Panel rules
when authoring content under `client/src/components/HelpPanel/pages/`.

### Completion Summary

Print files-created counts, commands to run, URLs, admin username (not
the password), the rotate-the-password reminder, and a one-line note
about pgedge-docs availability if the user declined the hand-off.

## Editing an Existing Scaffolded Project

When the skill is reinvoked on a project containing
`.pgedge-webapp.json`, run a consistency check before any edit:

- Are all `<TextField>` instances using `SELECT_FIELD_SX` and
  `shrink: true`?
- Any hex values in JSX that should be theme references?
- Any `<IconButton>` missing `aria-label`?
- Any `os.Getenv` calls in the server?
- Is current coverage above 90% on both sides?

Report findings before applying the user's actual request. Apply the
conventions automatically during the edit.

## Anti-patterns

These are forbidden. Refuse to add or accept them:

- `os.Getenv` in `cmd/`, `internal/config/`, or `internal/auth/`.
- Binary-relative config lookup
  (`<binary-dir>/<binary>.yaml`). The binary's parent directory is not
  a configuration source.
- Bypassing `apiClient.ts` with raw `fetch`.
- Hex colors outside `pgedgeTheme.ts`.
- `<TextField>` without `SELECT_FIELD_SX` and `shrink: true`.
- `<IconButton>` without `aria-label`.
- Adding coverage exclusions to dodge the 90% gate.
- `console.*` outside `utils/logger.ts`.
- Adding packages that overlap with existing scaffold deps (a second
  router, HTTP client, testing library).
