import type { BaseLayoutProps } from 'fumadocs-ui/layouts/shared';
import { gitConfig } from './shared';
import { Nav } from '@/app/(home)/page';

export function baseOptions(): BaseLayoutProps {
  return {
    nav: { component: <Nav /> },
    themeSwitch: {
      enabled: false
    },

    githubUrl: `https://github.com/${gitConfig.user}/${gitConfig.repo}`,
  };
}
