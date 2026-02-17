---
name: scripture-ref
description: Add interactive scripture verse tooltips to template pages. Use when adding ESV verse tooltips to Bible citations (parenthetical or inline).
argument-hint: "[page-name]"
disable-model-invocation: true
---

# Scripture Reference Tooltip Skill

Add interactive tooltips that display ESV verse text on hover/tap to scripture citations in Templ template pages. Citations can be parenthetical (e.g., `(Rom. 5:18–19)`) or inline (e.g., `Ezekiel 1:5–10`).

## Arguments

- `/scripture-ref <page-name>` — Add scripture tooltips to citations on the specified page template

## How It Works

The project has a reusable tooltip system built with three components in `templates/components/scripture_ref.templ`:

- **`ScriptureScope()`** — Alpine.js wrapper that manages tooltip state (show/hide/position). Wrap the content area containing citations in this component. It provides `{ children... }` slot and automatically includes the shared tooltip element.
- **`ScriptureRef(id, ref, verseText, suffix)`** — Renders an interactive `<cite>` element. Hover (desktop), tap (mobile), and focus (keyboard) trigger the tooltip. The `ref` text is rendered as-is inside the `<cite>` tag.
- **`ScriptureTooltip()`** — Shared floating tooltip. Rendered automatically inside `ScriptureScope`.

CSS is already defined in `static/css/components.css` (`.scripture-ref`, `.scripture-tooltip`) and `static/css/print.css`.

## Process

### 1. Identify Citations

Read the target page template in `templates/pages/`. Find all scripture citations — both parenthetical like `(Rom. 5:18-19)` and inline like `Ezekiel 1:5–10`.

### 2. Look Up ESV Verse Text

For each citation, look up the exact ESV text. Use `WebFetch` with `https://www.esv.org/<Book>+<Chapter>:<Verses>/` to get accurate text from esv.org.

Each verse reference always gets its own individual `ScriptureRef` call and tooltip — even when multiple references appear within the same parenthetical group (e.g., `(1 Tim. 3:15; Matt. 28:19; 16:19)`). Fetch each reference separately.

For full-chapter references (e.g., `John 11`, `Acts 6`), select 1-4 key representative verses and add `<em>See [chapter] for full context.</em>` at the end.

### 3. Format Verse Text as HTML

The `verseText` parameter accepts HTML rendered via Alpine.js `x-html`. Each tooltip contains a single reference with a `<strong>` header followed by the verse text:

```
<strong>Romans 5:18&ndash;19</strong>Therefore, as one trespass led to condemnation for all men, so one act of righteousness leads to justification and life for all men. For as by the one man&rsquo;s disobedience the many were made sinners, so by the one man&rsquo;s obedience the many will be made righteous.
```

**Escaping rules** (the text passes through Templ attribute escaping → JS getAttribute → x-html):
- Use `&ldquo;` and `&rdquo;` for double quotes within verse text
- Use `&rsquo;` for apostrophes/single quotes within verse text
- `<strong>`, `<em>` HTML tags work as-is (Templ escapes them in the attribute, browser decodes them, Alpine renders them)

### 4. Wrap Content in ScriptureScope

Add the `components` import if not already present. Wrap the content area in `@components.ScriptureScope()`:

```go
@components.ScriptureScope() {
    <div class="content-section">
        // ... page content with ScriptureRef calls ...
    </div>
    <p class="esv-copyright">
        Scripture quotations are from the ESV&reg; Bible (The Holy Bible, English Standard Version&reg;), copyright &copy; 2001 by Crossway, a publishing ministry of Good News Publishers. Used by permission. All rights reserved.
    </p>
}
```

### 5. Replace Citations with ScriptureRef Calls

Replace each citation with a component call. Each call must be on its own line. **Important:** Templ does NOT insert whitespace between text nodes and component calls. Use `&nbsp;` before citations and trailing spaces in suffixes to manage spacing.

**Parenthetical citation — single reference** — include `(` and `)` in the `ref` parameter:

Before:
```
...the believer (Rom. 5:18-19). The sole ground...
```

After:
```
...the believer&nbsp;
@components.ScriptureRef("unique-id", "(Rom. 5:18-19)", `<strong>Romans 5:18&ndash;19</strong>verse text...`, ". ")
The sole ground...
```

**Parenthetical citation — multi-reference group** — split into individual `ScriptureRef` calls. The opening `(` goes on the first ref, the closing `)` goes on the last ref. Semicolons stay with each ref's display text:

Before:
```
...rightful head (1 Tim. 3:15; Matt. 16:19; Matt. 28:19; 1 Cor. 11:24–26).
```

After:
```
...rightful head&nbsp;
@components.ScriptureRef("marks-church-tim", "(1 Tim. 3:15;", `<strong>1 Timothy 3:15</strong>verse text...`, " ")
@components.ScriptureRef("marks-church-matt-16", "Matt. 16:19;", `<strong>Matthew 16:19</strong>verse text...`, " ")
@components.ScriptureRef("marks-church-matt-28", "Matt. 28:19;", `<strong>Matthew 28:19</strong>verse text...`, " ")
@components.ScriptureRef("marks-church-cor", "1 Cor. 11:24–26)", `<strong>1 Corinthians 11:24&ndash;26</strong>verse text...`, ".")
```

