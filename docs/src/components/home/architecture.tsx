import { ArrowRight, Server, Code, Send, Database } from 'lucide-react';

const steps = [
  {
    icon: Code,
    title: 'Write Exploit',
    description: 'Write your Python exploit using the simple CookieFarm API',
    color: 'text-sky-500',
    bgColor: 'bg-sky-500/10',
  },
  {
    icon: Server,
    title: 'Go Engine Executes',
    description: 'High-performance Go core runs your exploit against all targets',
    color: 'text-emerald-500',
    bgColor: 'bg-emerald-500/10',
  },
  {
    icon: Database,
    title: 'Flags Extracted',
    description: 'Automatic flag extraction and deduplication',
    color: 'text-amber-500',
    bgColor: 'bg-amber-500/10',
  },
  {
    icon: Send,
    title: 'Auto Submit',
    description: 'Flags submitted to the scoreserver automatically',
    color: 'text-rose-500',
    bgColor: 'bg-rose-500/10',
  },
];

export function Architecture() {
  return (
    <section className="py-24 bg-card/50 relative overflow-hidden">
      <div className="absolute inset-0 grid-pattern opacity-20" />

      {/* Decorative elements */}
      <div className="absolute top-0 left-1/4 w-64 h-64 bg-amber-500/5 rounded-full blur-3xl" />
      <div className="absolute bottom-0 right-1/4 w-64 h-64 bg-emerald-500/5 rounded-full blur-3xl" />

      <div className=" relative z-10 px-4">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <span className="inline-block px-4 py-1.5 rounded-full bg-emerald-500/10 text-emerald-400 text-sm font-medium mb-4">
            Architecture
          </span>
          <h2 className="text-3xl md:text-4xl font-bold mb-4 text-balance">
            Hybrid <span className="text-emerald-500">Go</span> + <span className="text-amber-500">Python</span> Power
          </h2>
          <p className="text-lg text-muted-foreground">
            The best of both worlds: Go&apos;s performance for the core engine,
            Python&apos;s flexibility for rapid exploit development.
          </p>
        </div>

        {/* Flow diagram */}
        <div className="flex flex-col lg:flex-row items-center justify-center gap-4 lg:gap-0">
          {steps.map((step, index) => (
            <div key={step.title} className="flex items-center">
              <div className="relative group">
                <div className="absolute -inset-2 bg-linear-to-r from-amber-500/20 to-emerald-500/20 rounded-xl blur opacity-0 group-hover:opacity-100 transition-opacity" />
                <div className="relative flex flex-col items-center p-6 bg-background border border-border rounded-xl w-64 hover:border-amber-500/50 transition-colors">
                  <div className={`inline-flex p-4 rounded-xl ${step.bgColor} mb-4`}>
                    <step.icon className={`size-8 ${step.color}`} />
                  </div>
                  <span className="text-xs text-muted-foreground mb-2">Step {index + 1}</span>
                  <h3 className="font-semibold text-foreground mb-2">{step.title}</h3>
                  <p className="text-sm text-muted-foreground text-center">{step.description}</p>
                </div>
              </div>

              {index < steps.length - 1 && (
                <ArrowRight className="hidden lg:block mx-4 size-6 text-muted-foreground" />
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
