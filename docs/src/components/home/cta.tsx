import Link from 'next/link';
import { ArrowRight, BookOpen, Cookie } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { SiGithub } from '@icons-pack/react-simple-icons';
export function CTA() {
  return (
    <section className="py-24 relative overflow-hidden">
      {/* Background effects */}
      <div className="absolute inset-0 bg-gradient-to-b from-background via-card to-background" />
      <div className="absolute inset-0 grid-pattern opacity-20" />
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-amber-500/10 rounded-full blur-3xl" />

      <div className=" relative z-10 px-4">
        <div className="max-w-3xl mx-auto text-center">
          <div className="inline-flex items-center justify-center size-20 rounded-2xl bg-amber-500/10 border border-amber-500/20 mb-8">
            <Cookie className="size-10 text-amber-500" />
          </div>

          <h2 className="text-4xl md:text-5xl font-bold mb-6 text-balance">
            Ready to <span className="text-amber-500">capture</span> some flags?
          </h2>

          <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            Join the ByteTheCookies community and start dominating your next CTF competition
            with CookieFarm.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Button asChild size="lg" className="bg-amber-500 hover:bg-amber-600 text-black font-semibold w-full sm:w-auto">
              <Link href="/docs">
                <BookOpen className="mr-2 size-4" />
                Read the Docs
              </Link>
            </Button>
            <Button asChild variant="outline" size="lg" className="border-border hover:bg-secondary w-full sm:w-auto">
              <a href="https://github.com/ByteTheCookies/CookieFarm" target="_blank" rel="noopener noreferrer">
                <SiGithub className="mr-2 size-4" />
                Star on GitHub
                <ArrowRight className="ml-2 size-4" />
              </a>
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
}
