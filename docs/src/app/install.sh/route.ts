export const dynamic = 'force-dynamic'; // never cache at build time

const BRANCH = process.env.BRANCH || 'dev';

const GITHUB_RAW_URL =
  `https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/${BRANCH}/install.sh`;

export async function GET() {
  const res = await fetch(GITHUB_RAW_URL, {
    next: { revalidate: 60 },
  });

  if (!res.ok) {
    return new Response('Failed to fetch script', { status: 502 });
  }

  const content = await res.text();

  return new Response(content, {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      'Cache-Control': 'no-store',
      'X-Content-Type-Options': 'nosniff',
    },
  });
}
