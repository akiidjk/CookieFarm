import { NextRequest, NextResponse } from 'next/server';

const publicPaths = ['/login', '/dashboard'];

export default function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  const isPublicPath = publicPaths.some(path => pathname === path || pathname.startsWith(path + '/'));

  const token = request.cookies.get('auth-token')?.value;

  const isLoggedIn = !!token;

  if (!isLoggedIn && !isPublicPath && pathname !== '/') {
    const redirectUrl = new URL('/login', request.url);
    return NextResponse.redirect(redirectUrl);
  }

  if (isLoggedIn && (isPublicPath || pathname === '/')) {
    const redirectUrl = new URL('/dashboard', request.url);
    return NextResponse.redirect(redirectUrl);
  }

  return NextResponse.next();
}

// Configure middleware to run on specific paths
export const config = {
  matcher: [
    /*
     * Match all paths except:
     * 1. /api routes
     * 2. /_next (Next.js internals)
     * 3. /fonts, /icons, /images (static files)
     * 4. all root files inside public (robots.txt, favicon.ico, etc.)
     */
    '/((?!api|_next|fonts|icons|images|[\\w-]+\\.\\w+).*)',
  ],
};
