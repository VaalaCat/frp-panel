import { Server } from "@/lib/pb/common";
import { TypedClientPluginOptions } from "@/types/plugin";
import { HTTPProxyConfig, TypedProxyConfig } from "@/types/proxy";
import { ServerConfig } from "@/types/server";
import { ArrowRight, Globe, Monitor } from 'lucide-react';

export function VisitPreview({ server, typedProxyConfig, direction, withIcon = true }:
  {
    server: Server;
    typedProxyConfig: TypedProxyConfig;
    direction?: "row" | "column";
    withIcon?: boolean
  }) {
  return (
    <div className={"flex items-start sm:items-center justify-start p-2 text-xs font-mono text-nowrap " + (
      !direction ? "flex-wrap" : direction == "row" ? "flex-row" : "flex-col"
    )}>
      <ServerSideVisitPreview server={server} typedProxyConfig={typedProxyConfig} withIcon={withIcon} />
      <ArrowRight className="hidden sm:block w-4 h-4 text-gray-400 mx-2 flex-shrink-0" />
      <ClientSideVisitPreview typedProxyConfig={typedProxyConfig} withIcon={withIcon} />
    </div>
  );
}

export function ServerSideVisitPreview({ server, typedProxyConfig, withIcon = true }: { server: Server; typedProxyConfig: TypedProxyConfig; withIcon?: boolean }) {
  const serverCfg = JSON.parse(server?.config || '{}') as ServerConfig;
  const serverAddress = server.ip || serverCfg.bindAddr || 'Unknown';
  const serverPort = getServerPort(typedProxyConfig, serverCfg);

  return <div className="flex items-center mb-2 sm:mb-0">
    {withIcon && <Globe className="w-4 h-4 text-blue-500 mr-2 flex-shrink-0" />}
    <span className="text-nowrap">{typedProxyConfig.type == "http" ? "http://" : ""}{
      typedProxyConfig.type == "http" ? (
        getServerAuth(typedProxyConfig as HTTPProxyConfig) + getServerHost(typedProxyConfig as HTTPProxyConfig, serverCfg, serverAddress)
      ) : serverAddress
    }:{serverPort || "?"}{
        typedProxyConfig.type == "http" ?
          getServerPath(typedProxyConfig as HTTPProxyConfig) : ""
      }</span>
  </div>
}

export function ClientSideVisitPreview({ typedProxyConfig, withIcon = true }: { typedProxyConfig: TypedProxyConfig, withIcon?: boolean }) {
  const localAddress = typedProxyConfig.localIP || '127.0.0.1';
  const localPort = typedProxyConfig.localPort;
  const clientPlugin = typedProxyConfig.plugin;

  return <div className="flex items-center mb-2 sm:mb-0">
    {withIcon && <Monitor className="w-4 h-4 text-green-500 mr-2 flex-shrink-0" />}
    {clientPlugin && clientPlugin.type.length > 0 ?
      <PluginLocalDist plugins={clientPlugin} /> :
      <span className="text-nowrap">{localAddress}:{localPort}</span>}
  </div>
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

function getServerAuth(httpProxyConfig: HTTPProxyConfig) {
  if (!httpProxyConfig.httpUser || !httpProxyConfig.httpPassword) {
    return "";
  }
  return `${httpProxyConfig.httpUser}:${httpProxyConfig.httpPassword}@`
}

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
  return `[${httpProxyConfig.locations.join(", ")}]`;
}

function getServerHost(httpProxyConfig: HTTPProxyConfig, serverCfg: ServerConfig, serverAddress: string) {
  let allHosts = []
  if (httpProxyConfig.subdomain) {
    allHosts.push(`${httpProxyConfig.subdomain}.${serverCfg.subDomainHost}`);
  }

  allHosts.push(...(httpProxyConfig.customDomains || []));

  if (allHosts.length == 0) {
    return serverAddress;
  }

  if (allHosts.length == 1) {
    return allHosts[0];
  }

  return `[${allHosts.join(", ")}]`;
}

function PluginLocalDist({ plugins }: { plugins: TypedClientPluginOptions }) {
  return (<>
    {
      plugins.type === "unix_domain_socket" ? (
        <span className="text-nowrap">{plugins.unixPath}</span>
      ) : plugins.type === "static_file" ? (
        <span className="text-nowrap">{plugins.localPath}</span>
      ) : plugins.type === "http2https" || plugins.type === "https2http" || plugins.type === "https2https" ? (
        <span className="text-nowrap">{plugins.localAddr}</span>
      ) : (
        <span className="text-nowrap">{JSON.stringify(plugins)}</span>
      )
    }</>)
}