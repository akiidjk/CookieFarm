import Link from "next/link";
import {
  Terminal,
  Zap,
  Code2,
  Flag,
  LayoutDashboard,
  Users,
  ArrowRight,
  Server,
} from "lucide-react";

import { SiGithub } from '@icons-pack/react-simple-icons';
// ─── Nav ────────────────────────────────────────────────────────────────────

export function Nav() {
  return (
    <header className="sticky top-0 z-50 border-b border-(--surface-border) bg-background/80 backdrop-blur-sm">
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-2">
          <span className="font-mono text-base font-bold text-(--green)">
            CookieFarm
          </span>
        </Link>

        {/* Nav links */}
        <nav className="flex items-center gap-6">
          <Link
            href="/docs"
            className="font-mono text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            Docs
          </Link>
          <Link
            href="https://github.com/ByteTheCookies/CookieFarm"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1.5 font-mono text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <SiGithub size={14} />
            GitHub
          </Link>
        </nav>
      </div>
    </header>
  );
}

// ─── Terminal Block ──────────────────────────────────────────────────────────

function TerminalBlock() {
  return (
    <div className="w-full overflow-hidden rounded-lg border border-(--surface-border) bg-(--surface)">
      {/* Window chrome */}
      <div className="flex items-center gap-1.5 border-b border-(--surface-border) px-4 py-2.5">
        <span className="h-3 w-3 rounded-full bg-[#ff5f57]" />
        <span className="h-3 w-3 rounded-full bg-[#febc2e]" />
        <span className="h-3 w-3 rounded-full bg-[#28c840]" />
        <span className="ml-3 font-mono text-xs text-muted-foreground">
          exploit.py
        </span>
      </div>
      {/* Code */}
      <pre className="overflow-x-auto p-5 text-sm leading-relaxed">
        <code className="font-mono">
          <span className="text-muted-foreground">#!/usr/bin/env python3</span>
          {"\n"}
          <span className="text-[oklch(0.7_0.12_260)]">import</span>
          <span className="text-foreground"> requests</span>
          {"\n"}
          <span className="text-[oklch(0.7_0.12_260)]">from</span>
          <span className="text-foreground"> cookiefarm </span>
          <span className="text-[oklch(0.7_0.12_260)]">import</span>
          <span className="text-foreground"> exploit_manager</span>
          {"\n\n"}
          <span className="text-foreground">@exploit_manager</span>
          {"\n"}
          <span className="text-[oklch(0.7_0.12_260)]">def</span>
          <span className="text-[oklch(0.85_0.14_90)]"> exploit</span>
          <span className="text-foreground">{"(ip: str, port: int, name: str):"}</span>
          {"\n"}
          <span className="text-foreground">{"    base_url = f\"http://{ip}:{port}\""}</span>
          {"\n\n"}
          <span className="text-foreground">{"    # Service 1"}</span>
          {"\n"}
          <span className="text-foreground">{"    r = requests.get(f\"{base_url}/get-flag\")"}</span>
          {"\n"}
          <span className="text-foreground">{"    print(r.text)"}</span>
          {"\n\n"}
          <span className="text-muted-foreground"># run: </span>
          <span className="text-(--green)">
            cookiefarm run exploit.py --config config.toml
          </span>
        </code>
      </pre>
    </div>
  );
}

// ─── Hero ────────────────────────────────────────────────────────────────────

function Hero() {
  return (
    <section className="mx-auto max-w-6xl px-6 pb-20 pt-24 md:pt-32">
      <div className="mb-4 inline-flex items-center gap-2 rounded-full border border-(--green)/30 bg-(--green)/5 px-3 py-1">
        <span className="h-1.5 w-1.5 rounded-full bg-(--green)" />
        <span className="font-mono text-xs text-(--green)">
          Attack / Defense CTF Framework
        </span>
      </div>

      <h1 className="mb-6 max-w-3xl text-balance text-4xl font-bold leading-tight tracking-tight text-foreground md:text-5xl lg:text-6xl">
        The exploit farm that stays{" "}
        <span className="text-(--green)">out of your way.</span>
      </h1>

      <p className="mb-10 max-w-2xl text-pretty text-lg leading-relaxed text-muted-foreground">
        CookieFarm is an A/D CTF framework with a Go server, Python client SDK,
        and zero-config flag submission. Write the exploit. We handle the rest.
      </p>

      {/* CTAs */}
      <div className="mb-14 flex flex-wrap items-center gap-3">
        <Link
          href="/docs"
          className="inline-flex items-center gap-2 rounded-md bg-(--green) px-5 py-2.5 font-mono text-sm font-semibold text-[oklch(0.1_0_0)] transition-opacity hover:opacity-90"
        >
          Get Started
          <ArrowRight size={14} />
        </Link>
        <Link
          href="https://github.com/ByteTheCookies/CookieFarm"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-md border border-(--surface-border) bg-(--surface) px-5 py-2.5 font-mono text-sm text-foreground transition-colors hover:border-(--green)/50 hover:text-(--green)"
        >
          <SiGithub size={14} />
          View on GitHub
        </Link>
      </div>

      {/* Terminal */}
      <TerminalBlock />
    </section>
  );
}

