/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    return {
      fallback: [
        {
          source: '/api/:path*',
          destination: `http://localhost:9000/api/:path*`,
        },
      ],
    }
  },
}

module.exports = nextConfig
