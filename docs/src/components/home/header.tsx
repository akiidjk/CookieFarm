'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Cookie, Menu, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { SiGithub } from '@icons-pack/react-simple-icons';

const navigation = [
  { name: 'Documentation', href: '/docs' },
  { name: 'Installation ', href: '/docs/installation' },
  { name: 'Exploits', href: '/docs/exploits' },
];

export function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="fixed top-0 left-0 right-0 z-50 border-b border-border/50 bg-background/80 backdrop-blur-xl">
      <nav className="px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2 font-bold text-lg">
            <Cookie className="size-6 text-amber-500" />
            <span className="text-foreground">Cookie
              <span className="text-amber-500">Farm</span>
            </span>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center gap-8">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
              >
                {item.name}
              </Link>
            ))}
          </div>

          {/* Desktop Actions */}
          <div className="hidden md:flex items-center gap-4">
            <Button asChild variant="ghost" size="sm">
              <a
                href="https://github.com/ByteTheCookies/CookieFarm"
                target="_blank"
                rel="noopener noreferrer"
              >
                <SiGithub className="size-4 mr-2" />
                GitHub
              </a>
            </Button>
            <Button asChild size="sm" className="bg-amber-500 hover:bg-amber-600 text-black font-semibold">
              <Link href="/docs">Get Started</Link>
            </Button>
          </div>

          {/* Mobile Menu Button */}
          <button
            type="button"
            className="md:hidden p-2 text-muted-foreground hover:text-foreground"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          >
            {mobileMenuOpen ? <X className="size-6" /> : <Menu className="size-6" />}
          </button>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <div className="md:hidden py-4 space-y-4 border-t border-border">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className="block text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
                onClick={() => setMobileMenuOpen(false)}
              >
                {item.name}
              </Link>
            ))}
            <div className="flex flex-col gap-2 pt-4 border-t border-border">
              <Button asChild variant="outline" size="sm" className="justify-center">
                <a
                  href="https://github.com/ByteTheCookies/CookieFarm"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <SiGithub className="size-4 mr-2" />
                  GitHub
                </a>
              </Button>
              <Button asChild size="sm" className="bg-amber-500 hover:bg-amber-600 text-black font-semibold justify-center">
                <Link href="/docs" onClick={() => setMobileMenuOpen(false)}>
                  Get Started
                </Link>
              </Button>
            </div>
          </div>
        )}
      </nav>
    </header>
  );
}
