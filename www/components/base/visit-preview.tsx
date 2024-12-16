import { Server } from "@/lib/pb/common";
import { HTTPProxyConfig, TypedProxyConfig } from "@/types/proxy";
import { ServerConfig } from "@/types/server";
import { ArrowRight, Globe, Monitor } from 'lucide-react';

export function VisitPreview({ server, typedProxyConfig }: { server: Server; typedProxyConfig: TypedProxyConfig }) {
  const serverCfg = JSON.parse(server?.config || '{}') as ServerConfig;
  const serverAddress = server.ip || serverCfg.bindAddr || 'Unknown';
  const serverPort = getServerPort(typedProxyConfig, serverCfg);
  const localAddress = typedProxyConfig.localIP || '127.0.0.1';
  const localPort = typedProxyConfig.localPort;

  function getServerPath(httpProxyConfig: HTTPProxyConfig) {
    if (!httpProxyConfig.locations) {
      return "";
    }
    if (httpProxyConfig.locations.length == 0) {
      return "";
    }
    if (httpProxyConfig.locations.length == 1) {
      return httpProxyConfig.locations[0];
    }
    return `[${httpProxyConfig.locations.join(",")}]`;
  }

  return (
    <div className="flex flex-col sm:flex-row items-start sm:items-center justify-start p-2 text-sm font-mono text-nowrap">
      <div className="flex items-center mb-2 sm:mb-0">
        <Globe className="w-4 h-4 text-blue-500 mr-2 flex-shrink-0" />
        <span className="text-nowrap">{typedProxyConfig.type == "http" ? "http://" : ""}{
          typedProxyConfig.type == "http" ? `${(typedProxyConfig as HTTPProxyConfig).subdomain}.${serverCfg.subDomainHost}` : serverAddress}:{
            serverPort}{typedProxyConfig.type == "http" ? getServerPath(typedProxyConfig as HTTPProxyConfig) : ""}</span>
      </div>
      <ArrowRight className="hidden sm:block w-4 h-4 text-gray-400 mx-2 flex-shrink-0" />
      <div className="flex items-center mb-2 sm:mb-0">
        <Monitor className="w-4 h-4 text-green-500 mr-2 flex-shrink-0" />
        <span className="text-nowrap">{localAddress}:{localPort}</span>
      </div>
    </div>
  );
}

function getServerPort(proxyConfig: TypedProxyConfig, serverConfig: ServerConfig): number | undefined {
  switch (proxyConfig.type) {
    case 'tcp':
      return (proxyConfig as any).remotePort;
    case 'udp':
      return (proxyConfig as any).remotePort;
    case 'http':
      return serverConfig.vhostHTTPPort;
    case 'https':
      return serverConfig.vhostHTTPSPort;
    default:
      return undefined;
  }
}