// ─── Feature Grid ─────────────────────────────────────────────────────────

const features = [
  {
    icon: Server,
    title: "Go Client and Server Core",
    description:
      "High-performance scheduler written in Go. Handles exploit parallelism, flag collection and timed execution cycles without breaking a sweat.",
  },
  {
    icon: Code2,
    title: "Python SDK",
    description:
      "A dead-simple client library. Import, subclass, write your attack logic. That's it.",
  },
  {
    icon: Zap,
    title: "Zero Distraction",
    description:
      "No YAML sprawl. No boilerplate. Your only job is the exploit function.",
  },
  {
    icon: Flag,
    title: "Auto Flag Submission",
    description:
      "Flags are detected, deduplicated, and submitted to the scoreboard automatically every tick.",
  },
  {
    icon: LayoutDashboard,
    title: "Live Dashboard",
    description:
      "Monitor exploit runs, flag counts and errors from a clean web UI in real time.",
  },
  {
    icon: Users,
    title: "Team-Ready",
    description:
      "Designed for competition environments. Deploys fast, scales with your team.",
  },
];

function FeatureGrid() {
  return (
    <section className="border-t border-(--surface-border)">
      <div className="mx-auto max-w-6xl px-6 py-20">
        <p className="mb-2 font-mono text-xs uppercase tracking-widest text-(--green)">
          Features
        </p>
        <h2 className="mb-12 text-balance text-2xl font-bold tracking-tight text-foreground md:text-3xl">
          Everything handled. Nothing in the way.
        </h2>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((f) => (
            <div
              key={f.title}
              className="group rounded-lg border border-(--surface-border) bg-(--surface) p-6 transition-colors hover:border-(--green)/40"
            >
              <div className="mb-4 inline-flex h-9 w-9 items-center justify-center rounded-md border border-(--green)/20 bg-(--green)/10 text-(--green)">
                <f.icon size={16} />
              </div>
              <h3 className="mb-2 font-mono text-sm font-semibold text-foreground">
                {f.title}
              </h3>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {f.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ─── Code Showcase ────────────────────────────────────────────────────────

function CodeShowcase() {
  return (
    <section className="border-t border-(--surface-border)">
      <div className="mx-auto max-w-6xl px-6 py-20">
        <div className="grid gap-12 lg:grid-cols-2 lg:items-center">
          {/* Left */}
          <div>
            <p className="mb-2 font-mono text-xs uppercase tracking-widest text-(--green)">
              SDK
            </p>
            <h2 className="mb-4 text-balance text-2xl font-bold tracking-tight text-foreground md:text-3xl">
              Write an exploit in under 10 lines.
            </h2>
            <p className="text-pretty leading-relaxed text-muted-foreground">
              {"CookieFarm's Python SDK handles everything from target iteration to flag submission. Just decorator "}
              <code className="rounded bg-(--surface-raised) px-1.5 py-0.5 font-mono text-sm text-(--green)">
                @exploit_manager
              </code>
              {","}
              <code className="rounded bg-(--surface-raised) px-1.5 py-0.5 font-mono text-sm text-(--green)">
                print(flag)
              </code>
              {", and run."}
            </p>

            <div className="mt-8 flex flex-col gap-3">
              {["Target iteration handled automatically", "Flags deduplicated before submission", "Parallel execution across all IPs"].map(
                (item) => (
                  <div key={item} className="flex items-start gap-2.5">
                    <span className="mt-0.5 flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-(--green)/15 text-(--green)">
                      <svg width="8" height="8" viewBox="0 0 8 8" fill="none">
                        <path
                          d="M1.5 4L3 5.5L6.5 2"
                          stroke="currentColor"
                          strokeWidth="1.5"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </svg>
                    </span>
                    <span className="text-sm leading-relaxed text-muted-foreground">
                      {item}
                    </span>
                  </div>
                )
              )}
            </div>
          </div>

          {/* Right: code block */}
          <div className="overflow-hidden rounded-lg border border-(--surface-border) bg-(--surface)">
            <div className="flex items-center gap-1.5 border-b border-(--surface-border) px-4 py-2.5">
              <span className="h-3 w-3 rounded-full bg-[#ff5f57]" />
              <span className="h-3 w-3 rounded-full bg-[#febc2e]" />
              <span className="h-3 w-3 rounded-full bg-[#28c840]" />
              <span className="ml-3 font-mono text-xs text-muted-foreground">
                exploit.py
              </span>
            </div>
            <pre className="p-5 text-sm leading-relaxed">
              <code className="font-mono">
                <span className="text-muted-foreground">#!/usr/bin/env python3</span>
                {"\n"}
                <span className="text-[oklch(0.7_0.12_260)]">import</span>
                <span className="text-foreground"> requests</span>
                {"\n"}
                <span className="text-[oklch(0.7_0.12_260)]">from</span>
                <span className="text-foreground"> cookiefarm </span>
                <span className="text-[oklch(0.7_0.12_260)]">import</span>
                <span className="text-foreground"> exploit_manager</span>
                {"\n\n"}
                <span className="text-foreground">@exploit_manager</span>
                {"\n"}
                <span className="text-[oklch(0.7_0.12_260)]">def</span>
                <span className="text-[oklch(0.85_0.14_90)]"> exploit</span>
                <span className="text-foreground">{"(ip: str, port: int, name: str):"}</span>
                {"\n"}
                <span className="text-foreground">{"    base_url = f\"http://{ip}:{port}\""}</span>
                {"\n\n"}
                <span className="text-foreground">{"    r = requests.get(f\"{base_url}/get-flag\")"}</span>
                {"\n"}
                <span className="text-foreground">{"    print(r.text)"}</span>
                {"\n\n"}
                <span className="text-muted-foreground"># run: </span>
                <span className="text-(--green)">
                  ckc exploit run -e exploit -n service
                </span>
              </code>
            </pre>
          </div>
        </div>
      </div>
    </section>
  );
}

// ─── Architecture ─────────────────────────────────────────────────────────

function Architecture() {
  const nodes = [
    {
      label: "Python Exploit",
      sub: "You write this",
      icon: Code2,
    },
    {
      label: "Go Client",
      sub: "CookieFarm runs this",
      icon: Server,
    },
    {
      label: "Scoreboard",
      sub: "Flags land here",
      icon: Flag,
    },
  ];

  return (
    <section className="border-t border-(--surface-border)">
      <div className="mx-auto max-w-6xl px-6 py-20">
        <p className="mb-2 font-mono text-xs uppercase tracking-widest text-(--green)">
          Architecture
        </p>
        <h2 className="mb-12 text-balance text-2xl font-bold tracking-tight text-foreground md:text-3xl">
          Simple by design.
        </h2>

        <div className="flex flex-col items-center gap-4 sm:flex-row sm:items-stretch sm:justify-center">
          {nodes.map((node, i) => (
            <div key={node.label} className="flex items-center gap-4">
              {/* Node box */}
              <div className="flex w-48 flex-col items-center rounded-lg border border-(--surface-border) bg-(--surface) p-6 text-center">
                <div className="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-md border border-(--green)/20 bg-(--green)/10 text-(--green)">
                  <node.icon size={18} />
                </div>
                <p className="font-mono text-sm font-semibold text-foreground">
                  {node.label}
                </p>
                <p className="mt-1 font-mono text-xs text-muted-foreground">
                  {node.sub}
                </p>
              </div>
              {/* Arrow between nodes */}
              {i < nodes.length - 1 && (
                <div className="flex shrink-0 items-center text-(--green)">
                  <svg
                    width="32"
                    height="16"
                    viewBox="0 0 32 16"
                    fill="none"
                    className="hidden sm:block"
                  >
                    <path
                      d="M0 8H28M28 8L20 2M28 8L20 14"
                      stroke="currentColor"
                      strokeWidth="1.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                  {/* Mobile: down arrow */}
                  <svg
                    width="16"
                    height="32"
                    viewBox="0 0 16 32"
                    fill="none"
                    className="block sm:hidden"
                  >
                    <path
                      d="M8 0V28M8 28L2 20M8 28L14 20"
                      stroke="currentColor"
                      strokeWidth="1.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ─── Footer CTA ───────────────────────────────────────────────────────────

function FooterCTA() {
  return (
    <section className="border-t border-(--surface-border)">
      <div className="mx-auto max-w-6xl px-6 py-24 text-center">
        <Terminal
          size={32}
          className="mx-auto mb-6 text-(--green)"
          strokeWidth={1.5}
        />
        <h2 className="mb-6 text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl">
          Ready to eat some cookies?
        </h2>
        <Link
          href="/docs"
          className="inline-flex items-center gap-2 rounded-md bg-(--green) px-6 py-3 font-mono text-sm font-semibold text-[oklch(0.1_0_0)] transition-opacity hover:opacity-90"
        >
          Read the Docs
          <ArrowRight size={14} />
        </Link>
      </div>
    </section>
  );
}

// ─── Footer ───────────────────────────────────────────────────────────────

function Footer() {
  return (
    <footer className="border-t border-(--surface-border)">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 py-8 sm:flex-row">
        <span className="font-mono text-sm font-semibold text-(--green)">
          CookieFarm
        </span>
        <p className="font-mono text-xs text-muted-foreground">
          Built by{" "}
          <Link
            href="https://github.com/ByteTheCookies"
            target="_blank"
            rel="noopener noreferrer"
            className="text-foreground transition-colors hover:text-(--green)"
          >
            ByteTheCookies
          </Link>
          {" · "}
          <Link
            href="https://github.com/ByteTheCookies/CookieFarm"
            target="_blank"
            rel="noopener noreferrer"
            className="text-foreground transition-colors hover:text-(--green)"
          >
            GitHub
          </Link>
        </p>
      </div>
    </footer>
  );
}

// ─── Page ─────────────────────────────────────────────────────────────────

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <main>
        <Hero />
        <FeatureGrid />
        <CodeShowcase />
        <Architecture />
        <FooterCTA />
      </main>
      <Footer />
    </div>
  );
}
