import type { BaseLayoutProps, LinkItemType } from 'fumadocs-ui/layouts/shared';
import { appName, gitConfig } from './shared';
import { Cookie } from 'lucide-react';

export function baseOptions(): BaseLayoutProps {
  return {
    nav: {
      // JSX supported
      title: (
        <>
          <Cookie className="size-6 text-amber-500" />
          <span className="text-foreground">Cookie
            <span className="text-amber-500">Farm</span>
          </span>
        </>
      ),
      url: '/',
    },
    links: [
      {
        text: 'Documentation',
        url: '/docs',
      },
      {
        text: 'Installation',
        url: `/docs/installation`,
      },
      {
        text: 'Exploits',
        url: `/docs/exploits`,
      },
    ] as LinkItemType[],
    githubUrl: `https://github.com/${gitConfig.user}/${gitConfig.repo}`,
  };
}
