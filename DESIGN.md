# DESIGN.md — Alumni

← Back to the [README](README.md)

How the Alumni pages should look and feel, written in the
[DESIGN.md format](https://github.com/VoltAgent/awesome-design-md) so that both people and
AI coding agents can build new pages from it.

> **Status: target design.** This is the direction the pages are moving in, not a description
> of today's pages. Where they differ (radius, dividers, buttons), the pages have not caught up yet.

**References.** Two styles from [Refero Styles](https://styles.refero.design) shaped this file:
[monopo saigon](https://styles.refero.design/style/3e52dd36-6ab1-48c6-bc40-47ef6d33abc2)
(monochrome editorial, alternating light and dark bands, pill buttons, no shadows) and
[ORYZO AI](https://styles.refero.design/style/1f204e95-454a-437e-845b-c1b169d35607)
(warm-dark canvas, cream instead of white, dashed hairlines, small uppercase labels).
Their ideas are taken; their fonts, colors and content are not. Alumni keeps its own
maroon, gold, Spectral and Inter.

## 1. Visual Theme & Atmosphere

An **editorial registry**: a university yearbook laid out like a gallery catalogue.

- Calm, spacious, institutional. Large serif headlines, a lot of empty paper around them.
- The page alternates **paper bands** (light) and **ink bands** (dark) instead of boxes and cards.
- Campus photography is the only image, and it lives in the hero.
- Nothing decorates: no icons for their own sake, no gradients on controls, no shadows.

## 2. Color Palette & Roles

| Name | Hex | Token | Role |
|---|---|---|---|
| Paper | `#e9e4d8` | `--paper` | canvas of light bands |
| Paper Deep | `#ddd6c4` | `--paper-deep` | a quieter light band, inset panels |
| Ink | `#211c18` | `--ink` | text on paper; canvas of dark bands |
| Muted | `#5b5248` | `--muted` | secondary text, dates, captions |
| Cream | `#f4ece2` | `--cream` | text on dark bands; **never pure `#fff`** |
| Maroon | `#4a0e1a` | `--maroon` | links, the one filled button, the closed-network band |
| Maroon Deep | `#320a12` | `--maroon-deep` | text on gold buttons, hover on maroon |
| Gold | `#a9832d` | `--gold` | small accents on paper: numerals, signatures |
| Light Gold | `#e2c27c` | `--gold-light` | accents on dark: kicker, emphasis, the filled button on dark |
| Line | `rgba(33,28,24,0.16)` | `--line` | hairlines on paper |
| Line on dark | `rgba(244,236,226,0.25)` | `--line-dark` | hairlines on ink |

Maroon carries the brand; gold is an accent and never fills a large area.

## 3. Typography Rules

- **Spectral** (serif): display, headings, the wordmark, large numbers. Weights 400 and 500;
  italic for emphasis. Never 600+ above 40px.
- **Inter** (sans-serif): body, forms, tables, buttons, labels. Weights 400, 500, 600.
- Letter-spacing 0 everywhere.

| Role | Font | Size | Weight | Line height |
|---|---|---|---|---|
| Display (hero) | Spectral | 64px | 400 | 1.05 |
| Heading | Spectral | 40px | 500 | 1.1 |
| Subheading | Spectral | 26px | 500 | 1.2 |
| Body large (lede) | Inter | 18px | 400 | 1.55 |
| Body | Inter | 16px | 400 | 1.6 |
| Small | Inter | 14px | 400 | 1.5 |
| Label | Inter | 12px | 500, UPPERCASE | 1.2 |

Labels (nav items, kickers, column headers, tags) are the only uppercase text.
Body copy is always sentence case and left-aligned.

## 4. Component Stylings

**Buttons** — the only rounded thing on the page (pill, `999px`).
- *Filled*: maroon fill, cream text on paper; light-gold fill, maroon-deep text on ink.
  Padding `12px 28px`, Inter 15px/500. **At most one per section.**
- *Ghost*: transparent, 1px border in ink (on paper) or cream at 50% (on ink). Same size.
- Hover: border or fill darkens; no lift, no shadow.

**Links** — maroon on paper, light gold on ink. Underline on hover only.

**Lists of records** (graduates, postings, announcements) — ruled rows, not cards:
title in Inter 16px/600, one muted meta line under it, a 1px `--line` hairline between rows.

**Section dividers** — 1px **dashed** `--line` between sections of the same band;
a change of band (paper → ink) needs no line.

**Inputs** — 0 radius, `--paper` slightly lighter fill (`#f6f2ea`), 1px `--line` border,
maroon 2px focus outline.

**Navigation** — wordmark (*İ.Ü.* maroon + *Alumni*) left; at most five links right,
as 12px uppercase labels. Transparent over the hero, hairline under it elsewhere.

**Centered growing object** — the one signature moment, used once per page, right after the hero.
An ink band about 2.5 screens tall whose content sticks while the visitor scrolls through it:
- In the middle: the Alumni medallion, *İ.Ü.* in light gold Spectral with *1453* under it and
  "İSTANBUL ÜNİVERSİTESİ · MEZUN AĞI · ALUMNI" around the edge. It is Alumni's own mark,
  not the university's official logo.
- As the band scrolls, the medallion stands up from tilted (`rotateX` 58° → 0°), grows
  (scale 0.42 → 1) and its ring of text turns (0° → 140°).
- A Spectral heading sits on its left and a short paragraph on its right, 18px gutters;
  both fade from 35% to full as the object grows, never from invisible.

**Tags** (application status) — pill, 12px label, soft fill:
accepted `#dfe8d6` / `#2f4a1f`, rejected `#f3dcd8` / maroon.

## 5. Layout Principles

- Spacing scale (px): **4, 8, 12, 16, 24, 32, 48, 64, 96, 128**. Sections are 96px apart
  on desktop, 64px on phones.
- Content column: **1080px** max for the main page, **720px** for reading and forms.
  Bands are full-bleed; the content inside them is not.
- Left-aligned text. Only the hero headline may sit centered.
- Two-column rows alternate sides (text left / image right, then the reverse).
- The growing-object band is the one three-column row: text · object · text.
- Generous whitespace is the design: when in doubt, add space, not a border.

## 6. Depth & Elevation

**None.** No `box-shadow`, no blur, no glass. Depth comes from:
1. switching band color (paper ↔ ink ↔ maroon), and
2. hairlines.

The hero photo sits under a dark gradient so text stays readable; that gradient is the only overlay.

## 7. Do's and Don'ts

**Do**
- Use cream (`#f4ece2`) for text on dark; keep pure white out of the UI.
- Keep corners square except buttons and tags.
- Give the hero headline room; one idea per band.
- Write labels in UPPERCASE, everything else in sentence case.

**Don't**
- Don't add drop shadows, gradients on buttons, or cards inside cards.
- Don't use a border-radius between 1px and 998px.
- Don't put more than one filled button in a section.
- Don't use the growing object more than once on a page, or for anything but the medallion.
- Don't center body copy or use gold for large areas.
- Don't use stock photos; only real campus photography from `internal/view/static/img/`.

## 8. Responsive Behavior

- Breakpoint: **640px**. Below it: display 40px, heading 28px; two-column rows stack;
  nav links wrap under the wordmark.
- Touch targets at least 44px tall.
- Wide tables scroll sideways inside their own box; the page itself never scrolls sideways.
- The growing-object band stacks on phones: medallion on top (240px), heading, then paragraph.
- With *prefers-reduced-motion*, the hero slideshow stops on the first photo, transitions are off,
  and the growing-object band becomes a normal-height band with the medallion already full size.
- Motion otherwise: `cubic-bezier(0.19, 1, 0.22, 1)`, 0.4s for color, 0.8s for movement.

## 9. Agent Prompt Guide

**Quick reference:** paper `#e9e4d8` · ink `#211c18` · cream `#f4ece2` · maroon `#4a0e1a` ·
gold `#a9832d` / `#e2c27c` · Spectral headings · Inter body · pill buttons only · no shadows.

**Prompts that work with this file:**
- "Build the `/graduates` page following DESIGN.md: a paper band with a 40px Spectral heading,
  then the graduates as ruled rows, one filled maroon pill button to add a graduate."
- "Add an ink band to the main page with three numbers in Spectral 64px and uppercase labels,
  following DESIGN.md."
- "Check this template against DESIGN.md and list every rule it breaks."
