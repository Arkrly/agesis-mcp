# Aegis-MCP Frontend

A production-ready React frontend for the Aegis-MCP security gateway. Built with the "HackCulture" design system.

## Tech Stack

- **Framework:** React 19 (Vite)
- **Styling:** Tailwind CSS v4
- **Animations:** Framer Motion
- **Icons:** Lucide React
- **Forms:** React Hook Form + Zod
- **API Client:** Axios
- **Routing:** React Router v7

## Getting Started

### Prerequisites

- Node.js (v18+)
- npm

### Installation

```bash
cd frontend
npm install
```

### Development

1. Create a `.env` file (already provided with defaults):
   ```
   VITE_API_URL=http://localhost:8080
   ```
2. Run the development server:
   ```bash
   npm run dev
   ```

### Production Build

```bash
npm run build
```

The production assets will be generated in the `dist/` directory.

## Project Structure

- `src/api/`: API client and interceptors.
- `src/components/`: Reusable UI components (Layout, etc.).
- `src/pages/`: Page-level components (Dashboard, Policies, etc.).
- `src/index.css`: Tailwind v4 configuration and global styles.

## Design System (HackCulture)

The frontend implements the "HackCulture" theme from `DESIGN.md`:
- **Typography:** Space Grotesk
- **Colors:** Electric blue-violet primary accents with clean editorial spacing.
- **Components:** High-contrast buttons, card-based layouts, and subtle animations.
