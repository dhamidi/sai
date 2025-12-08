# UI Design Guidelines

This project uses a **neo-brutalist** design style. Follow these conventions when creating or modifying UI components.

## Core Principles

- **High contrast**: Black (#000) and white (#fff) as primary colors
- **Bold borders**: 2-3px solid black borders on containers and inputs
- **Monospace typography**: SF Mono, Menlo, Monaco, Courier New
- **Uppercase text**: Headers, labels, and status indicators use `text-transform: uppercase`
- **Letter spacing**: 1-2px on uppercase text for readability
- **No rounded corners**: All elements have sharp, rectangular edges
- **No shadows**: Flat design without drop shadows or gradients

## Typography

- Font family: `"SF Mono", "Menlo", "Monaco", "Courier New", monospace`
- Line height: 1.5
- Labels: 12px, uppercase, bold, 1px letter-spacing
- Body text: 13-14px

## Interactive Elements

### Buttons
- Black background, white text
- 2px solid black border
- Uppercase, bold, 1px letter-spacing
- Hover: invert colors (white background, black text)

### Links
- Black text with underline
- `text-underline-offset: 2px`
- Hover: invert colors (black background, white text)

### Inputs
- 2px solid black border
- White background
- Focus: light gray background (#f0f0f0), no outline

## Status Indicators

- Pending: italic
- In progress: bold
- Completed: strikethrough
- Failed: uppercase

## Layout

- Container max-width: 960px
- Padding: 20-40px
- Sidebar width: 280px with 3px border
- Use CSS grid for definition lists
- Use flexbox for layouts

## Feedback Elements

### Progress bars
- 2px black border container
- Solid black fill
- Centered percentage text with `mix-blend-mode: difference`

### Errors
- 3px solid black border
- Uppercase, bold
- Prefixed with "ERROR: " via `::before`
