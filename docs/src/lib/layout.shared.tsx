import type { BaseLayoutProps, LinkItemType } from 'fumadocs-ui/layouts/shared';
import { appName, gitConfig } from './shared';
import { Book, BookIcon, Cookie, DownloadIcon, Syringe, SyringeIcon } from 'lucide-react';
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