**Inline citation** — use plain `ref` without parentheses:

Before:
```
...the "four creatures" in Ezekiel 1:5–10, Ezekiel 10:14, and Revelation 4:6–7 to represent...
```

After:
```
...the &ldquo;four creatures&rdquo; in&nbsp;
@components.ScriptureRef("four-creatures", "Ezekiel 1:5–10", `...`, ", ")
@components.ScriptureRef("four-creatures-faces", "Ezekiel 10:14", `...`, ", ")
and&nbsp;
@components.ScriptureRef("four-creatures-rev", "Revelation 4:6–7", `...`, " ")
to represent...
```

### 6. Whitespace Rules

| Position | Pattern | Example |
|----------|---------|---------|
| Before citation (text precedes) | Add `&nbsp;` at end of preceding text line | `...the believer&nbsp;` |
| After citation (text follows, mid-paragraph) | Add trailing space to suffix | `". "`, `", "` |
| After citation (end of paragraph) | No trailing space needed | `"."`, `","` |
| Between adjacent citations | Use suffix trailing space (no `&nbsp;` line needed) | First: `", "`, second starts on next line |
| Before citation (word precedes, like "and") | Add `&nbsp;` at end of that word line | `and&nbsp;` |

### 7. Parameter Reference

| Parameter | Description | Example |
|-----------|-------------|---------|
| `id` | Unique identifier for this citation | `"sola-fide-1"`, `"total-depravity"` |
| `ref` | Display text rendered inside `<cite>` tag — include `(` `)` for parenthetical, plain for inline. For multi-ref groups, `(` goes on first ref and `)` on last ref, with semicolons on each ref | `"(Rom. 5:18-19)"`, `"(1 Tim. 3:15;"`, `"Matt. 28:19)"`, `"Ezekiel 1:5–10"` |
| `verseText` | HTML string with ESV verse text (use backtick raw string) | `` `<strong>Romans 5:18</strong>Therefore...` `` |
| `suffix` | Punctuation rendered immediately after `</cite>` — include trailing space if text follows | `". "`, `","`, `" "` |

### 8. ID Naming Convention

Use descriptive kebab-case IDs that relate to the theological topic:
- `sola-fide-1`, `sola-fide-2` — numbered when multiple citations support one point
- `total-depravity`, `unconditional-election` — doctrine names
- `baptism-circumcision`, `lords-supper-institution` — sacrament-specific
- `elders-plurality`, `deacons-qualifications` — office-specific
- `four-creatures`, `four-creatures-rev` — descriptive of the reference content

## Example: Complete Citation Replacement

```go
// Parenthetical — single reference
@components.ScriptureRef("covenant-promise", "(Gen. 3:15)", `<strong>Genesis 3:15</strong>I will put enmity between you and the woman, and between your offspring and her offspring; he shall bruise your head, and you shall bruise his heel.`, ".")

// Parenthetical — multi-reference group (each verse gets its own tooltip)
@components.ScriptureRef("marks-church-tim", "(1 Tim. 3:15;", `<strong>1 Timothy 3:15</strong>...if I delay, you may know how one ought to behave in the household of God, which is the church of the living God, a pillar and buttress of the truth.`, " ")
@components.ScriptureRef("marks-church-matt-16", "Matt. 16:19;", `<strong>Matthew 16:19</strong>I will give you the keys of the kingdom of heaven, and whatever you bind on earth shall be bound in heaven, and whatever you loose on earth shall be loosed in heaven.`, " ")
@components.ScriptureRef("marks-church-matt-28", "Matt. 28:19)", `<strong>Matthew 28:19</strong>Go therefore and make disciples of all nations, baptizing them in the name of the Father and of the Son and of the Holy Spirit,`, ".")

// Parenthetical — mid-sentence (comma suffix with trailing space)
@components.ScriptureRef("baptism-circumcision", "(Col. 2:11-12)", `<strong>Colossians 2:11-12</strong>In him also you were circumcised with a circumcision made without hands...`, ", ")

// Inline — no parentheses in ref
@components.ScriptureRef("four-creatures", "Ezekiel 1:5–10", `<strong>Ezekiel 1:5&ndash;10</strong>And from the midst of it came the likeness...`, ", ")
```

## Existing Implementations

- `templates/pages/about_beliefs.templ` — 46 scripture reference tooltips (parenthetical, individually split)
- `templates/pages/about_sanctuary.templ` — 4 citations (3 inline, 1 parenthetical)

## After Adding Tooltips

1. Run `make generate` to compile the updated `.templ` file
2. Run `go build ./...` to verify the build succeeds
3. Test in browser: hover citations (desktop), tap (mobile), Tab key (keyboard), Escape to dismiss
4. Verify print preview hides tooltips
