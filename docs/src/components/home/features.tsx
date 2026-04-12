import {
  Zap,
  Shield,
  Code2,
  Gauge,
  Users,
  Terminal,
  GitBranch,
  Flag
} from 'lucide-react';

const features = [
  {
    icon: Zap,
    title: 'Zero Distraction',
    description: 'Focus only on writing exploit logic. CookieFarm handles flag submission, team management, and everything else.',
    color: 'text-amber-500',
    bgColor: 'bg-amber-500/10',
  },
  {
    icon: Code2,
    title: 'Go + Python Hybrid',
    description: 'High-performance Go core for network operations combined with Python flexibility for exploit development.',
    color: 'text-emerald-500',
    bgColor: 'bg-emerald-500/10',
  },
  {
    icon: Gauge,
    title: 'Real-time Dashboard',
    description: 'Monitor your attacks, track flag submissions, and view team statistics in real-time.',
    color: 'text-sky-500',
    bgColor: 'bg-sky-500/10',
  },
  {
    icon: Shield,
    title: 'Battle Tested',
    description: 'Used by ByteTheCookies in numerous CTF competitions. Proven reliability under pressure.',
    color: 'text-rose-500',
    bgColor: 'bg-rose-500/10',
  },
  {
    icon: Users,
    title: 'Team Coordination',
    description: 'Built-in tools for team collaboration, exploit sharing, and coordinated attacks.',
    color: 'text-violet-500',
    bgColor: 'bg-violet-500/10',
  },
  {
    icon: Terminal,
    title: 'Simple CLI',
    description: 'Powerful command-line interface for managing exploits, targets, and submissions.',
    color: 'text-orange-500',
    bgColor: 'bg-orange-500/10',
  },
  {
    icon: GitBranch,
    title: 'Open Source',
    description: 'Fully open source under MIT license. Contribute, customize, and make it your own.',
    color: 'text-teal-500',
    bgColor: 'bg-teal-500/10',
  },
  {
    icon: Flag,
    title: 'Flag Management',
    description: 'Automatic flag extraction, deduplication, and submission with configurable patterns.',
    color: 'text-pink-500',
    bgColor: 'bg-pink-500/10',
  },
];

export function Features() {
  return (
    <section className="py-24 relative">
      <div className="absolute inset-0 grid-pattern opacity-30" />

      <div className="relative z-10 px-16">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4 text-balance">
            Everything you need to{' '}
            <span className="text-amber-500">dominate</span> the competition
          </h2>
          <p className="text-lg text-muted-foreground">
            CookieFarm provides all the tools you need for Attack/Defense CTF competitions,
            so you can focus on what matters most: finding vulnerabilities and capturing flags.
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          {features.map((feature, index) => (
            <div
              key={feature.title}
              className="group relative p-6 rounded-xl bg-card border border-border hover:border-amber-500/50 transition-all duration-300"
              style={{ animationDelay: `${index * 100}ms` }}
            >
              <div className={`inline-flex p-3 rounded-lg ${feature.bgColor} mb-4`}>
                <feature.icon className={`size-6 ${feature.color}`} />
              </div>

              <h3 className="text-lg font-semibold mb-2 text-foreground group-hover:text-amber-500 transition-colors">
                {feature.title}
              </h3>

              <p className="text-sm text-muted-foreground leading-relaxed">
                {feature.description}
              </p>

              {/* Hover glow effect */}
              <div className="absolute inset-0 rounded-xl bg-linear-to-br from-amber-500/5 to-emerald-500/5 opacity-0 group-hover:opacity-100 transition-opacity -z-10" />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
