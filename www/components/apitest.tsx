import { login, register } from '@/api/auth'
import { Button } from './ui/button'
import { deleteClient, getClient, initClient, listClient } from '@/api/client'
import { deleteServer, getServer, initServer, listServer } from '@/api/server'
import { updateFRPC, updateFRPS } from '@/api/frp'
import { ClientConfig } from '@/types/client'
import { ServerConfig } from '@/types/server'
import { getUserInfo, updateUserInfo } from '@/api/user'
import { Separator } from './ui/separator'
import { useState } from 'react'
import { Input } from './ui/input'
import { Label } from '@radix-ui/react-label'

export const APITest = () => {
  const [serverID, setServerID] = useState<string>('admin.server')
  const [clientID, setClientID] = useState<string>('admin.client')
  const [username, setUsername] = useState<string>('admin')
  const [password, setPassword] = useState<string>('admin')
  const [email, setEmail] = useState<string>('admin@localhost')

  return (
    <div className="flex flex-col w-full p-10 lg:w-1/2">
      <div className="grid grid-cols-2 sm:grid-cols-5 gap-4 my-4">
        <div>
          <Label>username</Label>
          <Input value={username} onChange={(e) => setUsername(e.target.value)} />
        </div>
        <div>
          <Label>password</Label>
          <Input value={password} onChange={(e) => setPassword(e.target.value)} />
        </div>
        <div>
          <Label>email</Label>
          <Input value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <div>
          <Label>clientID</Label>
          <Input value={clientID} onChange={(e) => setClientID(e.target.value)} />
        </div>
        <div>
          <Label>serverID</Label>
          <Input value={serverID} onChange={(e) => setServerID(e.target.value)} />
        </div>
      </div>
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 my-4">
        <Button
          onClick={async () => {
            console.log('attempting login:', await login({ username: username, password: password }))
          }}
        >
          login
        </Button>
        <Button
          onClick={async () => {
            console.log(
              'attempting register:',
              await register({ username: username, password: password, email: email }),
            )
          }}
        >
          register
        </Button>
        <Button
          onClick={async () => {
            console.log(
              'attempting update user:',
              await updateUserInfo({
                userInfo: { token: '123123' },
              }),
            )
          }}
        >
          update user
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting get user:', await getUserInfo({}))
          }}
        >
          get user
        </Button>
      </div>
      <Separator />
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 my-4">
        <Button
          onClick={async () => {
            console.log(
              'attempting init server:',
              await initServer({ serverId: serverID.replace(username + '.', ''), serverIp: '127.0.0.1' }),
            )
          }}
        >
          init server
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting delete server:', await deleteServer({ serverId: serverID }))
          }}
        >
          delete server
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting list servers:', await listServer({ page: 1, pageSize: 10 }))
          }}
        >
          list servers
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting get server:', await getServer({ serverId: serverID }))
          }}
        >
          get server
        </Button>
      </div>
      <Separator />
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 my-4">
        <Button
          onClick={async () => {
            console.log(
              'attempting update frps:',
              await updateFRPS({
                serverId: serverID,
                config: Buffer.from(
                  JSON.stringify({
                    bindPort: 1122,
                  } as ServerConfig),
                ),
              }),
            )
          }}
        >
          update frps
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting delete frps:', await deleteServer({ serverId: serverID }))
          }}
        >
          delete frps
        </Button>
      </div>
      <Separator />
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 my-4">
        <Button
          onClick={async () => {
            console.log('attempting init client:', await initClient({ clientId: clientID.replace(username + '.', '') }))
          }}
        >
          init client
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting delete client:', await deleteClient({ clientId: clientID }))
          }}
        >
          delete client
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting list clients:', await listClient({ page: 1, pageSize: 10 }))
          }}
        >
          list clients
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting get client:', await getClient({ clientId: clientID }))
          }}
        >
          get client
        </Button>
      </div>
      <Separator />
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 my-4">
        <Button
          onClick={async () => {
            console.log(
              'attempting update frpc:',
              await updateFRPC({
                clientId: clientID,
                serverId: serverID,
                config: Buffer.from(
                  JSON.stringify({
                    proxies: [{ name: 'test', type: 'tcp', localIP: '127.0.0.1', localPort: 1234, remotePort: 4321 }],
                  } as ClientConfig),
                ),
              }),
            )
          }}
        >
          update frpc
        </Button>
        <Button
          onClick={async () => {
            console.log('attempting delete frpc:', await deleteClient({ clientId: clientID }))
          }}
        >
          delete frpc
        </Button>
      </div>
    </div>
  )
}
