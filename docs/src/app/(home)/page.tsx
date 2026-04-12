import { Header } from '@/components/home/header';
import { Hero } from '@/components/home/hero';
import { Features } from '@/components/home/features';
import { Architecture } from '@/components/home/architecture';
import { Quickstart } from '@/components/home/quickstart';
import { CTA } from '@/components/home/cta';
import { Footer } from '@/components/home/footer';

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background">
      {/*<Header />*/}
      <main className="pt-16">
        <Hero />
        <Features />
        <Architecture />
        <Quickstart />
        <CTA />
      </main>
      <Footer />
    </div>
  );
}
