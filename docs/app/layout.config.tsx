import type { BaseLayoutProps } from "fumadocs-ui/layouts/shared";

/**
 * Shared layout configurations
 *
 * you can customize layouts individually from:
 * Home Layout: app/(home)/layout.tsx
 * Docs Layout: app/docs/layout.tsx
 */
export const baseOptions: BaseLayoutProps = {
  nav: {
    title: (
      <>
        <img src="https://raw.githubusercontent.com/yesoreyeram/uql/refs/heads/main/public/logo.svg" width="24" height="24" aria-label="UQL Logo" />
        UQL
      </>
    ),
  },
  searchToggle: {
    enabled: false,
  },
};
