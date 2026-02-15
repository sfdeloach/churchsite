---
name: templ
description: Create or modify Templ templates following the project design system. Use when building page templates or reusable components.
argument-hint: "[page|component] [name]"
disable-model-invocation: true
---

# Templ Template Generator

Create or modify Templ templates (`.templ` files) following the project's design system and existing patterns.

## Arguments

- `/templ page <name>` — Create a page template in `templates/pages/`
- `/templ component <name>` — Create a reusable component in `templates/components/`

## Page Templates

Page templates live in `templates/pages/` and render full pages wrapped in the base layout.

### Pattern

```go
package pages

import "github.com/sfdeloach/churchsite/templates/layouts"
import "github.com/sfdeloach/churchsite/templates/components"
// import models if data-driven: import "github.com/sfdeloach/churchsite/internal/models"

templ PageName() {
    @layouts.Base("Page Title - Saint Andrew's Chapel") {
        @components.PageHeader("Page Title", "Optional subtitle text")
        <section class="content-section">
            <div class="container">
                // Page content here
            </div>
        </section>
    }
}
```

### Conventions

- Always wrap content in `@layouts.Base("Title")` — the base layout provides the HTML shell, nav, and footer
- Use `@components.PageHeader()` for the banner at the top of content pages
- Use semantic HTML: `<section>`, `<article>`, `<aside>`, `<nav>`
- Wrap content in `<div class="container">` for consistent max-width
- For data-driven pages, accept model types or slices as parameters
- Use `templ.SafeURL()` for dynamic URLs in `href` attributes

### Data-Driven Example

```go
templ AboutStaff(grouped map[models.StaffCategory][]models.StaffMember) {
    @layouts.Base("Pastors & Staff - Saint Andrew's Chapel") {
        @components.PageHeader("Pastors & Staff", "")
        <section class="about-content">
            <div class="container">
                for _, cat := range models.OrderedStaffCategories() {
                    if members, ok := grouped[cat]; ok {
                        <h2>{ models.StaffCategories[cat].Label }</h2>
                        <div class="staff-grid">
                            for _, member := range members {
                                @components.StaffCard(member)
                            }
                        </div>
                    }
                }
            </div>
        </section>
    }
}
```

## Component Templates

Components live in `templates/components/` and are reusable building blocks.

### Pattern

```go
package components

// import models if needed: import "github.com/sfdeloach/churchsite/internal/models"

templ ComponentName(param1 string, param2 int) {
    <div class="component-name">
        // Component markup
    </div>
}
```

### Existing Components to Reference

- `nav.templ` — Header + Alpine.js dropdown + responsive hamburger
- `footer.templ` — Footer with church info, dynamic copyright year
- `service_times.templ` — Sunday/Wednesday service time cards
- `event_card.templ` — Event card with date/time formatting
- `staff_card.templ` — Staff member card with photo placeholder
- `page_header.templ` — Reusable banner component

## Design System Reference

### CSS Classes Available

- `.container` — max-width wrapper with horizontal padding
- `.content-section` — standard section spacing
- `.page-header` — banner with title and optional subtitle
- `.staff-grid`, `.events-grid` — responsive grids
- `.btn`, `.btn--primary`, `.btn--secondary` — button styles
- `.card` — generic card component

### Brand Colors (CSS Custom Properties)

- `--color-primary: #89191C` (crimson)
- `--color-accent: #B8860B` (warm gold)
- Neutral grays are warm-tinted (not pure gray)

### Typography

- Headings: EB Garamond (serif) via `--font-heading`
- Body: Nunito (sans-serif) via `--font-body`

### HTMX Patterns

For dynamic content, use HTMX attributes:
- `hx-get="/path"` — fetch and swap content
- `hx-target="#element"` — where to put the response
- `hx-swap="innerHTML"` — how to insert (innerHTML, outerHTML, beforeend, etc.)
- `hx-trigger="click"` — what triggers the request

### Alpine.js Patterns

For client-side interactivity:
- `x-data="{ open: false }"` — component state
- `x-show="open"` — conditional display
- `x-on:click="open = !open"` — event handling
- `x-transition` — animation

## After Creating Templates

Remind the user to:
1. Run `make generate` to compile `.templ` files to Go (or rely on `air` in dev mode)
2. Add corresponding CSS classes to the appropriate file in `static/css/`
3. Create a handler method that renders the template
