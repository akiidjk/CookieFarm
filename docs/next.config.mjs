import { createMDX } from 'fumadocs-mdx/next';

const withMDX = createMDX();

/** @type {import('next').NextConfig} */
const config = {
  typedRoutes: true,
  reactCompiler: true,
  reactStrictMode: true,
  compress: true,
  output: 'standalone',
};

export default withMDX(config);
