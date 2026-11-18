# Brand

← Back to the [README](README.md)

How Alumni looks and sounds. Everything here is already in the pages under
`internal/view/templates/`; this file writes it down so new pages match.

## Contents

- [Name and wordmark](#name-and-wordmark)
- [Tagline and voice](#tagline-and-voice)
- [Colors](#colors)
- [Typography](#typography)
- [Photography](#photography)
- [Layout](#layout)
- [Accessibility](#accessibility)

## Name and wordmark

- The product is **Alumni**; in full, *İstanbul Üniversitesi Mezun Ağı*.
- The wordmark is text, not an image: **İ.Ü.** in maroon, semibold, followed by
  *Alumni* in ink, both in Spectral. It sits top left in every page's nav and links to `/main`.
- Always write *İ.Ü.* with the dotted capital İ and Ü, never *IU* or *I.U.*

## Tagline and voice

> Mezuniyet bir son değil, *yeni bir sicilin* başlangıcıdır.

- The interface speaks **Turkish**; code, comments and documentation are in **English**.
- Calm and institutional, not a startup: short sentences, no exclamation marks, no emoji in the UI.
- Honest about what exists. Unfinished features are labelled as planned
  (*geliştirme aşamasında*), never shown as working.
- The one promise the brand makes: everyone in the feed is a verified İ.Ü. graduate.

## Colors

| Token | Hex | Use |
|---|---|---|
| `--paper` | `#e9e4d8` | page background |
| `--paper-deep` | `#ddd6c4` | a slightly darker paper; defined, not used yet |
| `--ink` | `#211c18` | body text, wordmark, hero background |
| `--muted` | `#5b5248` | secondary text, table headers, dates |
| `--maroon` | `#4a0e1a` | links, buttons, *İ.Ü.*, the closed-network panel |
| `--maroon-deep` | `#320a12` | a darker maroon; defined, not used yet |
| `--gold` | `#a9832d` | small accents on paper: roman numerals, signatures |
| light gold | `#e2c27c` | accents on dark backgrounds: hero kicker, emphasis, the year |
| `--line` | `rgba(33,28,24,0.16)` | hairline borders and dividers |

Maroon is the university's color and does the work; gold is only an accent.
Status tags are the one exception: accepted `#dfe8d6` / `#2f4a1f`, rejected `#f3dcd8` / maroon.

## Typography

Both from Google Fonts.

- **Spectral** (serif; 400, 500, 600, italic): headings, the wordmark, large numbers
  and italic accents. Headings use weight 500, never bold.
- **Inter** (sans-serif; 400, 500, 600): body text, forms, tables, buttons.
- Monospace (`ui-monospace`, Consolas) only for routes and code.

## Photography

- Real İstanbul University campus photos, taken by the project's author, in
  `internal/view/static/img/`: resized and with metadata stripped before they are committed.
- They appear behind the main page's hero under a dark gradient so the text stays readable.
- No stock photos. People's faces appear only on the main page, and only for
  graduates who opt in (see the [Roadmap](README.md#roadmap)).

## Layout

- Registry, not dashboard: content is listed in ruled rows (hairline `--line` dividers),
  not in cards with shadows.
- One bold moment per page at most, like the maroon closed-network panel on the main page.
- Small corner radius (3px) on buttons and inputs.
- Narrow reading width: 56rem on the main page, 38rem on `/about`, 44rem on the others.

## Accessibility

- Every link, button and input has a visible focus outline: 2px maroon, 3px offset.
- With *reduced motion* switched on, the hero slideshow stops on the first photo.
- Pages are `lang="tr"` and must work at phone width (640px and below).
