'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { ArrowRight, Zap, Cookie } from 'lucide-react';
import { SiGithub } from '@icons-pack/react-simple-icons';
import { Button } from '@/components/ui/button';

const codeLines = [
  '#!/usr/bin/env python3',
  'from cookiefarm import exploit_manager',
  '',
  '@exploit_manager',
  'def exploit(ip, port, name_service, flag_ids: list):',
  '    # Run your exploit here',
  '    base_url = f"http://{ip}:{port}"',
  '',
  '    r = requests.get(f"{base_url}/get-flag")',
  '    print(r.text)',
];

function TypewriterCode() {
  const [displayedLines, setDisplayedLines] = useState<string[]>([]);
  const [currentLine, setCurrentLine] = useState(0);
  const [currentChar, setCurrentChar] = useState(0);

  useEffect(() => {
    if (currentLine >= codeLines.length) {
      // Reset after a delay
      const timeout = setTimeout(() => {
        setDisplayedLines([]);
        setCurrentLine(0);
        setCurrentChar(0);
      }, 3000);
      return () => clearTimeout(timeout);
    }

    const line = codeLines[currentLine];

    if (currentChar <= line.length) {
      const timeout = setTimeout(() => {
        setDisplayedLines(prev => {
          const newLines = [...prev];
          newLines[currentLine] = line.slice(0, currentChar);
          return newLines;
        });
        setCurrentChar(prev => prev + 1);
      }, 30);
      return () => clearTimeout(timeout);
    } else {
      setCurrentLine(prev => prev + 1);
      setCurrentChar(0);
    }
  }, [currentLine, currentChar]);

  return (
    <div className="font-mono text-sm leading-relaxed">
      {displayedLines.map((line, i) => (
        <div key={i} className="flex">
          <span className="text-muted-foreground w-8 select-none">{i + 1}</span>
          <span className={
            line.includes('class') ? 'text-amber-400' :
              line.includes('def') ? 'text-emerald-400' :
                line.includes('from') || line.includes('import') ? 'text-sky-400' :
                  line.includes('return') ? 'text-rose-400' :
                    'text-foreground'
          }>
            {line}
            {i === currentLine - 1 || (i === displayedLines.length - 1 && currentLine < codeLines.length) ? (
              <span className="inline-block w-2 h-4 bg-amber-500 animate-pulse ml-0.5" />
            ) : null}
          </span>
        </div>
      ))}
    </div>
  );
}

export function Hero() {
  return (
    <section className="relative min-h-[90vh] flex items-center justify-center overflow-hidden">
      {/* Background Effects */}
      <div className="absolute inset-0 grid-pattern opacity-50" />
      <div className="absolute top-1/4 -left-32 w-96 h-96 bg-amber-500/10 rounded-full blur-3xl" />
      <div className="absolute bottom-1/4 -right-32 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl" />

      <div className="container relative z-10 px-4 py-20">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left Content */}
          <div className="space-y-8">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-amber-500/10 border border-amber-500/20 text-amber-400 text-sm font-medium">
              <Cookie className="size-4" />
              <span>By ByteTheCookies</span>
            </div>

            <h1 className="text-5xl md:text-6xl lg:text-7xl font-bold tracking-tight text-balance">
              <span className="text-foreground">Cookie</span>
              <span className="text-amber-500">Farm</span>
            </h1>

            <p className="text-xl md:text-2xl text-muted-foreground max-w-lg leading-relaxed">
              Attack/Defense CTF framework with a{' '}
              <span className="text-emerald-400 font-semibold">zero distraction</span>{' '}
              approach. Your only task:{' '}
              <span className="text-foreground font-semibold">write the exploit logic!</span>
            </p>

            <div className="flex flex-wrap gap-4">
              <Button asChild size="lg" className="bg-amber-500 hover:bg-amber-600 text-black font-semibold animate-pulse-glow">
                <Link href="/docs">
                  Get Started
                  <ArrowRight className="ml-2 size-4" />
                </Link>
              </Button>
              <Button asChild variant="outline" size="lg" className="border-border hover:bg-secondary">
                <a href="https://github.com/ByteTheCookies/CookieFarm" target="_blank" rel="noopener noreferrer">
                  <SiGithub className="mr-2 size-4" />
                  View on GitHub
                </a>
              </Button>
            </div>

            <div className="flex items-center gap-6 pt-4 text-sm text-muted-foreground">
              <div className="flex items-center gap-2">
                <div className="size-2 rounded-full bg-emerald-500" />
                <span>Go + Python</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="size-2 rounded-full bg-amber-500" />
                <span>Open Source</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="size-2 rounded-full bg-sky-500" />
                <span>GPL-3.0 License</span>
              </div>
            </div>
          </div>

          {/* Right Content - Code Preview */}
          <div className="relative">
            <div className="absolute -inset-4 bg-linear-to-r from-amber-500/20 via-emerald-500/20 to-sky-500/20 rounded-2xl blur-xl opacity-50" />
            <div className="relative bg-card border border-border rounded-xl overflow-hidden shadow-2xl">
              {/* Terminal Header */}
              <div className="flex items-center gap-2 px-4 py-3 bg-secondary/50 border-b border-border">
                <div className="flex gap-2">
                  <div className="size-3 rounded-full bg-red-500/80" />
                  <div className="size-3 rounded-full bg-amber-500/80" />
                  <div className="size-3 rounded-full bg-emerald-500/80" />
                </div>
                <span className="ml-4 text-sm text-muted-foreground font-mono">exploit.py</span>
              </div>

              {/* Code Content */}
              <div className="p-6 min-h-[280px] bg-background/50">
                <TypewriterCode />
              </div>

              {/* Terminal Footer */}
              <div className="px-4 py-3 bg-secondary/30 border-t border-border flex items-center gap-2">
                <Zap className="size-4 text-amber-500" />
                <span className="text-xs text-muted-foreground font-mono">
                  Ready to capture flags
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
