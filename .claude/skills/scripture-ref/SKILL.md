---
name: scripture-ref
description: Add interactive scripture verse tooltips to template pages. Use when adding ESV verse tooltips to parenthetical Bible citations.
argument-hint: "[page-name]"
disable-model-invocation: true
---

# Scripture Reference Tooltip Skill

Add interactive tooltips that display ESV verse text on hover/tap to parenthetical scripture citations in Templ template pages.

## Arguments

- `/scripture-ref <page-name>` — Add scripture tooltips to citations on the specified page template

## How It Works

The project has a reusable tooltip system built with three components in `templates/components/scripture_ref.templ`:

- **`ScriptureScope()`** — Alpine.js wrapper that manages tooltip state (show/hide/position). Wrap the content area containing citations in this component. It provides `{ children... }` slot and automatically includes the shared tooltip element.
- **`ScriptureRef(id, ref, verseText, suffix)`** — Renders an interactive `<cite>` element. Hover (desktop), tap (mobile), and focus (keyboard) trigger the tooltip.
- **`ScriptureTooltip()`** — Shared floating tooltip. Rendered automatically inside `ScriptureScope`.

CSS is already defined in `static/css/components.css` (`.scripture-ref`, `.scripture-tooltip`) and `static/css/print.css`.

## Process

### 1. Identify Citations

Read the target page template in `templates/pages/`. Find all parenthetical scripture citations — text like `(Rom. 5:18-19)` or `(John 3:16; John 10:27-30)`.

### 2. Look Up ESV Verse Text

For each citation, look up the exact ESV text. Use `WebFetch` with `https://www.esv.org/<Book>+<Chapter>:<Verses>/` to get accurate text from esv.org.

For multi-verse groups (e.g., `1 Tim. 3:15; Matt. 28:19; 16:19`), fetch each reference separately.

For full-chapter references (e.g., `John 11`, `Acts 6`), select 1-4 key representative verses and add `<em>See [chapter] for full context.</em>` at the end.

### 3. Format Verse Text as HTML

The `verseText` parameter accepts HTML rendered via Alpine.js `x-html`. Format multi-reference tooltips using `<strong>` tags as section headers:

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

Replace each `(Reference)` with a component call. Each call must be on its own line. The `suffix` parameter handles trailing punctuation (`.`, `,`, `;`) to avoid whitespace issues.

**Before:**
```
...the believer (Rom. 5:18-19). The sole ground...
```

**After:**
```
...the believer
@components.ScriptureRef("unique-id", "Rom. 5:18-19", `<strong>Romans 5:18&ndash;19</strong>verse text...`, ".")
The sole ground...
```

### 6. Parameter Reference

| Parameter | Description | Example |
|-----------|-------------|---------|
| `id` | Unique identifier for this citation | `"sola-fide-1"`, `"total-depravity"` |
| `ref` | Display text shown in parentheses | `"Rom. 5:18-19"`, `"John 3:16; John 10:27-30"` |
| `verseText` | HTML string with ESV verse text (use backtick raw string) | `` `<strong>Romans 5:18</strong>Therefore...` `` |
| `suffix` | Punctuation after the closing paren | `"."`, `","`, `""` |

### 7. ID Naming Convention

Use descriptive kebab-case IDs that relate to the theological topic:
- `sola-fide-1`, `sola-fide-2` — numbered when multiple citations support one point
- `total-depravity`, `unconditional-election` — doctrine names
- `baptism-circumcision`, `lords-supper-institution` — sacrament-specific
- `elders-plurality`, `deacons-qualifications` — office-specific

## Example: Complete Citation Replacement

```go
// Single reference
@components.ScriptureRef("covenant-promise", "Gen. 3:15", `<strong>Genesis 3:15</strong>I will put enmity between you and the woman, and between your offspring and her offspring; he shall bruise your head, and you shall bruise his heel.`, ".")

// Multi-reference group
@components.ScriptureRef("marks-church", "1 Tim. 3:15; Matt. 28:19; 16:19", `<strong>1 Timothy 3:15</strong>...if I delay, you may know how one ought to behave in the household of God, which is the church of the living God, a pillar and buttress of the truth. <strong>Matthew 28:19</strong>Go therefore and make disciples of all nations, baptizing them in the name of the Father and of the Son and of the Holy Spirit, <strong>Matthew 16:19</strong>I will give you the keys of the kingdom of heaven, and whatever you bind on earth shall be bound in heaven, and whatever you loose on earth shall be loosed in heaven.`, ".")

// Mid-sentence citation (comma suffix)
@components.ScriptureRef("baptism-circumcision", "Col. 2:11-12", `<strong>Colossians 2:11-12</strong>In him also you were circumcised with a circumcision made without hands...`, ",")
```

## Existing Implementation

See `templates/pages/about_beliefs.templ` for a complete example with 20 scripture reference tooltips.

## After Adding Tooltips

1. Run `make generate` to compile the updated `.templ` file
2. Run `go build ./...` to verify the build succeeds
3. Test in browser: hover citations (desktop), tap (mobile), Tab key (keyboard), Escape to dismiss
4. Verify print preview hides tooltips
