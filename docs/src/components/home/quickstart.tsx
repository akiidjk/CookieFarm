'use client';

import { useState } from 'react';
import { Check, Copy, Terminal } from 'lucide-react';
import { Button } from '@/components/ui/button';

const installCommands = {
  docker: 'docker pull bythethecookies/cookiefarm:latest',
  git: 'git clone https://github.com/ByteTheCookies/CookieFarm.git',
};

export function Quickstart() {
  const [copied, setCopied] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'docker' | 'git'>('docker');

  const copyToClipboard = (text: string, key: string) => {
    navigator.clipboard.writeText(text);
    setCopied(key);
    setTimeout(() => setCopied(null), 2000);
  };

  return (
    <section className="py-24 relative">
      <div className="absolute inset-0 grid-pattern opacity-30" />

      <div className=" relative z-10 px-4">
        <div className="text-center max-w-3xl mx-auto mb-12">
          <span className="inline-block px-4 py-1.5 rounded-full bg-amber-500/10 text-amber-400 text-sm font-medium mb-4">
            Quick Start
          </span>
          <h2 className="text-3xl md:text-4xl font-bold mb-4 text-balance">
            Get started in <span className="text-amber-500">seconds</span>
          </h2>
          <p className="text-lg text-muted-foreground">
            Choose your preferred installation method and start capturing flags immediately.
          </p>
        </div>

        <div className="max-w-2xl mx-auto">
          {/* Tab buttons */}
          <div className="flex gap-2 mb-4">
            <button
              onClick={() => setActiveTab('docker')}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${activeTab === 'docker'
                ? 'bg-amber-500 text-black'
                : 'bg-secondary text-muted-foreground hover:text-foreground'
                }`}
            >
              Docker
            </button>
            <button
              onClick={() => setActiveTab('git')}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${activeTab === 'git'
                ? 'bg-amber-500 text-black'
                : 'bg-secondary text-muted-foreground hover:text-foreground'
                }`}
            >
              Git Clone
            </button>
          </div>

          {/* Command box */}
          <div className="relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-amber-500/20 to-emerald-500/20 rounded-xl blur opacity-50 group-hover:opacity-100 transition-opacity" />
            <div className="relative bg-card border border-border rounded-xl overflow-hidden">
              <div className="flex items-center justify-between px-4 py-3 bg-secondary/50 border-b border-border">
                <div className="flex items-center gap-2">
                  <Terminal className="size-4 text-amber-500" />
                  <span className="text-sm text-muted-foreground font-mono">terminal</span>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => copyToClipboard(installCommands[activeTab], activeTab)}
                  className="h-8 px-3 text-muted-foreground hover:text-foreground"
                >
                  {copied === activeTab ? (
                    <Check className="size-4 text-emerald-500" />
                  ) : (
                    <Copy className="size-4" />
                  )}
                </Button>
              </div>

              <div className="p-6">
                <code className="font-mono text-sm md:text-base text-foreground">
                  <span className="text-emerald-400">$</span>{' '}
                  {installCommands[activeTab]}
                </code>
              </div>
            </div>
          </div>

          {/* Additional steps */}
          <div className="mt-8 grid gap-4">
            <div className="flex items-start gap-4 p-4 rounded-lg bg-card/50 border border-border">
              <div className="flex-shrink-0 size-8 rounded-full bg-amber-500/10 flex items-center justify-center text-amber-500 font-bold text-sm">
                1
              </div>
              <div>
                <h4 className="font-medium text-foreground mb-1">Configure your competition</h4>
                <p className="text-sm text-muted-foreground">
                  Set up flag format, teams, and submission endpoint in <code className="text-amber-400">config.yaml</code>
                </p>
              </div>
            </div>

            <div className="flex items-start gap-4 p-4 rounded-lg bg-card/50 border border-border">
              <div className="flex-shrink-0 size-8 rounded-full bg-emerald-500/10 flex items-center justify-center text-emerald-500 font-bold text-sm">
                2
              </div>
              <div>
                <h4 className="font-medium text-foreground mb-1">Write your exploits</h4>
                <p className="text-sm text-muted-foreground">
                  Create Python exploit files in the <code className="text-amber-400">exploits/</code> directory
                </p>
              </div>
            </div>

            <div className="flex items-start gap-4 p-4 rounded-lg bg-card/50 border border-border">
              <div className="flex-shrink-0 size-8 rounded-full bg-sky-500/10 flex items-center justify-center text-sky-500 font-bold text-sm">
                3
              </div>
              <div>
                <h4 className="font-medium text-foreground mb-1">Start the farm</h4>
                <p className="text-sm text-muted-foreground">
                  Run <code className="text-amber-400">./cookiefarm start</code> and watch the flags roll in
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
