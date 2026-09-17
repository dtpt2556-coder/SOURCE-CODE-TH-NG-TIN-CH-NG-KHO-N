import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /**
   * Sinh `.next/standalone/server.js` de `Dockerfile.web` chay `node server.js`
   * ma khong can copy `node_modules` (architecture.md muc 2: container `web`).
   */
  output: "standalone",
};

export default nextConfig;
