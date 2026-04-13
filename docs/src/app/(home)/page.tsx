"use client"
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
import { motion, useInView } from "framer-motion";

import { SiGithub } from '@icons-pack/react-simple-icons';
import { Carousel } from "@/components/carousel";
import { useRef } from "react";
import { MotionSection } from "@/components/motionwrapper";


// ─── Animation Variants ──────────────────────────────────────────────────────

const fadeInUp = {
  hidden: { opacity: 0, y: 30 },
  visible: { opacity: 1, y: 0 },
};

const fadeIn = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
};

const scaleIn = {
  hidden: { opacity: 0, scale: 0.95 },
  visible: { opacity: 1, scale: 1 },
};

const slideInLeft = {
  hidden: { opacity: 0, x: -40 },
  visible: { opacity: 1, x: 0 },
};

const slideInRight = {
  hidden: { opacity: 0, x: 40 },
  visible: { opacity: 1, x: 0 },
};

const staggerContainer = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.1,
    },
  },
};

const staggerItem = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
};





// ─── Nav ────────────────────────────────────────────────────────────────────

export function Nav() {
  return (
    <motion.header
      className="sticky top-0 z-50 border-b border-(--surface-border) bg-background/80 backdrop-blur-sm"
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
    >
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-2">
          <motion.span
            className="font-mono text-base font-bold text-(--green)"
            whileHover={{ scale: 1.05 }}
            transition={{ type: "spring", stiffness: 400, damping: 10 }}
          >
            CookieFarm
          </motion.span>
        </Link>

        {/* Nav links */}
        <nav className="flex items-center gap-6">
          <motion.div whileHover={{ y: -2 }} transition={{ type: "spring", stiffness: 400, damping: 10 }}>
            <Link
              href="/docs"
              className="font-mono text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              Docs
            </Link>
          </motion.div>
          <motion.div whileHover={{ y: -2 }} transition={{ type: "spring", stiffness: 400, damping: 10 }}>
            <Link
              href="https://github.com/ByteTheCookies/CookieFarm"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-1.5 font-mono text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              <SiGithub size={14} />
              GitHub
            </Link>
          </motion.div>
        </nav>
      </div>
    </motion.header>
  );
}



// ─── Hero ────────────────────────────────────────────────────────────────────

