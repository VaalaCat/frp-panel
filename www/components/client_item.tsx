import { Client } from "@/lib/pb/common"

export interface ClientItemProps {
    Client: Client
}
export const ClientItem: React.FC<ClientItemProps> = ({ Client }) => {
    return (<>
        <p className="text-sm text-muted-foreground">
            {Client.id}
        </p>
    </>)
}