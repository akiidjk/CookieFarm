import { Cookie } from 'lucide-react';
import Link from 'next/link';
import { SiGithub, SiX } from '@icons-pack/react-simple-icons';

const navigation = {
  docs: [
    { name: 'Introduction', href: '/docs' },
    { name: 'Installation', href: '/docs/installation' },
    { name: 'Writing Exploits', href: '/docs/exploits' },
    { name: 'Configuration', href: '/docs/configuration' },
  ],
  community: [
    { name: 'GitHub', href: 'https://github.com/ByteTheCookies/CookieFarm' },
    { name: 'Discord', href: '#' },
    { name: 'Twitter', href: '#' },
  ],
  team: [
    { name: 'ByteTheCookies', href: 'https://github.com/ByteTheCookies' },
    { name: 'Contributors', href: 'https://github.com/ByteTheCookies/CookieFarm/graphs/contributors' },
  ],
};

export function Footer() {
  return (
    <footer className="border-t border-border bg-card/50">
      <div className=" px-4 py-12">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 px-5">
          {/* Brand */}
          <div className="col-span-2 md:col-span-1">
            <Link href="/" className="flex items-center gap-2 mb-4">
              <Cookie className="size-6 text-amber-500" />
              <span className="font-bold text-lg">CookieFarm</span>
            </Link>
            <p className="text-sm text-muted-foreground mb-4">
              Attack/Defense CTF Framework by ByteTheCookies
            </p>
            <div className="flex items-center gap-4">
              <a
                href="https://github.com/ByteTheCookies/CookieFarm"
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                <SiGithub className="size-5" />
              </a>
              <a
                href="#"
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                <SiX className="size-5" />
              </a>
            </div>
          </div>

          {/* Docs */}
          <div>
            <h3 className="font-semibold text-foreground mb-4">Documentation</h3>
            <ul className="space-y-3">
              {navigation.docs.map((item) => (
                <li key={item.name}>
                  <Link
                    href={item.href}
                    className="text-sm text-muted-foreground hover:text-amber-500 transition-colors"
                  >
                    {item.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Community */}
          <div>
            <h3 className="font-semibold text-foreground mb-4">Community</h3>
            <ul className="space-y-3">
              {navigation.community.map((item) => (
                <li key={item.name}>
                  <a
                    href={item.href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-muted-foreground hover:text-amber-500 transition-colors"
                  >
                    {item.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          {/* Team */}
          <div>
            <h3 className="font-semibold text-foreground mb-4">Team</h3>
            <ul className="space-y-3">
              {navigation.team.map((item) => (
                <li key={item.name}>
                  <a
                    href={item.href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-muted-foreground hover:text-amber-500 transition-colors"
                  >
                    {item.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Bottom */}
        <div className="mt-12 pt-8 border-t border-border px-5">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <p className="text-sm text-muted-foreground">
              &copy; {new Date().getFullYear()} ByteTheCookies. Released under MIT License.
            </p>
            <p className="text-sm text-muted-foreground">
              Made with <span className="text-red-500">❤️</span> in Italy
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
}