function Hero() {
  return (
    <section className="mx-auto max-w-6xl px-6 pb-20 pt-24 md:pt-32">
      <motion.div
        initial="hidden"
        animate="visible"
        variants={staggerContainer}
      >
        <motion.div
          variants={fadeInUp}
          transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
          className="mb-4 inline-flex items-center gap-2 rounded-full border border-(--green)/30 bg-(--green)/5 px-3 py-1"
        >
          <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-(--green)" />
          <span className="font-mono text-xs text-(--green)">
            Attack / Defense CTF Framework
          </span>
        </motion.div>

        <motion.h1
          variants={fadeInUp}
          transition={{ duration: 0.6, delay: 0.1, ease: [0.22, 1, 0.36, 1] }}
          className="mb-6 max-w-3xl text-balance text-4xl font-bold leading-tight tracking-tight text-foreground md:text-5xl lg:text-6xl"
        >
          The exploit farm that stays{" "}
          <span className="text-(--green)">out of your way.</span>
        </motion.h1>

        <motion.p
          variants={fadeInUp}
          transition={{ duration: 0.6, delay: 0.2, ease: [0.22, 1, 0.36, 1] }}
          className="mb-10 max-w-2xl text-pretty text-lg leading-relaxed text-muted-foreground"
        >
          CookieFarm is an A/D CTF framework with a Go server, Python client SDK,
          and zero-config flag submission. Write the exploit. We handle the rest.
        </motion.p>

        {/* CTAs */}
        <motion.div
          variants={fadeInUp}
          transition={{ duration: 0.6, delay: 0.3, ease: [0.22, 1, 0.36, 1] }}
          className="mb-14 flex flex-wrap items-center gap-3"
        >
          <Link
            href="/docs"
            className="group inline-flex items-center gap-2 rounded-md bg-(--green) px-5 py-2.5 font-mono text-sm font-semibold text-[oklch(0.1_0_0)] transition-all duration-300 hover:scale-105 hover:shadow-[0_0_20px_4px_oklch(0.55_0.15_145_/_0.3)]"
          >
            Get Started
            <ArrowRight size={14} className="transition-transform duration-300 group-hover:translate-x-1" />
          </Link>
          <Link
            href="https://github.com/ByteTheCookies/CookieFarm"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 rounded-md border border-(--surface-border) bg-(--surface) px-5 py-2.5 font-mono text-sm text-foreground transition-all duration-300 hover:scale-105 hover:border-[var(--green)]/50 hover:text-(--green)"
          >
            <SiGithub size={14} />
            View on GitHub
          </Link>
        </motion.div>

        {/* Carousel */}
        <motion.div
          variants={scaleIn}
          transition={{ duration: 0.7, delay: 0.4, ease: [0.22, 1, 0.36, 1] }}
        >
          <Carousel />
        </motion.div>
      </motion.div>
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
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="border-t border-(--surface-border)">
      <div className="mx-auto max-w-6xl px-6 py-20">
        <MotionSection>
          <p className="mb-2 font-mono text-xs uppercase tracking-widest text-(--green)">
            Features
          </p>
          <h2 className="mb-12 text-balance text-2xl font-bold tracking-tight text-foreground md:text-3xl">
            Everything handled. Nothing in the way.
          </h2>
        </MotionSection>

        <motion.div
          ref={ref}
          initial="hidden"
          animate={isInView ? "visible" : "hidden"}
          variants={staggerContainer}
          className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3"
        >
          {features.map((f, index) => (
            <motion.div
              key={f.title}
              variants={staggerItem}
              transition={{ duration: 0.5, delay: index * 0.08, ease: [0.22, 1, 0.36, 1] }}
              className="group rounded-lg border border-(--surface-border) bg-(--surface) p-6 transition-all duration-300 hover:-translate-y-1 hover:border-[var(--green)]/40 hover:shadow-[0_4px_20px_-4px_oklch(0.55_0.15_145_/_0.15)]"
            >
              <motion.div
                className="mb-4 inline-flex h-9 w-9 items-center justify-center rounded-md border border-[var(--green)]/20 bg-(--green)/10 text-(--green) transition-all duration-300 group-hover:scale-110 group-hover:bg-(--green)/20"
                whileHover={{ rotate: [0, -10, 10, 0] }}
                transition={{ duration: 0.4 }}
              >
                <f.icon size={16} />
              </motion.div>
              <h3 className="mb-2 font-mono text-sm font-semibold text-foreground">
                {f.title}
              </h3>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {f.description}
              </p>
            </motion.div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}


// ─── Code Showcase ────────────────────────────────────────────────────────

function CodeShowcase() {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });
  const checklistItems = ["Target iteration handled automatically", "Flags deduplicated before submission", "Parallel execution across all IPs"];

  return (
    <section className="border-t border-(--surface-border)" ref={ref}>
      <div className="mx-auto max-w-6xl px-6 py-20">
        <div className="grid gap-12 lg:grid-cols-2 lg:items-center">
          {/* Left */}
          <motion.div
            initial="hidden"
            animate={isInView ? "visible" : "hidden"}
            variants={slideInLeft}
            transition={{ duration: 0.7, ease: [0.22, 1, 0.36, 1] }}
          >
            <p className="mb-2 font-mono text-xs uppercase tracking-widest text-(--green)">
              SDK
            </p>
            <h2 className="mb-4 text-balance text-2xl font-bold tracking-tight text-foreground md:text-3xl">
              Write an exploit in under 10 lines.
            </h2>
            <p className="text-pretty leading-relaxed text-muted-foreground">
              {"CookieFarm's Python SDK handles everything from target iteration to flag submission. Just subclass "}
              <code className="rounded bg-(--surface-raised) px-1.5 py-0.5 font-mono text-sm text-(--green)">
                ExploitBase
              </code>
              {", implement "}
              <code className="rounded bg-(--surface-raised) px-1.5 py-0.5 font-mono text-sm text-(--green)">
                attack()
              </code>
              {", and run."}
            </p>

            <motion.div
              className="mt-8 flex flex-col gap-3"
              initial="hidden"
              animate={isInView ? "visible" : "hidden"}
              variants={staggerContainer}
            >
              {checklistItems.map((item, index) => (
                <motion.div
                  key={item}
                  className="flex items-start gap-2.5"
                  variants={staggerItem}
                  transition={{ duration: 0.4, delay: 0.3 + index * 0.1 }}
                >
                  <motion.span
                    className="mt-0.5 flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-(--green)/15 text-(--green)"
                    initial={{ scale: 0 }}
                    animate={isInView ? { scale: 1 } : { scale: 0 }}
                    transition={{ duration: 0.3, delay: 0.5 + index * 0.1, type: "spring" }}
                  >
                    <svg width="8" height="8" viewBox="0 0 8 8" fill="none">
                      <path
                        d="M1.5 4L3 5.5L6.5 2"
                        stroke="currentColor"
                        strokeWidth="1.5"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      />
                    </svg>
                  </motion.span>
                  <span className="text-sm leading-relaxed text-muted-foreground">
                    {item}
                  </span>
                </motion.div>
              ))}
            </motion.div>
          </motion.div>

          {/* Right: code block */}
          <motion.div
            className="overflow-hidden rounded-lg border border-(--surface-border) bg-(--surface) transition-all duration-300 hover:border-(--green)/30 hover:shadow-[0_4px_30px_-8px_oklch(0.55_0.15_145_/_0.2)]"
            initial="hidden"
            animate={isInView ? "visible" : "hidden"}
            variants={slideInRight}
            transition={{ duration: 0.7, delay: 0.2, ease: [0.22, 1, 0.36, 1] }}
          >
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
          </motion.div>
        </div>
      </div>
    </section>
  );
}

// ─── Architecture ─────────────────────────────────────────────────────────

function Architecture() {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  const nodes = [
    {
      label: "Python Exploit",
      sub: "You write this",
      icon: Code2,
    },
    {
      label: "Go Server",
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
    <section className="border-t border-(--surface-border)" ref={ref}>
      <div className="mx-auto max-w-6xl px-6 py-20">
        <MotionSection>
          <p className="mb-2 font-mono text-xs uppercase tracking-widest text-(--green)">
            Architecture
          </p>
          <h2 className="mb-12 text-balance text-2xl font-bold tracking-tight text-foreground md:text-3xl">
            Simple by design.
          </h2>
        </MotionSection>

        <div className="flex flex-col items-center gap-4 sm:flex-row sm:items-stretch sm:justify-center">
          {nodes.map((node, i) => (
            <motion.div
              key={node.label}
              className="flex items-center gap-4"
              initial="hidden"
              animate={isInView ? "visible" : "hidden"}
              variants={fadeInUp}
              transition={{ duration: 0.6, delay: 0.2 + i * 0.15, ease: [0.22, 1, 0.36, 1] }}
            >
              {/* Node box */}
              <motion.div
                className="group flex w-48 flex-col items-center rounded-lg border border-(--surface-border) bg-(--surface) p-6 text-center transition-all duration-300 hover:-translate-y-1 hover:border-[var(--green)]/40 hover:shadow-[0_4px_20px_-4px_oklch(0.55_0.15_145_/_0.15)]"
                whileHover={{ scale: 1.02 }}
              >
                <motion.div
                  className="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-md border border-[var(--green)]/20 bg-(--green)/10 text-(--green) transition-transform duration-300 group-hover:scale-110"
                  initial={{ rotate: -10, opacity: 0 }}
                  animate={isInView ? { rotate: 0, opacity: 1 } : { rotate: -10, opacity: 0 }}
                  transition={{ duration: 0.5, delay: 0.4 + i * 0.15 }}
                >
                  <node.icon size={18} />
                </motion.div>
                <p className="font-mono text-sm font-semibold text-foreground">
                  {node.label}
                </p>
                <p className="mt-1 font-mono text-xs text-muted-foreground">
                  {node.sub}
                </p>
              </motion.div>
              {/* Arrow between nodes */}
              {i < nodes.length - 1 && (
                <motion.div
                  className="flex shrink-0 items-center text-(--green)"
                  initial={{ opacity: 0, scale: 0.5 }}
                  animate={isInView ? { opacity: 1, scale: 1 } : { opacity: 0, scale: 0.5 }}
                  transition={{ duration: 0.4, delay: 0.5 + i * 0.2 }}
                >
                  <motion.svg
                    width="32"
                    height="16"
                    viewBox="0 0 32 16"
                    fill="none"
                    className="hidden sm:block"
                    animate={{ x: [0, 4, 0] }}
                    transition={{ duration: 1.5, repeat: Infinity, ease: "easeInOut" }}
                  >
                    <path
                      d="M0 8H28M28 8L20 2M28 8L20 14"
                      stroke="currentColor"
                      strokeWidth="1.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </motion.svg>
                  {/* Mobile: down arrow */}
                  <motion.svg
                    width="16"
                    height="32"
                    viewBox="0 0 16 32"
                    fill="none"
                    className="block sm:hidden"
                    animate={{ y: [0, 4, 0] }}
                    transition={{ duration: 1.5, repeat: Infinity, ease: "easeInOut" }}
                  >
                    <path
                      d="M8 0V28M8 28L2 20M8 28L14 20"
                      stroke="currentColor"
                      strokeWidth="1.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </motion.svg>
                </motion.div>
              )}
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ─── Footer CTA ───────────────────────────────────────────────────────────

function FooterCTA() {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="border-t border-(--surface-border)" ref={ref}>
      <div className="mx-auto max-w-6xl px-6 py-24 text-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 20 }}
          transition={{ duration: 0.6, ease: [0.22, 1, 0.36, 1] }}
        >
          <motion.div
            animate={{ y: [0, -8, 0] }}
            transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
          >
            <Terminal
              size={32}
              className="mx-auto mb-6 text-(--green)"
              strokeWidth={1.5}
            />
          </motion.div>
          <motion.h2
            className="mb-6 text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl"
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 20 }}
            transition={{ duration: 0.6, delay: 0.15, ease: [0.22, 1, 0.36, 1] }}
          >
            Ready to dominate the scoreboard?
          </motion.h2>
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={isInView ? { opacity: 1, scale: 1 } : { opacity: 0, scale: 0.9 }}
            transition={{ duration: 0.5, delay: 0.3, type: "spring" }}
          >
            <Link
              href="/docs"
              className="group inline-flex items-center gap-2 rounded-md bg-(--green) px-6 py-3 font-mono text-sm font-semibold text-[oklch(0.1_0_0)] transition-all duration-300 hover:scale-105 hover:shadow-[0_0_24px_6px_oklch(0.55_0.15_145_/_0.35)]"
            >
              Read the Docs
              <ArrowRight size={14} className="transition-transform duration-300 group-hover:translate-x-1" />
            </Link>
          </motion.div>
        </motion.div>
      </div>
    </section>
  );
}

// ─── Footer ───────────────────────────────────────────────────────────────

function Footer() {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-50px" });

  return (
    <motion.footer
      ref={ref}
      className="border-t border-(--surface-border)"
      initial={{ opacity: 0 }}
      animate={isInView ? { opacity: 1 } : { opacity: 0 }}
      transition={{ duration: 0.5 }}
    >
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 py-8 sm:flex-row">
        <motion.span
          className="font-mono text-sm font-semibold text-(--green)"
          whileHover={{ scale: 1.05 }}
        >
          CookieFarm
        </motion.span>
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
    </motion.footer>
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
